package configly

import (
	"os"
)

type EnvSource struct{}

func FromEnv() Source {
	return &EnvSource{}
}

func (s *EnvSource) Name() string {
	return "env"
}

func (s *EnvSource) GetValue(key string) (string, bool, error) {
	val, found := os.LookupEnv(key)
	return val, found, nil
}

func (s *EnvSource) GetPartialConfig(keys []string) (map[string]string, error) {
	config := map[string]string{}
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			config[key] = value
		}
	}
	return config, nil
}
