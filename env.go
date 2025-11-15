package configly

import (
	"os"

	"github.com/zanedma/configly/sources"
)

type envSource struct{}

func FromEnv() sources.Source {
	return &envSource{}
}

func (s *envSource) Name() string {
	return "env"
}

func (s *envSource) GetValue(key string) (string, bool, error) {
	val, found := os.LookupEnv(key)
	return val, found, nil
}

func (s *envSource) GetPartialConfig(keys []string) (map[string]string, error) {
	config := map[string]string{}
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			config[key] = value
		}
	}
	return config, nil
}
