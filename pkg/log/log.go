package log

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	once sync.Once
	base zerolog.Logger
)

// TODO pass options
func GetBase() zerolog.Logger {
	once.Do(func() {
		zerolog.SetGlobalLevel(zerolog.GlobalLevel())
		base = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).
			With().
			Timestamp().
			Str("package", "configly").
			Logger().
			Level(zerolog.GlobalLevel())
	})
	return base
}
