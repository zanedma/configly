package sources

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type FileSource struct {
	kvMap    map[string]string
	filePath string
}

func FromFile(path string) (*FileSource, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	split := strings.Split(path, ".")
	if len(split) < 2 || split[len(split)-1] == "" {
		return nil, fmt.Errorf("file has no extension: %s", path)
	}

	switch split[len(split)-1] {
	case "json":
		kvMap, err := unmarshalFile(bytes, "json", json.Unmarshal)
		if err != nil {
			return nil, err
		}
		return &FileSource{
			kvMap:    kvMap,
			filePath: path,
		}, nil
	case "yml", "yaml":
		kvMap, err := unmarshalFile(bytes, "yaml", yaml.Unmarshal)
		if err != nil {
			return nil, err
		}
		return &FileSource{
			kvMap:    kvMap,
			filePath: path,
		}, nil
	}

	// Check if this is an env file: extension is "env" OR "env" appears in the middle
	// Examples: .env, .env.local, config.env
	ext := split[len(split)-1]
	isEnvFile := ext == "env" || slices.Contains(split[1:len(split)-1], "env")
	if !isEnvFile {
		return nil, errors.New("unsupported file type")
	}

	kvMap, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing env file: %w", err)
	}

	return &FileSource{
		kvMap:    kvMap,
		filePath: path,
	}, nil
}

func unmarshalFile(bytes []byte, fileType string, unmarshalFunc func(bytes []byte, out any) error) (map[string]string, error) {
	var out map[string]interface{}
	err := unmarshalFunc(bytes, &out)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s file: %w", fileType, err)
	}

	// Convert to map[string]string, filtering out non-scalar values
	result := make(map[string]string)
	for key, value := range out {
		// Only include scalar values (strings, numbers, booleans)
		// Objects, arrays, and null are ignored
		switch v := value.(type) {
		case string:
			result[key] = v
		case bool:
			result[key] = fmt.Sprintf("%t", v)
		case float64: // JSON numbers are float64
			result[key] = fmt.Sprintf("%v", v)
		case int, int8, int16, int32, int64:
			result[key] = fmt.Sprintf("%d", v)
		case uint, uint8, uint16, uint32, uint64:
			result[key] = fmt.Sprintf("%d", v)
		// Ignore: maps, slices, nil (objects, arrays, null)
		}
	}

	return result, nil
}

func (fs *FileSource) Name() string {
	return fmt.Sprintf("file:%s", fs.filePath)
}

func (fs *FileSource) GetValue(key string) (string, bool, error) {
	val, found := fs.kvMap[key]
	return val, found, nil
}
