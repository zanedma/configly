// Package configly provides a type-safe, validation-focused configuration loader
// for Go applications with support for multiple configuration sources.
//
// # Overview
//
// Configly uses Go generics to provide compile-time type safety for your configuration.
// It loads values from multiple sources (environment variables, JSON, YAML, .env files)
// with priority-based resolution, and validates constraints using struct tags.
//
// # Quick Start
//
// Define your configuration structure with validation tags:
//
//	type Config struct {
//	    Port     int           `configly:"PORT,default=8080"`
//	    Host     string        `configly:"HOST,default=localhost"`
//	    Database string        `configly:"DB_URL,required"`
//	    Timeout  time.Duration `configly:"TIMEOUT,default=30s"`
//	}
//
// Create a loader and load your configuration:
//
//	loader, err := configly.New[Config](configly.LoaderConfig{
//	    Sources: []sources.Source{
//	        configly.FromEnv(),
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	config, err := loader.Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Validation Tags
//
// Available struct tag options:
//   - required: Field must have a value
//   - default=VALUE: Default value if not found
//   - min=N: Minimum value for numbers
//   - max=N: Maximum value for numbers
//   - minLen=N: Minimum string length
//   - maxLen=N: Maximum string length
//
// # Multiple Sources
//
// Configure multiple sources with priority ordering (first source wins):
//
//	loader, err := configly.New[Config](configly.LoaderConfig{
//	    Sources: []sources.Source{
//	        configly.FromEnv(),              // Highest priority
//	        sources.FromFile(".env.local"),  // Second priority
//	        sources.FromFile("config.yaml"), // Lowest priority
//	    },
//	})
//
// # Supported Types
//
// Configly supports the following field types:
//   - string
//   - bool
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - time.Duration
//
// See the sources subpackage for available configuration sources including
// FromFile() for JSON, YAML, and .env files.
package configly
