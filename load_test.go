package configly

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zanedma/configly/sources"
)

type validConfig struct {
	Value string `configly:"value,required"`
}

type configWithDefaults struct {
	Host string `configly:"host,default=localhost"`
	Port int    `configly:"port,default=8080"`
}

type configWithValidation struct {
	Age      int    `configly:"age,min=0,max=120"`
	Username string `configly:"username,minLen=3,maxLen=20"`
}

type configWithAllTypes struct {
	StringVal   string        `configly:"string_val"`
	IntVal      int           `configly:"int_val"`
	Int8Val     int8          `configly:"int8_val"`
	Int16Val    int16         `configly:"int16_val"`
	Int32Val    int32         `configly:"int32_val"`
	Int64Val    int64         `configly:"int64_val"`
	UintVal     uint          `configly:"uint_val"`
	Uint8Val    uint8         `configly:"uint8_val"`
	Uint16Val   uint16        `configly:"uint16_val"`
	Uint32Val   uint32        `configly:"uint32_val"`
	Uint64Val   uint64        `configly:"uint64_val"`
	Float32Val  float32       `configly:"float32_val"`
	Float64Val  float64       `configly:"float64_val"`
	BoolVal     bool          `configly:"bool_val"`
	DurationVal time.Duration `configly:"duration_val"`
}

type configWithUnexported struct {
	Public  string `configly:"public"`
	private string `configly:"private"`
}

type configWithNoTags struct {
	Field1 string
	Field2 int
}

func TestNew(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		l, err := New[validConfig](LoaderConfig{
			Sources: []sources.Source{&sources.MockSource{SourceName: "test"}},
		})
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if l == nil {
			t.Error("expected returned loader to be non-nil")
		}
	})

	t.Run("non-struct generic", func(t *testing.T) {
		l, err := New[string](LoaderConfig{
			Sources: []sources.Source{&sources.MockSource{SourceName: "test"}},
		})
		if err == nil {
			t.Error("expected error to be non-nil")
		}
		if l != nil {
			t.Error("expected returned loader to be nil")
		}
	})

	t.Run("no sources", func(t *testing.T) {
		l, err := New[validConfig](LoaderConfig{})
		if err == nil {
			t.Error("expected error to be non-nil")
		}
		if l != nil {
			t.Error("expected returned loader to be nil")
		}
	})

	t.Run("custom tag key", func(t *testing.T) {
		type customTagConfig struct {
			Value string `env:"value"`
		}
		l, err := New[customTagConfig](LoaderConfig{
			TagKey:  "env",
			Sources: []sources.Source{&sources.MockSource{SourceName: "test"}},
		})
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if l == nil {
			t.Error("expected returned loader to be non-nil")
			t.FailNow()
		}
		if l.tagKey != "env" {
			t.Errorf("expected tagKey to be 'env', got: %s", l.tagKey)
		}
	})

	t.Run("default tag key", func(t *testing.T) {
		l, err := New[validConfig](LoaderConfig{
			Sources: []sources.Source{&sources.MockSource{SourceName: "test"}},
		})
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if l.tagKey != defaultTagKey {
			t.Errorf("expected tagKey to be '%s', got: %s", defaultTagKey, l.tagKey)
		}
	})

	t.Run("multiple sources", func(t *testing.T) {
		l, err := New[validConfig](LoaderConfig{
			Sources: []sources.Source{
				&sources.MockSource{SourceName: "source1"},
				&sources.MockSource{SourceName: "source2"},
			},
		})
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if len(l.sources) != 2 {
			t.Errorf("expected 2 sources, got: %d", len(l.sources))
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("load with required field present", func(t *testing.T) {
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"value": "test-value"},
		}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg == nil {
			t.Fatal("expected config to be non-nil")
		}
		if cfg.Value != "test-value" {
			t.Errorf("expected Value to be 'test-value', got: %s", cfg.Value)
		}
	})

	t.Run("load with required field missing", func(t *testing.T) {
		source := &sources.MockSource{SourceName: "test", Values: map[string]string{}}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err == nil {
			t.Error("expected error to be non-nil")
		}
		if cfg != nil {
			t.Error("expected config to be nil")
		}
	})

	t.Run("load with default values", func(t *testing.T) {
		source := &sources.MockSource{SourceName: "test", Values: map[string]string{}}
		l, _ := New[configWithDefaults](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg.Host != "localhost" {
			t.Errorf("expected Host to be 'localhost', got: %s", cfg.Host)
		}
		if cfg.Port != 8080 {
			t.Errorf("expected Port to be 8080, got: %d", cfg.Port)
		}
	})

	t.Run("load with source value overriding default", func(t *testing.T) {
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"host": "example.com"},
		}
		l, _ := New[configWithDefaults](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg.Host != "example.com" {
			t.Errorf("expected Host to be 'example.com', got: %s", cfg.Host)
		}
		if cfg.Port != 8080 {
			t.Errorf("expected Port to use default 8080, got: %d", cfg.Port)
		}
	})

	t.Run("load from multiple sources with priority", func(t *testing.T) {
		source1 := &sources.MockSource{
			SourceName: "source1",
			Values:     map[string]string{"value": "from-source1"},
		}
		source2 := &sources.MockSource{
			SourceName: "source2",
			Values:     map[string]string{"value": "from-source2"},
		}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source1, source2}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg.Value != "from-source1" {
			t.Errorf("expected Value from first source, got: %s", cfg.Value)
		}
	})

	t.Run("load skips source with error", func(t *testing.T) {
		source1 := &sources.MockSource{
			SourceName: "source1",
			Err:        errors.New("source error"),
		}
		source2 := &sources.MockSource{
			SourceName: "source2",
			Values:     map[string]string{"value": "from-source2"},
		}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source1, source2}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg.Value != "from-source2" {
			t.Errorf("expected Value from second source, got: %s", cfg.Value)
		}
	})

	t.Run("load with invalid tag format", func(t *testing.T) {
		type badConfig struct {
			Value int `configly:"value,min=abc"`
		}
		source := &sources.MockSource{SourceName: "test", Values: map[string]string{}}
		l, _ := New[badConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for invalid tag format")
		}
	})

	t.Run("load with unexported field", func(t *testing.T) {
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"public": "public-value", "private": "private-value"},
		}
		l, _ := New[configWithUnexported](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg.Public != "public-value" {
			t.Errorf("expected Public to be 'public-value', got: %s", cfg.Public)
		}
		// private field should remain empty as it's unexported
		if cfg.private != "" {
			t.Errorf("expected private to be empty, got: %s", cfg.private)
		}
	})

	t.Run("load with no tags", func(t *testing.T) {
		source := &sources.MockSource{SourceName: "test", Values: map[string]string{}}
		l, _ := New[configWithNoTags](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if cfg == nil {
			t.Error("expected config to be non-nil")
		}
	})
}

func TestSetField(t *testing.T) {
	t.Run("set all supported types", func(t *testing.T) {
		source := &sources.MockSource{
			SourceName: "test",
			Values: map[string]string{
				"string_val":   "hello",
				"int_val":      "42",
				"int8_val":     "8",
				"int16_val":    "16",
				"int32_val":    "32",
				"int64_val":    "64",
				"uint_val":     "42",
				"uint8_val":    "8",
				"uint16_val":   "16",
				"uint32_val":   "32",
				"uint64_val":   "64",
				"float32_val":  "3.14",
				"float64_val":  "2.71828",
				"bool_val":     "true",
				"duration_val": "5s",
			},
		}
		l, _ := New[configWithAllTypes](LoaderConfig{Sources: []sources.Source{source}})

		cfg, err := l.Load()
		if err != nil {
			t.Fatalf("expected err to be nil, got: %s", err)
		}

		if cfg.StringVal != "hello" {
			t.Errorf("expected StringVal to be 'hello', got: %s", cfg.StringVal)
		}
		if cfg.IntVal != 42 {
			t.Errorf("expected IntVal to be 42, got: %d", cfg.IntVal)
		}
		if cfg.Int8Val != 8 {
			t.Errorf("expected Int8Val to be 8, got: %d", cfg.Int8Val)
		}
		if cfg.Int16Val != 16 {
			t.Errorf("expected Int16Val to be 16, got: %d", cfg.Int16Val)
		}
		if cfg.Int32Val != 32 {
			t.Errorf("expected Int32Val to be 32, got: %d", cfg.Int32Val)
		}
		if cfg.Int64Val != 64 {
			t.Errorf("expected Int64Val to be 64, got: %d", cfg.Int64Val)
		}
		if cfg.UintVal != 42 {
			t.Errorf("expected UintVal to be 42, got: %d", cfg.UintVal)
		}
		if cfg.Float32Val != 3.14 {
			t.Errorf("expected Float32Val to be 3.14, got: %f", cfg.Float32Val)
		}
		if cfg.Float64Val != 2.71828 {
			t.Errorf("expected Float64Val to be 2.71828, got: %f", cfg.Float64Val)
		}
		if cfg.BoolVal != true {
			t.Errorf("expected BoolVal to be true, got: %v", cfg.BoolVal)
		}
		if cfg.DurationVal != 5*time.Second {
			t.Errorf("expected DurationVal to be 5s, got: %v", cfg.DurationVal)
		}
	})

	t.Run("set field with invalid int", func(t *testing.T) {
		type intConfig struct {
			Value int `configly:"value"`
		}
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"value": "not-a-number"},
		}
		l, _ := New[intConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for invalid int value")
		}
	})

	t.Run("set field with invalid uint", func(t *testing.T) {
		type uintConfig struct {
			Value uint `configly:"value"`
		}
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"value": "-1"},
		}
		l, _ := New[uintConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for invalid uint value")
		}
	})

	t.Run("set field with invalid float", func(t *testing.T) {
		type floatConfig struct {
			Value float64 `configly:"value"`
		}
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"value": "not-a-float"},
		}
		l, _ := New[floatConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for invalid float value")
		}
	})

	t.Run("set field with invalid bool", func(t *testing.T) {
		type boolConfig struct {
			Value bool `configly:"value"`
		}
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"value": "not-a-bool"},
		}
		l, _ := New[boolConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for invalid bool value")
		}
	})

	t.Run("set field with invalid duration", func(t *testing.T) {
		type durationConfig struct {
			Value time.Duration `configly:"value"`
		}
		source := &sources.MockSource{
			SourceName: "test",
			Values:     map[string]string{"value": "not-a-duration"},
		}
		l, _ := New[durationConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for invalid duration value")
		}
	})
}

func TestParseTag(t *testing.T) {
	l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{&sources.MockSource{SourceName: "test"}}})

	t.Run("parse simple key", func(t *testing.T) {
		opts, errs := l.parseTag("my_key")
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if opts.key != "my_key" {
			t.Errorf("expected key to be 'my_key', got: %s", opts.key)
		}
		if opts.required {
			t.Error("expected required to be false")
		}
	})

	t.Run("parse required tag", func(t *testing.T) {
		opts, errs := l.parseTag("my_key,required")
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if !opts.required {
			t.Error("expected required to be true")
		}
	})

	t.Run("parse default value", func(t *testing.T) {
		opts, errs := l.parseTag("my_key,default=hello")
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if opts.defaultValue != "hello" {
			t.Errorf("expected defaultValue to be 'hello', got: %s", opts.defaultValue)
		}
	})

	t.Run("parse min/max values", func(t *testing.T) {
		opts, errs := l.parseTag("my_key,min=0,max=100")
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if opts.min == nil || *opts.min != 0 {
			t.Error("expected min to be 0")
		}
		if opts.max == nil || *opts.max != 100 {
			t.Error("expected max to be 100")
		}
	})

	t.Run("parse minLen/maxLen values", func(t *testing.T) {
		opts, errs := l.parseTag("my_key,minLen=5,maxLen=50")
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if opts.minLen == nil || *opts.minLen != 5 {
			t.Error("expected minLen to be 5")
		}
		if opts.maxLen == nil || *opts.maxLen != 50 {
			t.Error("expected maxLen to be 50")
		}
	})

	t.Run("parse invalid min value", func(t *testing.T) {
		_, errs := l.parseTag("my_key,min=invalid")
		if len(errs) == 0 {
			t.Error("expected error for invalid min value")
		}
	})

	t.Run("parse invalid max value", func(t *testing.T) {
		_, errs := l.parseTag("my_key,max=invalid")
		if len(errs) == 0 {
			t.Error("expected error for invalid max value")
		}
	})

	t.Run("parse invalid minLen value", func(t *testing.T) {
		_, errs := l.parseTag("my_key,minLen=invalid")
		if len(errs) == 0 {
			t.Error("expected error for invalid minLen value")
		}
	})

	t.Run("parse invalid maxLen value", func(t *testing.T) {
		_, errs := l.parseTag("my_key,maxLen=invalid")
		if len(errs) == 0 {
			t.Error("expected error for invalid maxLen value")
		}
	})

	t.Run("parse complex tag with whitespace", func(t *testing.T) {
		opts, errs := l.parseTag("my_key, required, default=test, min=0, max=100")
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if !opts.required {
			t.Error("expected required to be true")
		}
		if opts.defaultValue != "test" {
			t.Errorf("expected defaultValue to be 'test', got: %s", opts.defaultValue)
		}
	})
}

func TestParseAllTags(t *testing.T) {
	t.Run("parse all valid tags", func(t *testing.T) {
		l, _ := New[configWithDefaults](LoaderConfig{Sources: []sources.Source{&sources.MockSource{SourceName: "test"}}})
		var cfg configWithDefaults
		val := reflect.ValueOf(&cfg).Elem()
		typ := val.Type()

		t.Log(typ.NumField(), val)
		tagMap, err := l.parseAllTags(typ.NumField(), val)
		t.Log(tagMap, err)
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}
		if tagMap == nil {
			t.Error("expected tagMap to be non-nil")
		}
	})

	t.Run("parse with mixed valid and invalid tags", func(t *testing.T) {
		type mixedConfig struct {
			Valid   string `configly:"valid"`
			Invalid int    `configly:"invalid,min=abc"`
		}
		var cfg mixedConfig
		val := reflect.ValueOf(&cfg).Elem()
		typ := val.Type()
		l, _ := New[mixedConfig](LoaderConfig{Sources: []sources.Source{&sources.MockSource{SourceName: "test"}}})

		_, err := l.parseAllTags(typ.NumField(), val)
		if err == nil {
			t.Error("expected error for invalid tag")
		}
	})
}

func TestGetValueFromSources(t *testing.T) {
	t.Run("get value from first source", func(t *testing.T) {
		source1 := &sources.MockSource{
			SourceName: "source1",
			Values:     map[string]string{"key": "value1"},
		}
		source2 := &sources.MockSource{
			SourceName: "source2",
			Values:     map[string]string{"key": "value2"},
		}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source1, source2}})

		val, sourceName, found := l.getValueFromSources("key")
		if !found {
			t.Error("expected value to be found")
		}
		if val != "value1" {
			t.Errorf("expected value to be 'value1', got: %s", val)
		}
		if sourceName != "source1" {
			t.Errorf("expected sourceName to be 'source1', got: %s", sourceName)
		}
	})

	t.Run("get value not found", func(t *testing.T) {
		source := &sources.MockSource{SourceName: "test", Values: map[string]string{}}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, _, found := l.getValueFromSources("nonexistent")
		if found {
			t.Error("expected value not to be found")
		}
	})

	t.Run("get value with source error", func(t *testing.T) {
		source := &sources.MockSource{
			SourceName: "test",
			Err:        errors.New("source error"),
		}
		l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, _, found := l.getValueFromSources("key")
		if found {
			t.Error("expected value not to be found when source has error")
		}
	})
}

func TestMultipleValidationErrors(t *testing.T) {
	t.Run("multiple required fields missing", func(t *testing.T) {
		type multiRequiredConfig struct {
			Field1 string `configly:"field1,required"`
			Field2 string `configly:"field2,required"`
			Field3 string `configly:"field3,required"`
		}
		source := &sources.MockSource{SourceName: "test", Values: map[string]string{}}
		l, _ := New[multiRequiredConfig](LoaderConfig{Sources: []sources.Source{source}})

		_, err := l.Load()
		if err == nil {
			t.Error("expected error for missing required fields")
		}
		// Check that error mentions multiple fields (using errors.Join)
		errStr := err.Error()
		if !contains(errStr, "field1") || !contains(errStr, "field2") || !contains(errStr, "field3") {
			t.Errorf("expected error to mention all missing fields, got: %s", errStr)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestValidateField(t *testing.T) {
	l, _ := New[validConfig](LoaderConfig{Sources: []sources.Source{&sources.MockSource{SourceName: "test"}}})

	t.Run("validate string with minLen", func(t *testing.T) {
		minLen := 5
		opts := tagOptions{
			key:    "test",
			minLen: &minLen,
		}

		// Valid string
		validStr := reflect.ValueOf("hello")
		err := l.validateField(validStr, opts)
		if err != nil {
			t.Errorf("expected no error for valid string, got: %s", err)
		}

		// Invalid string (too short)
		invalidStr := reflect.ValueOf("hi")
		err = l.validateField(invalidStr, opts)
		if err == nil {
			t.Error("expected error for string shorter than minLen")
		}
	})

	t.Run("validate string with maxLen", func(t *testing.T) {
		maxLen := 10
		opts := tagOptions{
			key:    "test",
			maxLen: &maxLen,
		}

		// Valid string
		validStr := reflect.ValueOf("hello")
		err := l.validateField(validStr, opts)
		if err != nil {
			t.Errorf("expected no error for valid string, got: %s", err)
		}

		// Invalid string (too long)
		invalidStr := reflect.ValueOf("this is a very long string")
		err = l.validateField(invalidStr, opts)
		if err == nil {
			t.Error("expected error for string longer than maxLen")
		}
	})

	t.Run("validate string with minLen and maxLen", func(t *testing.T) {
		minLen := 3
		maxLen := 20
		opts := tagOptions{
			key:    "test",
			minLen: &minLen,
			maxLen: &maxLen,
		}

		// Valid string
		validStr := reflect.ValueOf("username")
		err := l.validateField(validStr, opts)
		if err != nil {
			t.Errorf("expected no error for valid string, got: %s", err)
		}

		// Too short
		shortStr := reflect.ValueOf("ab")
		err = l.validateField(shortStr, opts)
		if err == nil {
			t.Error("expected error for string shorter than minLen")
		}

		// Too long
		longStr := reflect.ValueOf("this is way too long for the maximum length")
		err = l.validateField(longStr, opts)
		if err == nil {
			t.Error("expected error for string longer than maxLen")
		}
	})

	t.Run("validate int with min", func(t *testing.T) {
		min := int64(0)
		opts := tagOptions{
			key: "test",
			min: &min,
		}

		// Valid int
		validInt := reflect.ValueOf(42)
		err := l.validateField(validInt, opts)
		if err != nil {
			t.Errorf("expected no error for valid int, got: %s", err)
		}

		// Invalid int (below min)
		invalidInt := reflect.ValueOf(-5)
		err = l.validateField(invalidInt, opts)
		if err == nil {
			t.Error("expected error for int below minimum")
		}
	})

	t.Run("validate int with max", func(t *testing.T) {
		max := int64(100)
		opts := tagOptions{
			key: "test",
			max: &max,
		}

		// Valid int
		validInt := reflect.ValueOf(50)
		err := l.validateField(validInt, opts)
		if err != nil {
			t.Errorf("expected no error for valid int, got: %s", err)
		}

		// Invalid int (above max)
		invalidInt := reflect.ValueOf(150)
		err = l.validateField(invalidInt, opts)
		if err == nil {
			t.Error("expected error for int above maximum")
		}
	})

	t.Run("validate int with min and max", func(t *testing.T) {
		min := int64(0)
		max := int64(120)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		// Valid int
		validInt := reflect.ValueOf(25)
		err := l.validateField(validInt, opts)
		if err != nil {
			t.Errorf("expected no error for valid int, got: %s", err)
		}

		// Below min
		belowMin := reflect.ValueOf(-10)
		err = l.validateField(belowMin, opts)
		if err == nil {
			t.Error("expected error for int below minimum")
		}

		// Above max
		aboveMax := reflect.ValueOf(150)
		err = l.validateField(aboveMax, opts)
		if err == nil {
			t.Error("expected error for int above maximum")
		}
	})

	t.Run("validate int8 with min and max", func(t *testing.T) {
		min := int64(0)
		max := int64(100)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		validInt8 := reflect.ValueOf(int8(50))
		err := l.validateField(validInt8, opts)
		if err != nil {
			t.Errorf("expected no error for valid int8, got: %s", err)
		}
	})

	t.Run("validate int16 with min and max", func(t *testing.T) {
		min := int64(0)
		max := int64(1000)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		validInt16 := reflect.ValueOf(int16(500))
		err := l.validateField(validInt16, opts)
		if err != nil {
			t.Errorf("expected no error for valid int16, got: %s", err)
		}
	})

	t.Run("validate int32 with min and max", func(t *testing.T) {
		min := int64(0)
		max := int64(100000)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		validInt32 := reflect.ValueOf(int32(50000))
		err := l.validateField(validInt32, opts)
		if err != nil {
			t.Errorf("expected no error for valid int32, got: %s", err)
		}
	})

	t.Run("validate int64 with min and max", func(t *testing.T) {
		min := int64(0)
		max := int64(1000000)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		validInt64 := reflect.ValueOf(int64(500000))
		err := l.validateField(validInt64, opts)
		if err != nil {
			t.Errorf("expected no error for valid int64, got: %s", err)
		}
	})

	t.Run("validate uint with min and max", func(t *testing.T) {
		min := int64(10)
		max := int64(100)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		validUint := reflect.ValueOf(uint(50))
		err := l.validateField(validUint, opts)
		if err != nil {
			t.Errorf("expected no error for valid uint, got: %s", err)
		}
	})

	t.Run("validate float32 with min", func(t *testing.T) {
		min := int64(0)
		opts := tagOptions{
			key: "test",
			min: &min,
		}

		// Valid float
		validFloat := reflect.ValueOf(float32(3.14))
		err := l.validateField(validFloat, opts)
		if err != nil {
			t.Errorf("expected no error for valid float32, got: %s", err)
		}

		// Invalid float (below min)
		invalidFloat := reflect.ValueOf(float32(-1.5))
		err = l.validateField(invalidFloat, opts)
		if err == nil {
			t.Error("expected error for float32 below minimum")
		}
	})

	t.Run("validate float64 with max", func(t *testing.T) {
		max := int64(100)
		opts := tagOptions{
			key: "test",
			max: &max,
		}

		// Valid float
		validFloat := reflect.ValueOf(float64(50.5))
		err := l.validateField(validFloat, opts)
		if err != nil {
			t.Errorf("expected no error for valid float64, got: %s", err)
		}

		// Invalid float (above max)
		invalidFloat := reflect.ValueOf(float64(150.5))
		err = l.validateField(invalidFloat, opts)
		if err == nil {
			t.Error("expected error for float64 above maximum")
		}
	})

	t.Run("validate float with min and max", func(t *testing.T) {
		min := int64(0)
		max := int64(100)
		opts := tagOptions{
			key: "test",
			min: &min,
			max: &max,
		}

		// Valid float
		validFloat := reflect.ValueOf(float64(50.5))
		err := l.validateField(validFloat, opts)
		if err != nil {
			t.Errorf("expected no error for valid float64, got: %s", err)
		}

		// Below min
		belowMin := reflect.ValueOf(float64(-10.5))
		err = l.validateField(belowMin, opts)
		if err == nil {
			t.Error("expected error for float64 below minimum")
		}

		// Above max
		aboveMax := reflect.ValueOf(float64(150.5))
		err = l.validateField(aboveMax, opts)
		if err == nil {
			t.Error("expected error for float64 above maximum")
		}
	})

	t.Run("validate field with no constraints", func(t *testing.T) {
		opts := tagOptions{
			key: "test",
		}

		// Should pass for any value
		err := l.validateField(reflect.ValueOf("any string"), opts)
		if err != nil {
			t.Errorf("expected no error for unconstrained field, got: %s", err)
		}

		err = l.validateField(reflect.ValueOf(12345), opts)
		if err != nil {
			t.Errorf("expected no error for unconstrained field, got: %s", err)
		}
	})

	t.Run("validate bool field", func(t *testing.T) {
		opts := tagOptions{
			key: "test",
		}

		// Bool fields should not error (no validation for bool)
		err := l.validateField(reflect.ValueOf(true), opts)
		if err != nil {
			t.Errorf("expected no error for bool field, got: %s", err)
		}
	})
}
