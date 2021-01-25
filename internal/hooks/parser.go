package hooks

import (
	"github.com/spf13/viper"
)

type Config struct {
	Hooks []Hook
}

var config Config

func Parse() {
	viper.Unmarshal(&config)
}

func Defined() bool {
	return len(config.Hooks) > 0
}
