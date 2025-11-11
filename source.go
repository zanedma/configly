package configly

type Source interface {
	GetPartialConfig(keys []string) (map[string]string, error)
}
