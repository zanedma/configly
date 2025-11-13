package configly

type Source interface {
	Name() string
	GetValue(key string) (val string, found bool, err error)
	GetPartialConfig(keys []string) (map[string]string, error)
}
