package configly

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/zanedma/configly/pkg/log"
)

const (
	defaultTagKey = "configly"
)

type tagOptions struct {
	key          string
	fieldIdx     int
	required     bool
	defaultValue string
	min          *int64
	max          *int64
	minLen       *int
	maxLen       *int
	sourceName   string
	rawVal       string
	// TODO pattern
}

type Loader[T any] struct {
	tagKey  string
	sources []Source
	// val     reflect.Value
	logger zerolog.Logger
}

type LoaderConfig struct {
	TagKey  string
	Sources []Source
}

func New[T any](cfg LoaderConfig) (*Loader[T], error) {
	if len(cfg.Sources) == 0 {
		return nil, errors.New("at least one source is required")
	}

	var loaderCfgInstance T

	val := reflect.ValueOf(&loaderCfgInstance).Elem()
	valType := val.Type()
	loadLogger := log.GetBase().With().Str("component", "load").Logger()
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
		// val:     val,
		logger: logger,
	}, nil
}

func (l *Loader[T]) Load() (*T, error) {
	var cfg T
	val := reflect.ValueOf(&cfg).Elem()
	typ := val.Type()
	// parse all tags first, so that if there are any invalid/inproperly formatted
	// tags, we can return all errors in one
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

		// field := typ.Field(opts.fieldIdx)
		fieldValue := val.Field(opts.fieldIdx)
		if err := l.setField(&fieldValue, value); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("error setting %s (source %s): %w", opts.key, sourceName, err))
			continue
		}
	}

	if len(validationErrors) > 0 {
		return nil, errors.Join(validationErrors...)
	}

	return &cfg, nil
}

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

func parseMinMax(prefixKey, part string) (int64, error) {
	str := strings.TrimPrefix(part, fmt.Sprintf("%s=", prefixKey))
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func parseLen(prefixKey, part string) (int, error) {
	str := strings.TrimPrefix(part, fmt.Sprintf("%s=", prefixKey))
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (l *Loader[T]) getValuesFromSources(tagOpts map[string]*tagOptions) error {
	var sourceErrors []error
	for key := range tagOpts {
		for _, source := range l.sources {
			val, found, err := source.GetValue(key)
			if err != nil {
				sourceErrors = append(sourceErrors, fmt.Errorf("error getting value for key %s from source %s: %w", key, source.Name(), err))
				continue
			}
			if found {
				tagOpts[key].sourceName = source.Name()
				tagOpts[key].rawVal = val
				break
			}
		}
	}
	if len(sourceErrors) > 0 {
		return errors.Join(sourceErrors...)
	}
	return nil
}

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
