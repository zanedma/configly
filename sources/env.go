package sources

import (
	"os"
)

// EnvSource is a configuration source that reads from environment variables.
type EnvSource struct{}

// FromEnv creates a new environment variable configuration source.
func FromEnv() Source {
	return &EnvSource{}
}

// Name returns the name of this source.
func (s *EnvSource) Name() string {
	return "env"
}

// GetValue retrieves an environment variable by key.
func (s *EnvSource) GetValue(key string) (string, bool, error) {
	val, found := os.LookupEnv(key)
	return val, found, nil
}
