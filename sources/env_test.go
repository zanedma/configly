package sources

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	source := FromEnv()
	if source == nil {
		t.Fatal("expected FromEnv() to return non-nil source")
	}

	// Verify it returns an EnvSource
	_, ok := source.(*EnvSource)
	if !ok {
		t.Error("expected FromEnv() to return *EnvSource")
	}
}

func TestEnvSource_Name(t *testing.T) {
	source := &EnvSource{}
	name := source.Name()
	if name != "env" {
		t.Errorf("expected Name() to return 'env', got: %s", name)
	}
}

func TestEnvSource_GetValue(t *testing.T) {
	source := &EnvSource{}

	t.Run("get existing environment variable", func(t *testing.T) {
		// Set a test environment variable
		key := "TEST_ENV_VAR"
		expectedValue := "test-value"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)

		val, found, err := source.GetValue(key)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected variable to be found")
		}
		if val != expectedValue {
			t.Errorf("expected value to be '%s', got: %s", expectedValue, val)
		}
	})

	t.Run("get non-existent environment variable", func(t *testing.T) {
		key := "NONEXISTENT_ENV_VAR"
		// Make sure it doesn't exist
		os.Unsetenv(key)

		val, found, err := source.GetValue(key)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected variable not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("get empty environment variable", func(t *testing.T) {
		key := "EMPTY_ENV_VAR"
		os.Setenv(key, "")
		defer os.Unsetenv(key)

		val, found, err := source.GetValue(key)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected variable to be found even if empty")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("get environment variable with special characters", func(t *testing.T) {
		key := "SPECIAL_CHARS_VAR"
		expectedValue := "value with spaces and symbols: !@#$%"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)

		val, found, err := source.GetValue(key)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected variable to be found")
		}
		if val != expectedValue {
			t.Errorf("expected value to be '%s', got: %s", expectedValue, val)
		}
	})

	t.Run("get environment variable with newlines", func(t *testing.T) {
		key := "MULTILINE_VAR"
		expectedValue := "line1\nline2\nline3"
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)

		val, found, err := source.GetValue(key)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected variable to be found")
		}
		if val != expectedValue {
			t.Errorf("expected value to be '%s', got: %s", expectedValue, val)
		}
	})
}

func TestEnvSource_Integration(t *testing.T) {
	// Test the full workflow: FromEnv() -> Name() -> GetValue()
	source := FromEnv()

	// Set up test environment
	os.Setenv("INTEGRATION_TEST_HOST", "localhost")
	os.Setenv("INTEGRATION_TEST_PORT", "8080")
	os.Setenv("INTEGRATION_TEST_DEBUG", "true")
	defer func() {
		os.Unsetenv("INTEGRATION_TEST_HOST")
		os.Unsetenv("INTEGRATION_TEST_PORT")
		os.Unsetenv("INTEGRATION_TEST_DEBUG")
	}()

	// Test Name()
	if source.Name() != "env" {
		t.Errorf("expected source name to be 'env', got: %s", source.Name())
	}

	// Test GetValue() for each variable
	testCases := []struct {
		key           string
		expectedValue string
	}{
		{"INTEGRATION_TEST_HOST", "localhost"},
		{"INTEGRATION_TEST_PORT", "8080"},
		{"INTEGRATION_TEST_DEBUG", "true"},
	}

	for _, tc := range testCases {
		val, found, err := source.GetValue(tc.key)
		if err != nil {
			t.Errorf("expected no error for %s, got: %s", tc.key, err)
		}
		if !found {
			t.Errorf("expected %s to be found", tc.key)
		}
		if val != tc.expectedValue {
			t.Errorf("expected %s to be '%s', got: %s", tc.key, tc.expectedValue, val)
		}
	}

	// Test non-existent variable
	val, found, err := source.GetValue("NONEXISTENT_INTEGRATION_VAR")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	if found {
		t.Error("expected NONEXISTENT_INTEGRATION_VAR not to be found")
	}
	if val != "" {
		t.Errorf("expected empty value, got: %s", val)
	}
}
