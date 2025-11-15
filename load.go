package configly

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/zanedma/configly/sources"
)

const (
	// defaultTagKey is the struct tag key used when none is specified in LoaderConfig.
	defaultTagKey = "configly"
)

// tagOptions represents parsed options from a struct field's tag.
// It contains the configuration key, field index, validation constraints,
// and whether the field is required.
type tagOptions struct {
	key          string // The key to look up in configuration sources
	fieldIdx     int    // Index of the field in the struct
	required     bool   // Whether this field must have a value
	defaultValue string // Default value if not found in sources
	min          *int64 // Minimum value for numeric types
	max          *int64 // Maximum value for numeric types
	minLen       *int   // Minimum length for string types
	maxLen       *int   // Maximum length for string types
	// TODO pattern
}

// Loader is a generic configuration loader for type T.
// It retrieves values from multiple sources in priority order,
// validates constraints, and populates a struct instance.
type Loader[T any] struct {
	tagKey  string           // The struct tag key to use for field configuration
	sources []sources.Source // Configuration sources in priority order
	logger  zerolog.Logger   // Logger for debugging and warnings
}

// LoaderConfig contains configuration options for creating a new Loader.
type LoaderConfig struct {
	TagKey  string   // The struct tag key to use (defaults to "configly" if empty)
	Sources []sources.Source // Configuration sources in priority order (first source wins)
}

// New creates a new Loader instance for type T.
// It validates that T is a struct type and that at least one source is provided.
// The TagKey in LoaderConfig specifies which struct tag to use (defaults to "configly").
// Returns an error if no sources are provided or if T is not a struct type.
func New[T any](cfg LoaderConfig) (*Loader[T], error) {
	if len(cfg.Sources) == 0 {
		return nil, errors.New("at least one source is required")
	}

	var loaderCfgInstance T

	val := reflect.ValueOf(&loaderCfgInstance).Elem()
	valType := val.Type()
	loadLogger := getBaseLogger().With().Str("component", "load").Logger()
	loadLogger.Debug().Msgf("validating type '%s'", valType.Name())
	kind := valType.Kind()

	if kind != reflect.Struct {
		return nil, fmt.Errorf("invalid type for %s: %s (must be struct)", valType.Name(), kind)
	}

	logger := loadLogger.With().Str("type", valType.Name()).Logger()
	logger.Debug().Any("loaderConfig", cfg).Msg("successfully initialized")

	tagKey := cfg.TagKey
	if tagKey == "" {
		tagKey = defaultTagKey
	}

	return &Loader[T]{
		tagKey:  tagKey,
		sources: cfg.Sources,
		logger:  logger,
	}, nil
}

// Load loads configuration values from sources into a new instance of type T.
// It first parses all struct tags to identify fields and their constraints,
// then retrieves values from sources in priority order (first source wins),
// applies defaults when values are not found, and validates all constraints.
// Returns a fully populated and validated configuration instance or an error
// containing all validation failures joined together.
func (l *Loader[T]) Load() (*T, error) {
	var cfg T
	val := reflect.ValueOf(&cfg).Elem()
	typ := val.Type()
	// parse all tags first, so that if there are any invalid/inproperly formatted
	// tags, we can return all errors in one before attempting to fetch values from
	// the sources which adds unnecessary runtime (trade-off is that the generic T
	// must have valid tags before the user knows if there are any issues with the
	// actual values stored in the sources)
	tagOpts, err := l.parseAllTags(typ.NumField(), val)
	if err != nil {
		return nil, err
	}

	var validationErrors []error
	for _, opts := range tagOpts {
		value, sourceName, found := l.getValueFromSources(opts.key)
		if !found && opts.required {
			validationErrors = append(validationErrors, fmt.Errorf("required value %s not found in provided sources", opts.key))
			continue
		}

		if !found && opts.defaultValue != "" {
			value = opts.defaultValue
			found = true
		}

		if !found {
			continue
		}

		fieldValue := val.Field(opts.fieldIdx)
		if err := l.setField(&fieldValue, value); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("error setting %s (source %s): %w", opts.key, sourceName, err))
			continue
		}

		err = l.validateField(fieldValue, opts)
		if err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		return nil, errors.Join(validationErrors...)
	}

	return &cfg, nil
}

// parseAllTags parses struct tags for all fields in the configuration type.
// It skips unexported fields and fields without tags. If any tag has invalid
// formatting (e.g., invalid min/max values), all parsing errors are joined
// and returned together. Returns a slice of tagOptions for valid tagged fields.
func (l *Loader[T]) parseAllTags(numFields int, val reflect.Value) ([]tagOptions, error) {
	var parseErrors []error
	var allOpts []tagOptions
	for idx := range numFields {
		field := val.Type().Field(idx)
		fieldValue := val.Field(idx)

		if !fieldValue.CanSet() {
			l.logger.Debug().
				Str("key", field.Name).
				Msg("skipping unexported field")
			continue
		}

		tag := field.Tag.Get(l.tagKey)
		if tag == "" {
			l.logger.Debug().
				Str("field", field.Name).
				Msgf("no %s tag found, skipping", l.tagKey)
			continue
		}

		tagOpts, tagWarnings := l.parseTag(tag)
		if len(tagWarnings) > 0 {
			parseErrors = append(parseErrors, tagWarnings...)
		} else {
			tagOpts.fieldIdx = idx
			allOpts = append(allOpts, tagOpts)
		}
	}

	if len(parseErrors) > 0 {
		return nil, errors.Join(parseErrors...)
	}

	return allOpts, nil
}

// parseTag parses a single struct tag string into tagOptions.
// Tag format: "key,option1,option2=value"
// Supported options: required, default=value, min=int, max=int, minLen=int, maxLen=int
// Returns the parsed options and a slice of errors for any invalid option values.
// Whitespace around options is automatically trimmed.
func (l *Loader[T]) parseTag(tag string) (tagOptions, []error) {
	tagLogger := l.logger.With().Str("func", "parseTag").Str("tag", tag).Logger()
	parts := strings.Split(tag, ",")
	tagLogger.Debug().Strs("parts", parts).Send()
	opts := tagOptions{
		key: parts[0],
	}
	var warnings []error
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		switch {
		case part == "required":
			opts.required = true
		case strings.HasPrefix(part, "default="):
			opts.defaultValue = strings.TrimPrefix(part, "default=")
		case strings.HasPrefix(part, "min="):
			if val, err := parseMinMax("min", part); err != nil {
				warning := fmt.Errorf("invalid minimum value: %w", err)
				warnings = append(warnings, warning)
				tagLogger.Warn().Err(warning).Send()
			} else {
				opts.min = &val
			}
		case strings.HasPrefix(part, "max="):
			if val, err := parseMinMax("max", part); err != nil {
				warning := fmt.Errorf("invalid maximum value: %w", err)
				warnings = append(warnings, warning)
				tagLogger.Warn().Err(warning).Send()
			} else {
				opts.max = &val
			}
		case strings.HasPrefix(part, "minLen="):
			if val, err := parseLen("minLen", part); err != nil {
				warning := fmt.Errorf("invalid min length value %w", err)
				warnings = append(warnings, warning)
				tagLogger.Warn().Err(warning).Send()
			} else {
				opts.minLen = &val
			}
		case strings.HasPrefix(part, "maxLen="):
			if val, err := parseLen("maxLen", part); err != nil {
				warning := fmt.Errorf("invalid max length value %w", err)
				warnings = append(warnings, warning)
				tagLogger.Warn().Err(warning).Send()
			} else {
				opts.maxLen = &val
			}
		}
	}
	return opts, warnings
}

// parseMinMax parses a min or max value from a tag option part.
// The part should be in the format "min=123" or "max=456".
// Returns the parsed int64 value or an error if parsing fails.
func parseMinMax(prefixKey, part string) (int64, error) {
	str := strings.TrimPrefix(part, fmt.Sprintf("%s=", prefixKey))
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// parseLen parses a minLen or maxLen value from a tag option part.
// The part should be in the format "minLen=5" or "maxLen=50".
// Returns the parsed int value or an error if parsing fails.
func parseLen(prefixKey, part string) (int, error) {
	str := strings.TrimPrefix(part, fmt.Sprintf("%s=", prefixKey))
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// getValueFromSources retrieves a value for the given key from configured sources.
// Sources are checked in order, and the first source that returns a value wins.
// Sources that return errors are logged and skipped.
// Returns the value, the source name it came from, and whether a value was found.
func (l *Loader[T]) getValueFromSources(key string) (string, string, bool) {
	logger := l.logger.With().Str("func", "getValueFromSources").Str("key", key).Logger()
	for _, source := range l.sources {
		val, found, err := source.GetValue(key)
		if err != nil {
			logger.Warn().Str("source", source.Name()).Err(err)
			continue
		}
		if found {
			logger.Debug().Str("source", source.Name()).Msgf("found value %s", val)
			return val, source.Name(), true
		}
	}
	return "", "", false
}

// setField sets a struct field value by parsing a string value into the appropriate type.
// Supported types: string, all int types, all uint types, all float types, bool, and time.Duration.
// For time.Duration, the string must be in a format parseable by time.ParseDuration (e.g., "5s", "1h30m").
// Returns an error if the string cannot be parsed into the field's type.
func (l *Loader[T]) setField(value *reflect.Value, strVal string) error {
	switch value.Kind() {
	case reflect.String:
		value.SetString(strVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(strVal)
			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}
			value.SetInt(int64(duration))
			return nil
		}
		intVal, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer: %w", err)
		}
		value.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer: %w", err)
		}
		value.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return fmt.Errorf("invalid float: %w", err)
		}
		value.SetFloat(floatVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strVal)
		if err != nil {
			return fmt.Errorf("invalid boolean: %w", err)
		}
		value.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", value.Kind())
	}
	return nil
}

// validateField validates a field value against the constraints specified in its tag options.
// For strings: validates minLen and maxLen if specified.
// For integers (signed and unsigned): validates min and max if specified.
// For floats: validates min and max if specified.
// Other types (bool, etc.) have no validation constraints.
// Returns an error describing the first constraint violation, or nil if all constraints are satisfied.
func (l *Loader[T]) validateField(field reflect.Value, opts tagOptions) error {
	switch field.Kind() {
	case reflect.String:
		str := field.String()
		strLen := len(str)
		if opts.minLen != nil && strLen < *opts.minLen {
			return fmt.Errorf("string length %d less than minimum %d", strLen, *opts.minLen)
		}

		if opts.maxLen != nil && strLen > *opts.maxLen {
			return fmt.Errorf("string length %d exceeds maximum %d", strLen, *opts.maxLen)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := field.Int()

		if opts.min != nil && val < *opts.min {
			return fmt.Errorf("integer value %d is less than minimum %d", val, *opts.min)
		}

		if opts.max != nil && val > *opts.max {
			return fmt.Errorf("integer value %d exceeds maximum %d", val, *opts.max)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := field.Uint()

		if opts.min != nil && val < uint64(*opts.min) {
			return fmt.Errorf("unsigned integer value %d is less than minimum %d", val, *opts.min)
		}

		if opts.max != nil && val > uint64(*opts.max) {
			return fmt.Errorf("unsigned integer value %d exceeds maximum %d", val, *opts.max)
		}

	case reflect.Float32, reflect.Float64:
		val := field.Float()

		if opts.min != nil && val < float64(*opts.min) {
			return fmt.Errorf("float value %f is less than minimum %d", val, *opts.min)
		}

		if opts.max != nil && val > float64(*opts.max) {
			return fmt.Errorf("float value %f exceeds maximum %d", val, *opts.max)
		}
	}

	return nil
}
