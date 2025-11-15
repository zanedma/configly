package sources

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromFile(t *testing.T) {
	t.Run("create source from JSON file", func(t *testing.T) {
		tmpDir := t.TempDir()
		jsonFile := filepath.Join(tmpDir, "config.json")
		content := `{"host": "localhost", "port": "8080"}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if source == nil {
			t.Fatal("expected source to be non-nil")
		}
	})

	t.Run("create source from YAML file", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlFile := filepath.Join(tmpDir, "config.yaml")
		content := `host: localhost
port: 8080`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if source == nil {
			t.Fatal("expected source to be non-nil")
		}
	})

	t.Run("create source from YML file", func(t *testing.T) {
		tmpDir := t.TempDir()
		ymlFile := filepath.Join(tmpDir, "config.yml")
		content := `host: localhost`
		if err := os.WriteFile(ymlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(ymlFile)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if source == nil {
			t.Fatal("expected source to be non-nil")
		}
	})

	t.Run("create source from .env file", func(t *testing.T) {
		tmpDir := t.TempDir()
		envFile := filepath.Join(tmpDir, ".env")
		content := `HOST=localhost
PORT=8080`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if source == nil {
			t.Fatal("expected source to be non-nil")
		}
	})

	t.Run("create source from .env.local file", func(t *testing.T) {
		tmpDir := t.TempDir()
		envFile := filepath.Join(tmpDir, ".env.local")
		content := `HOST=localhost`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if source == nil {
			t.Fatal("expected source to be non-nil")
		}
	})

	t.Run("create source from env file without dot", func(t *testing.T) {
		tmpDir := t.TempDir()
		envFile := filepath.Join(tmpDir, "config.env")
		content := `HOST=localhost`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if source == nil {
			t.Fatal("expected source to be non-nil")
		}
	})

	t.Run("error when file does not exist", func(t *testing.T) {
		source, err := FromFile("/nonexistent/file.json")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
		if source != nil {
			t.Error("expected source to be nil on error")
		}
	})

	t.Run("error when file extension is unsupported", func(t *testing.T) {
		tmpDir := t.TempDir()
		txtFile := filepath.Join(tmpDir, "config.txt")
		if err := os.WriteFile(txtFile, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(txtFile)
		if err == nil {
			t.Error("expected error for unsupported file extension")
		}
		if source != nil {
			t.Error("expected source to be nil on error")
		}
	})

	t.Run("error when JSON is invalid", func(t *testing.T) {
		tmpDir := t.TempDir()
		jsonFile := filepath.Join(tmpDir, "invalid.json")
		content := `{invalid json}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
		if source != nil {
			t.Error("expected source to be nil on error")
		}
	})

	t.Run("error when YAML is invalid", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlFile := filepath.Join(tmpDir, "invalid.yaml")
		content := "invalid:\n  - yaml\n - bad indent"
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err == nil {
			t.Error("expected error for invalid YAML")
		}
		if source != nil {
			t.Error("expected source to be nil on error")
		}
	})

	t.Run("allow JSON with nested objects and arrays", func(t *testing.T) {
		tmpDir := t.TempDir()
		jsonFile := filepath.Join(tmpDir, "complex.json")
		content := `{
			"host": "localhost",
			"database": {"host": "db.local"},
			"servers": ["server1", "server2"]
		}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		// Should not error - file is valid JSON
		source, err := FromFile(jsonFile)
		if err != nil {
			t.Errorf("expected no error for complex JSON, got: %s", err)
		}
		if source == nil {
			t.Error("expected source to be non-nil")
		}
	})

	t.Run("allow YAML with nested objects and arrays", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlFile := filepath.Join(tmpDir, "complex.yaml")
		content := `host: localhost
database:
  host: db.local
servers:
  - server1
  - server2`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %s", err)
		}

		// Should not error - file is valid YAML
		source, err := FromFile(yamlFile)
		if err != nil {
			t.Errorf("expected no error for complex YAML, got: %s", err)
		}
		if source == nil {
			t.Error("expected source to be non-nil")
		}
	})

}

func TestFileSource_Name(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")
	content := `{"host": "localhost"}`
	if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %s", err)
	}

	source, err := FromFile(jsonFile)
	if err != nil {
		t.Fatalf("failed to create source: %s", err)
	}

	name := source.Name()
	if name == "" {
		t.Error("expected Name() to return non-empty string")
	}
	// Name should be "file:<path>"
	expectedName := "file:" + jsonFile
	if name != expectedName {
		t.Errorf("expected Name() to return '%s', got: %s", expectedName, name)
	}
}

func TestFileSource_GetValue_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "config.json")

	t.Run("get existing scalar values", func(t *testing.T) {
		content := `{
			"host": "localhost",
			"port": "8080",
			"debug": "true",
			"app_name": "MyApp"
		}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"host", "localhost"},
			{"port", "8080"},
			{"debug", "true"},
			{"app_name", "MyApp"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("get non-existent key", func(t *testing.T) {
		content := `{"host": "localhost"}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		val, found, err := source.GetValue("nonexistent")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected key not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("handle different scalar types as strings", func(t *testing.T) {
		content := `{
			"stringVal": "hello",
			"intVal": 42,
			"floatVal": 3.14,
			"boolVal": true
		}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"stringVal", "hello"},
			{"intVal", "42"},
			{"floatVal", "3.14"},
			{"boolVal", "true"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("objects and arrays are not found", func(t *testing.T) {
		content := `{
			"host": "localhost",
			"database": {"host": "db.local", "port": 5432},
			"servers": ["server1", "server2", "server3"]
		}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// Scalar value should be found
		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected host to be found")
		}
		if val != "localhost" {
			t.Errorf("expected 'localhost', got: %s", val)
		}

		// Object should not be found
		val, found, err = source.GetValue("database")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected object 'database' not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value for object, got: %s", val)
		}

		// Array should not be found
		val, found, err = source.GetValue("servers")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected array 'servers' not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value for array, got: %s", val)
		}
	})

	t.Run("handle null values", func(t *testing.T) {
		content := `{"nullVal": null, "emptyVal": ""}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// null values should not be found
		val, found, err := source.GetValue("nullVal")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected nullVal not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value for null, got: %s", val)
		}

		// empty string values should be found
		val, found, err = source.GetValue("emptyVal")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected emptyVal to be found")
		}
		if val != "" {
			t.Errorf("expected empty string, got: %s", val)
		}
	})

	t.Run("handle empty JSON object", func(t *testing.T) {
		content := `{}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		val, found, err := source.GetValue("anykey")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected key not to be found in empty file")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})
}

func TestFileSource_GetValue_YAML(t *testing.T) {
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")

	t.Run("get existing scalar values", func(t *testing.T) {
		content := `host: localhost
port: 8080
debug: true
app_name: MyApp`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"host", "localhost"},
			{"port", "8080"},
			{"debug", "true"},
			{"app_name", "MyApp"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("get non-existent key", func(t *testing.T) {
		content := `host: localhost`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		val, found, err := source.GetValue("nonexistent")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected key not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("handle different scalar types as strings", func(t *testing.T) {
		content := `stringVal: hello
intVal: 42
floatVal: 3.14
boolVal: true`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"stringVal", "hello"},
			{"intVal", "42"},
			{"floatVal", "3.14"},
			{"boolVal", "true"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("objects and arrays are not found", func(t *testing.T) {
		content := `host: localhost
database:
  host: db.local
  port: 5432
servers:
  - server1
  - server2`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// Scalar value should be found
		val, found, err := source.GetValue("host")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected host to be found")
		}
		if val != "localhost" {
			t.Errorf("expected 'localhost', got: %s", val)
		}

		// Object should not be found
		val, found, err = source.GetValue("database")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected object 'database' not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value for object, got: %s", val)
		}

		// Array should not be found
		val, found, err = source.GetValue("servers")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected array 'servers' not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value for array, got: %s", val)
		}
	})

	t.Run("handle null and empty values", func(t *testing.T) {
		content := `nullVal: null
emptyVal: ""`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// null values should not be found
		val, found, err := source.GetValue("nullVal")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected nullVal not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value for null, got: %s", val)
		}

		// empty string values should be found
		val, found, err = source.GetValue("emptyVal")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected emptyVal to be found")
		}
		if val != "" {
			t.Errorf("expected empty string, got: %s", val)
		}
	})

	t.Run("handle empty YAML file", func(t *testing.T) {
		content := ``
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		val, found, err := source.GetValue("anykey")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected key not to be found in empty file")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})
}

func TestFileSource_GetValue_Env(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	t.Run("get existing values", func(t *testing.T) {
		content := `HOST=localhost
PORT=8080
DEBUG=true
APP_NAME=MyApp`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"HOST", "localhost"},
			{"PORT", "8080"},
			{"DEBUG", "true"},
			{"APP_NAME", "MyApp"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("handle quoted values", func(t *testing.T) {
		content := `DOUBLE_QUOTED="value with spaces"
SINGLE_QUOTED='another value'
MIXED_QUOTES="value with 'quotes' inside"
EMPTY_QUOTED=""`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"DOUBLE_QUOTED", "value with spaces"},
			{"SINGLE_QUOTED", "another value"},
			{"MIXED_QUOTES", "value with 'quotes' inside"},
			{"EMPTY_QUOTED", ""},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("handle values with special characters", func(t *testing.T) {
		content := `URL=https://example.com:8080/path?query=value
PASSWORD=p@ssw0rd!#$%
PATH=/usr/local/bin:/usr/bin:/bin
CONNECTION_STRING=host=localhost;port=5432;db=mydb`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"URL", "https://example.com:8080/path?query=value"},
			{"PASSWORD", "p@ssw0rd!#$%"},
			{"PATH", "/usr/local/bin:/usr/bin:/bin"},
			{"CONNECTION_STRING", "host=localhost;port=5432;db=mydb"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("ignore comments and empty lines", func(t *testing.T) {
		content := `# This is a comment
HOST=localhost

# Another comment
PORT=8080
  # Indented comment

DEBUG=true`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"HOST", "localhost"},
			{"PORT", "8080"},
			{"DEBUG", "true"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("handle whitespace around equals", func(t *testing.T) {
		content := `NO_SPACE=value
SPACE_BEFORE =value
SPACE_AFTER= value
SPACE_BOTH = value
LOTS_OF_SPACE   =   value`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// All should have trimmed whitespace
		testCases := []struct {
			key      string
			expected string
		}{
			{"NO_SPACE", "value"},
			{"SPACE_BEFORE", "value"},
			{"SPACE_AFTER", "value"},
			{"SPACE_BOTH", "value"},
			{"LOTS_OF_SPACE", "value"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("handle empty values", func(t *testing.T) {
		content := `EMPTY=
EMPTY_WITH_SPACE=
QUOTED_EMPTY=""`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []string{"EMPTY", "EMPTY_WITH_SPACE", "QUOTED_EMPTY"}
		for _, key := range testCases {
			val, found, err := source.GetValue(key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", key)
			}
			if val != "" {
				t.Errorf("expected %s to be empty string, got: %s", key, val)
			}
		}
	})

	t.Run("get non-existent key", func(t *testing.T) {
		content := `HOST=localhost`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		val, found, err := source.GetValue("NONEXISTENT")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected key not to be found")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("handle multiline values", func(t *testing.T) {
		content := `SINGLE_LINE=value
MULTILINE="line1
line2
line3"
MULTILINE_SINGLE='line1
line2'`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"SINGLE_LINE", "value"},
			{"MULTILINE", "line1\nline2\nline3"},
			{"MULTILINE_SINGLE", "line1\nline2"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("handle escaped characters in quoted strings", func(t *testing.T) {
		content := `ESCAPED_QUOTE="value with \"escaped\" quotes"
ESCAPED_NEWLINE="line1\nline2"
ESCAPED_BACKSLASH="path\\to\\file"`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"ESCAPED_QUOTE", `value with "escaped" quotes`},
			{"ESCAPED_NEWLINE", "line1\nline2"},
			{"ESCAPED_BACKSLASH", `path\to\file`},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})

	t.Run("handle empty file", func(t *testing.T) {
		content := ``
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		val, found, err := source.GetValue("anykey")
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected key not to be found in empty file")
		}
		if val != "" {
			t.Errorf("expected empty value, got: %s", val)
		}
	})

	t.Run("handle export prefix", func(t *testing.T) {
		content := `export HOST=localhost
export PORT=8080
REGULAR=value`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		testCases := []struct {
			key      string
			expected string
		}{
			{"HOST", "localhost"},
			{"PORT", "8080"},
			{"REGULAR", "value"},
		}

		for _, tc := range testCases {
			val, found, err := source.GetValue(tc.key)
			if err != nil {
				t.Errorf("expected no error for key %s, got: %s", tc.key, err)
			}
			if !found {
				t.Errorf("expected %s to be found", tc.key)
			}
			if val != tc.expected {
				t.Errorf("expected %s='%s', got: %s", tc.key, tc.expected, val)
			}
		}
	})
}

func TestFileSource_Integration(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("full workflow with JSON", func(t *testing.T) {
		jsonFile := filepath.Join(tmpDir, "app.json")
		content := `{
			"APP_NAME": "MyApp",
			"VERSION": "1.0.0",
			"HOST": "0.0.0.0",
			"PORT": "3000"
		}`
		if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(jsonFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// Test Name()
		expectedName := "file:" + jsonFile
		if source.Name() != expectedName {
			t.Errorf("expected name '%s', got: %s", expectedName, source.Name())
		}

		// Test GetValue() for various keys
		tests := []struct {
			key      string
			expected string
			found    bool
		}{
			{"APP_NAME", "MyApp", true},
			{"VERSION", "1.0.0", true},
			{"HOST", "0.0.0.0", true},
			{"PORT", "3000", true},
			{"NONEXISTENT", "", false},
		}

		for _, tt := range tests {
			val, found, err := source.GetValue(tt.key)
			if err != nil {
				t.Errorf("unexpected error for %s: %s", tt.key, err)
			}
			if found != tt.found {
				t.Errorf("expected found=%v for %s, got: %v", tt.found, tt.key, found)
			}
			if val != tt.expected {
				t.Errorf("expected %s='%s', got: %s", tt.key, tt.expected, val)
			}
		}
	})

	t.Run("full workflow with YAML", func(t *testing.T) {
		yamlFile := filepath.Join(tmpDir, "app.yaml")
		content := `APP_NAME: MyApp
VERSION: 1.0.0
HOST: 0.0.0.0
PORT: 3000`
		if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %s", err)
		}

		source, err := FromFile(yamlFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// Test Name()
		expectedName := "file:" + yamlFile
		if source.Name() != expectedName {
			t.Errorf("expected name '%s', got: %s", expectedName, source.Name())
		}

		// Test GetValue() for various keys
		tests := []struct {
			key      string
			expected string
			found    bool
		}{
			{"APP_NAME", "MyApp", true},
			{"VERSION", "1.0.0", true},
			{"HOST", "0.0.0.0", true},
			{"PORT", "3000", true},
			{"NONEXISTENT", "", false},
		}

		for _, tt := range tests {
			val, found, err := source.GetValue(tt.key)
			if err != nil {
				t.Errorf("unexpected error for %s: %s", tt.key, err)
			}
			if found != tt.found {
				t.Errorf("expected found=%v for %s, got: %v", tt.found, tt.key, found)
			}
			if val != tt.expected {
				t.Errorf("expected %s='%s', got: %s", tt.key, tt.expected, val)
			}
		}
	})

	t.Run("full workflow with .env file", func(t *testing.T) {
		tempDir := t.TempDir()
		// Create a temporary .env file
		envFile := filepath.Join(tempDir, "test.env")
		content := `# Application configuration
APP_NAME=MyApp
VERSION="1.0.0"
HOST='0.0.0.0'
PORT=3000
DESCRIPTION="A test application"
`
		if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test .env file: %s", err)
		}

		// Create source from file
		source, err := FromFile(envFile)
		if err != nil {
			t.Fatalf("failed to create source: %s", err)
		}

		// Test Name()
		expectedName := "file:" + envFile
		if source.Name() != expectedName {
			t.Errorf("expected name '%s', got: %s", expectedName, source.Name())
		}

		// Test GetValue() for various keys
		tests := []struct {
			key      string
			expected string
			found    bool
		}{
			{"APP_NAME", "MyApp", true},
			{"VERSION", "1.0.0", true},
			{"HOST", "0.0.0.0", true},
			{"PORT", "3000", true},
			{"DESCRIPTION", "A test application", true},
			{"NONEXISTENT", "", false},
		}

		for _, tt := range tests {
			val, found, err := source.GetValue(tt.key)
			if err != nil {
				t.Errorf("unexpected error for %s: %s", tt.key, err)
			}
			if found != tt.found {
				t.Errorf("expected found=%v for %s, got: %v", tt.found, tt.key, found)
			}
			if val != tt.expected {
				t.Errorf("expected %s='%s', got: %s", tt.key, tt.expected, val)
			}
		}
	})
}
