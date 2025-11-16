# Configly

[![Go Reference](https://pkg.go.dev/badge/github.com/zanedma/configly.svg)](https://pkg.go.dev/github.com/zanedma/configly)
[![Go Report Card](https://goreportcard.com/badge/github.com/zanedma/configly)](https://goreportcard.com/report/github.com/zanedma/configly)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/zanedma/configly)](https://github.com/zanedma/configly/releases)

A safe, validation-focused configuration loader for Go with support for multiple configuration sources.

## Features

- **Type-Safe**: Uses Go generics to provide compile-time type safety for your configuration
- **Multiple Sources**: Load configuration from environment variables, JSON, YAML, and .env files
- **Priority-Based**: Define source priority - first source with a value wins
- **Validation Built-In**: Comprehensive validation with `required`, `min`, `max`, `minLen`, `maxLen` constraints
- **Default Values**: Specify default values directly in struct tags
- **Graceful Handling**: Files with complex structures (nested objects/arrays) are supported; only scalar values are loaded
- **Time Duration Support**: Native support for `time.Duration` parsing
- **Detailed Errors**: Get all validation errors at once, not just the first failure

## Installation

```bash
go get github.com/zanedma/configly
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/zanedma/configly"
    "github.com/zanedma/configly/sources"
)

// Define your configuration structure
type AppConfig struct {
    Port     int           `configly:"PORT,default=8080"`
    Host     string        `configly:"HOST,default=localhost"`
    Database string        `configly:"DB_URL,required"`
    Timeout  time.Duration `configly:"TIMEOUT,default=30s"`
    Debug    bool          `configly:"DEBUG,default=false"`
}

func main() {
    // Create a loader with environment variables as the source
    loader, err := configly.New[AppConfig](configly.LoaderConfig{
        Sources: []sources.Source{
            sources.FromEnv(),
        },
    })
    if err != nil {
        panic(err)
    }

    // Load and validate configuration
    config, err := loader.Load()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Server starting on %s:%d\n", config.Host, config.Port)
}
```

## Supported Configuration Sources

### Environment Variables

Load configuration from environment variables:

```go
sources.FromEnv()
```

### JSON Files

Load configuration from JSON files (supports nested structures, only scalar values are used):

```go
source, err := sources.FromFile("config.json")
```

**Example `config.json`:**
```json
{
  "PORT": 3000,
  "HOST": "0.0.0.0",
  "DEBUG": true,
  "nested": {
    "ignored": "object values are not loaded"
  }
}
```

### YAML Files

Load configuration from YAML files (`.yaml` or `.yml`):

```go
source, err := sources.FromFile("config.yaml")
```

**Example `config.yaml`:**
```yaml
PORT: 3000
HOST: 0.0.0.0
DEBUG: true
TIMEOUT: 45s
```

### .env Files

Load configuration from dotenv files (`.env`, `.env.local`, etc.):

```go
source, err := sources.FromFile(".env")
```

**Example `.env`:**
```bash
# Application configuration
PORT=3000
HOST="0.0.0.0"
DEBUG=true
DATABASE_URL='postgres://localhost/mydb'

# Supports multiline values
PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----"
```

## Multiple Sources with Priority

Configure multiple sources with priority ordering (first source wins):

```go
loader, err := configly.New[Config](configly.LoaderConfig{
    Sources: []sources.Source{
        sources.FromEnv(),              // Highest priority
        sources.FromFile(".env.local"), // Second priority
        sources.FromFile("config.yaml"), // Lowest priority
    },
})
```

## Struct Tags & Validation

Configly uses struct tags to define configuration behavior:

### Tag Format

```go
`configly:"KEY,option1,option2=value"`
```

### Available Options

| Option | Description | Example |
|--------|-------------|---------|
| `required` | Field must have a value | `configly:"API_KEY,required"` |
| `default=VALUE` | Default value if not found | `configly:"PORT,default=8080"` |
| `min=N` | Minimum value (numbers) | `configly:"PORT,min=1024"` |
| `max=N` | Maximum value (numbers) | `configly:"PORT,max=65535"` |
| `minLen=N` | Minimum length (strings) | `configly:"NAME,minLen=3"` |
| `maxLen=N` | Maximum length (strings) | `configly:"TOKEN,maxLen=256"` |

### Supported Types

- `string`
- `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `time.Duration`

### Validation Examples

```go
type Config struct {
    // Required field
    APIKey string `configly:"API_KEY,required"`

    // With default value
    Port int `configly:"PORT,default=8080"`

    // Numeric constraints
    MaxConnections int `configly:"MAX_CONN,min=1,max=1000,default=100"`

    // String length constraints
    Username string `configly:"USERNAME,required,minLen=3,maxLen=50"`

    // Time duration with default
    RequestTimeout time.Duration `configly:"TIMEOUT,default=30s"`

    // Boolean with default
    EnableMetrics bool `configly:"METRICS_ENABLED,default=true"`
}
```

## Custom Tag Keys

Use a custom struct tag key instead of `configly`:

```go
type Config struct {
    Port int `env:"PORT,default=8080"`
}

loader, err := configly.New[Config](configly.LoaderConfig{
    TagKey: "env",
    Sources: []sources.Source{sources.FromEnv()},
})
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/zanedma/configly"
    "github.com/zanedma/configly/sources"
)

type DatabaseConfig struct {
    Host            string        `configly:"DB_HOST,required"`
    Port            int           `configly:"DB_PORT,default=5432,min=1,max=65535"`
    Username        string        `configly:"DB_USER,required,minLen=1"`
    Password        string        `configly:"DB_PASS,required,minLen=8"`
    Database        string        `configly:"DB_NAME,required"`
    MaxConnections  int           `configly:"DB_MAX_CONN,default=25,min=1,max=100"`
    ConnectTimeout  time.Duration `configly:"DB_TIMEOUT,default=10s"`
}

type ServerConfig struct {
    Host           string        `configly:"SERVER_HOST,default=0.0.0.0"`
    Port           int           `configly:"SERVER_PORT,default=8080,min=1024"`
    ReadTimeout    time.Duration `configly:"READ_TIMEOUT,default=5s"`
    WriteTimeout   time.Duration `configly:"WRITE_TIMEOUT,default=10s"`
    MaxHeaderBytes int           `configly:"MAX_HEADER_BYTES,default=1048576"`
}

type AppConfig struct {
    Environment string         `configly:"ENV,default=development"`
    Debug       bool           `configly:"DEBUG,default=false"`
    Database    DatabaseConfig
    Server      ServerConfig
}

func main() {
    // Load from multiple sources with priority
    envFile, err := sources.FromFile(".env")
    if err != nil {
        log.Printf("Warning: couldn't load .env file: %v", err)
    }

    configFile, err := sources.FromFile("config.yaml")
    if err != nil {
        log.Printf("Warning: couldn't load config.yaml: %v", err)
    }

    sourcesToUse := []sources.Source{sources.FromEnv()}
    if envFile != nil {
        sourcesToUse = append(sourcesToUse, envFile)
    }
    if configFile != nil {
        sourcesToUse = append(sourcesToUse, configFile)
    }

    // Create loader
    loader, err := configly.New[AppConfig](configly.LoaderConfig{
        Sources: sourcesToUse,
    })
    if err != nil {
        log.Fatalf("Failed to create loader: %v", err)
    }

    // Load configuration
    config, err := loader.Load()
    if err != nil {
        log.Fatalf("Configuration validation failed: %v", err)
    }

    // Use configuration
    fmt.Printf("Starting %s server on %s:%d\n",
        config.Environment, config.Server.Host, config.Server.Port)
    fmt.Printf("Database: %s@%s:%d/%s\n",
        config.Database.Username, config.Database.Host,
        config.Database.Port, config.Database.Database)
}
```

## Error Handling

Configly returns all validation errors at once for better developer experience:

```go
config, err := loader.Load()
if err != nil {
    // Error message includes all validation failures:
    // "validation errors: field 'APIKey' is required but not found;
    //  field 'Port' value 99999 exceeds maximum 65535;
    //  field 'Username' length 2 is less than minimum 3"
    log.Fatal(err)
}
```

## Testing

Configly includes a `MockSource` for easy testing:

```go
import "github.com/zanedma/configly/sources"

func TestConfig(t *testing.T) {
    mockSource := &sources.MockSource{
        SourceName: "test",
        Values: map[string]string{
            "PORT": "8080",
            "HOST": "localhost",
        },
    }

    loader, _ := configly.New[Config](configly.LoaderConfig{
        Sources: []sources.Source{mockSource},
    })

    config, err := loader.Load()
    // ... assertions
}
```

## Advanced Features

### File Source Behavior

When loading from JSON or YAML files:
- Only scalar values (strings, numbers, booleans) are loaded
- Nested objects and arrays are gracefully ignored
- This allows you to maintain complex configuration files while only extracting the values you need

### Time Duration Parsing

Supports all Go duration formats:
- `"300ms"` → 300 milliseconds
- `"1.5s"` → 1.5 seconds
- `"2m"` → 2 minutes
- `"1h30m"` → 1 hour 30 minutes

### .env File Features

The .env file parser supports:
- Comments (lines starting with `#`)
- Empty lines (ignored)
- Single and double quoted values
- Multiline values (in quotes)
- Escaped characters (`\"`, `\\`, `\n`)
- Export prefix (e.g., `export VAR=value`)
- Whitespace around `=` is handled gracefully

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)

## Author

Built by [zanedma](https://github.com/zanedma)
