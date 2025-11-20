package sources

import (
	"os"
	"strings"
)

// CLISource is a configuration source that reads from command-line arguments.
type CLISource struct {
	flags map[string]string
}

// FromCLI creates a new command-line argument configuration source.
// It parses command-line flags in the format -key=value or --key=value.
func FromCLI() Source {
	return FromCLIArgs(nil)
}

// FromCLIArgs creates a new command-line argument configuration source
// with explicit arguments. If args is nil, os.Args[1:] is used.
func FromCLIArgs(args []string) Source {
	s := &CLISource{
		flags: make(map[string]string),
	}

	if args == nil {
		args = os.Args[1:]
	}

	// Parse simple key=value pairs from command line
	for _, arg := range args {
		trimmed := strings.TrimPrefix(arg, "--")
		trimmed = strings.TrimPrefix(trimmed, "-")

		if parts := strings.SplitN(trimmed, "=", 2); len(parts) == 2 {
			s.flags[parts[0]] = parts[1]
		}
	}

	return s
}

// Name returns the name of this source.
func (s *CLISource) Name() string {
	return "cli"
}

// GetValue retrieves a command-line flag value by key.
func (s *CLISource) GetValue(key string) (string, bool, error) {
	val, found := s.flags[key]
	return val, found, nil
}
