package configly

import (
	"testing"
)

func TestInvalidType(t *testing.T) {
	type StringConfig string

	_, err := Load[StringConfig](&EnvSource{})
	if err == nil {
		t.Error("error was nil when string specified as T")
	}

	if err != ErrInvalidType {
		t.Errorf("incorrect error received when string specified as T: %s\n", err)
	}
}

func TestInvalidTag(t *testing.T) {
	type Config struct {
		Value string `config:"test"`
	}

	Load[Config](&EnvSource{})
}
