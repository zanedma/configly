package sources

// MockSource is a mock configuration source for testing.
type MockSource struct {
	SourceName string            // Name of the source
	Values     map[string]string // Key-value pairs to return
	Err        error             // Error to return (if any)
}

func (m *MockSource) Name() string {
	return m.SourceName
}

func (m *MockSource) GetValue(key string) (string, bool, error) {
	if m.Err != nil {
		return "", false, m.Err
	}
	val, found := m.Values[key]
	return val, found, nil
}

func (m *MockSource) GetPartialConfig(keys []string) (map[string]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	result := make(map[string]string)
	for _, key := range keys {
		if val, found := m.Values[key]; found {
			result[key] = val
		}
	}
	return result, nil
}
