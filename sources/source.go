package sources

// Source is an interface for retrieving configuration values.
type Source interface {
	// Name returns the name of the configuration source.
	Name() string
	// GetValue retrieves a single configuration value by key.
	// Returns the value, whether it was found, and any error that occurred.
	GetValue(key string) (val string, found bool, err error)
}
