package sources

// MockSource is a mock configuration source for testing.
type MockSource struct {
	SourceName string            // Name of the source
	Values     map[string]string // Key-value pairs to return
	Err        error             // Error to return (if any)
}

// Name returns the name of this mock source.
func (m *MockSource) Name() string {
	return m.SourceName
}

// GetValue retrieves a value from the mock source.
func (m *MockSource) GetValue(key string) (string, bool, error) {
	if m.Err != nil {
		return "", false, m.Err
	}
	val, found := m.Values[key]
	return val, found, nil
}
