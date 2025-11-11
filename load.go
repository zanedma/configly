package configly

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/zanedma/configly/pkg/log"
)

const (
	defaultTagKey = "configly"
)

var (
	errNoSource    = errors.New("at least one source is required")
	ErrInvalidType = errors.New("generic type passed to Load must be a struct")
	loadLogger     = log.GetBase().With().Str("component", "load").Logger()
)

type tagOptions struct {
	key          string
	required     bool
	defaultValue string
	min          *int64
	max          *int64
	minLen       *int
	maxLen       *int
	// TODO pattern
}

type Loader[T any] struct {
	tagKey  string
	sources []Source
	val     reflect.Value
	valType reflect.Type
	logger  zerolog.Logger
}

type LoaderConfig struct {
	TagKey  string
	Sources []Source
}

func New[T any](cfg LoaderConfig) (*Loader[T], error) {
	if len(cfg.Sources) == 0 {
		return nil, errNoSource
	}

	var loaderCfgInstance T

	val := reflect.ValueOf(&loaderCfgInstance).Elem()
	valType := val.Type()
	loadLogger.Debug().Msgf("validating type '%s'", valType.Name())
	kind := reflect.TypeOf(val).Kind()

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
		tagKey:  defaultTagKey,
		sources: cfg.Sources,
		val:     val,
		valType: valType,
		logger:  logger,
	}, nil
}

func (l *Loader[T]) Load() (T, error) {
	var cfg T
  _, err := l.parseAllTags()
  if err != nil {
    return cfg, err
  }
	// get values from sources
  // validate values (could be done in above)
	return cfg, nil
}

func (l *Loader[T]) parseAllTags() ([]tagOptions, error) {
	var parseErrors []error
	var allOpts []tagOptions
	for idx := 0; idx < l.valType.NumField(); idx++ {
		field := l.valType.Field(idx)
		fieldValue := l.val.Field(idx)

		if !fieldValue.CanSet() {
			loadLogger.Debug().
				Str("key", field.Name).
				Msg("skipping unexported field")
			continue
		}

		tag := field.Tag.Get("configly")
		if tag == "" {
			loadLogger.Debug().
				Str("field", field.Name).
				Msg("no configly tag found, skipping")
			continue
		}

		tagOpts, tagWarnings := parseTag(tag)
		if len(tagWarnings) > 0 {
			parseErrors = append(parseErrors, tagWarnings...)
		} else {
			allOpts = append(allOpts, tagOpts)
		}
	}

	if len(parseErrors) > 0 {
		return nil, errors.Join(parseErrors...)
	}
	return allOpts, nil
}

func parseTag(tag string) (tagOptions, []error) {
	tagLogger := loadLogger.With().Str("func", "parseTag").Str("tag", tag).Logger()
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
				opts.min = &val
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
