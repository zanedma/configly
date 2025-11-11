package main

import (
	"fmt"
	"time"

	"github.com/zanedma/configly"
)

type Config struct {
	Port           int           `configly:"PORT,default=8080"`
	Database       string        `configly:"DB_URL,required"`
	Timeout        time.Duration `configly:"TIMEOUT,default=30s"`
	Min            int           `configly:"MIN,min=50"`
	Max            int           `configly:"MAX,max=50"`
	Invalid        int           `configly:"TEST,min=foo,minLen=bar"`
	AnotherInvalid int           `configly:"TEST,max=foor,maxLen=baz"`
}

func main() {
	loaderConfig := configly.LoaderConfig{
		Sources: []configly.Source{&configly.EnvSource{}},
	}
	loader, err := configly.New[Config](loaderConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(loader.Load())
}
