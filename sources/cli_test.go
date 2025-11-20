package sources

import (
	"os"
	"testing"
)

func TestFromCLI(t *testing.T) {
	source := FromCLI()
	if source == nil {
		t.Fatal("expected FromCLI() to return non-nil source")
	}

	// Verify it returns a CLISource
	_, ok := source.(*CLISource)
	if !ok {
		t.Error("expected FromCLI() to return *CLISource")
	}
}

func TestFromCLIArgs(t *testing.T) {
	t.Run("with nil args uses os.Args", func(t *testing.T) {
		// Save original os.Args and restore after test
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Set os.Args to simulate command-line arguments
		os.Args = []string{"program", "--test=value", "--foo=bar"}

		source := FromCLIArgs(nil)
		if source == nil {
			t.Fatal("expected FromCLIArgs(nil) to return non-nil source")
		}

		// Verify it returns a CLISource
		cliSource, ok := source.(*CLISource)
		if !ok {
			t.Error("expected FromCLIArgs(nil) to return *CLISource")
		}

		// Verify it actually parsed os.Args[1:]
		val, found, err := cliSource.GetValue("test")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'test' flag to be found from os.Args")
		}
		if val != "value" {
			t.Errorf("expected 'test' to be 'value', got: %s", val)
		}

		val, found, err = cliSource.GetValue("foo")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'foo' flag to be found from os.Args")
		}
		if val != "bar" {
			t.Errorf("expected 'foo' to be 'bar', got: %s", val)
		}
	})

	t.Run("with explicit args", func(t *testing.T) {
		args := []string{"--host=localhost", "--port=8080"}
		source := FromCLIArgs(args)
		if source == nil {
			t.Fatal("expected FromCLIArgs(args) to return non-nil source")
		}

		// Verify it returns a CLISource
		_, ok := source.(*CLISource)
		if !ok {
			t.Error("expected FromCLIArgs(args) to return *CLISource")
		}
	})
}

func TestCLISource_Name(t *testing.T) {
	source := &CLISource{
		flags: make(map[string]string),
	}
	name := source.Name()
	if name != "cli" {
		t.Errorf("expected Name() to return 'cli', got: %s", name)
	}
}

func TestCLISource_GetValue(t *testing.T) {
	t.Run("parse double dash flags", func(t *testing.T) {
		args := []string{"--host=localhost", "--port=8080", "--debug=true"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'host' to be found")
		}
		if val != "localhost" {
			t.Errorf("expected value to be 'localhost', got: %s", val)
		}

		val, found, err = source.GetValue("port")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'port' to be found")
		}
		if val != "8080" {
			t.Errorf("expected value to be '8080', got: %s", val)
		}

		val, found, err = source.GetValue("debug")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'debug' to be found")
		}
		if val != "true" {
			t.Errorf("expected value to be 'true', got: %s", val)
		}
	})

	t.Run("parse single dash flags", func(t *testing.T) {
		args := []string{"-host=localhost", "-port=8080"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'host' to be found")
		}
		if val != "localhost" {
			t.Errorf("expected value to be 'localhost', got: %s", val)
		}

		val, found, err = source.GetValue("port")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'port' to be found")
		}
		if val != "8080" {
			t.Errorf("expected value to be '8080', got: %s", val)
		}
	})

	t.Run("get non-existent flag", func(t *testing.T) {
		args := []string{"--host=localhost"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("nonexistent")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected flag not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("parse empty value", func(t *testing.T) {
		args := []string{"--empty="}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("empty")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected flag to be found even if empty")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("parse value with spaces", func(t *testing.T) {
		args := []string{"--message=hello world"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("message")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'message' to be found")
		}
		if val != "hello world" {
			t.Errorf("expected value to be 'hello world', got: %s", val)
		}
	})

	t.Run("parse value with special characters", func(t *testing.T) {
		args := []string{"--url=https://example.com:8080/path?query=value&foo=bar"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("url")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'url' to be found")
		}
		expectedValue := "https://example.com:8080/path?query=value&foo=bar"
		if val != expectedValue {
			t.Errorf("expected value to be '%s', got: %s", expectedValue, val)
		}
	})

	t.Run("parse value with equals sign", func(t *testing.T) {
		args := []string{"--config=key=value"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("config")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'config' to be found")
		}
		if val != "key=value" {
			t.Errorf("expected value to be 'key=value', got: %s", val)
		}
	})

	t.Run("ignore flags without equals sign", func(t *testing.T) {
		args := []string{"--host=localhost", "--verbose", "--port=8080"}
		source := FromCLIArgs(args)

		// Should parse host and port
		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'host' to be found")
		}
		if val != "localhost" {
			t.Errorf("expected value to be 'localhost', got: %s", val)
		}

		// Should ignore verbose (no equals sign)
		val, found, err = source.GetValue("verbose")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected 'verbose' not to be found (no equals sign)")
		}
	})

	t.Run("mixed single and double dash", func(t *testing.T) {
		args := []string{"-host=localhost", "--port=8080", "-debug=true"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'host' to be found")
		}
		if val != "localhost" {
			t.Errorf("expected value to be 'localhost', got: %s", val)
		}

		val, found, err = source.GetValue("port")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'port' to be found")
		}
		if val != "8080" {
			t.Errorf("expected value to be '8080', got: %s", val)
		}

		val, found, err = source.GetValue("debug")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'debug' to be found")
		}
		if val != "true" {
			t.Errorf("expected value to be 'true', got: %s", val)
		}
	})

	t.Run("last value wins for duplicate keys", func(t *testing.T) {
		args := []string{"--host=localhost", "--host=example.com"}
		source := FromCLIArgs(args)

		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected 'host' to be found")
		}
		if val != "example.com" {
			t.Errorf("expected last value 'example.com', got: %s", val)
		}
	})
}

func TestCLISource_Integration(t *testing.T) {
	// Test the full workflow: FromCLIArgs() -> Name() -> GetValue()
	args := []string{
		"--host=localhost",
		"--port=8080",
		"--debug=true",
	}
	source := FromCLIArgs(args)

	// Test Name()
	if source.Name() != "cli" {
		t.Errorf("expected source name to be 'cli', got: %s", source.Name())
	}

	// Test GetValue() for each flag
	testCases := []struct {
		key           string
		expectedValue string
	}{
		{"host", "localhost"},
		{"port", "8080"},
		{"debug", "true"},
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

	// Test non-existent flag
	val, found, err := source.GetValue("nonexistent")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	if found {
		t.Error("expected 'nonexistent' not to be found")
	}
	if val != "" {
		t.Errorf("expected empty value, got: %s", val)
	}
}
