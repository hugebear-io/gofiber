package fabric

import (
	"strings"

	"github.com/spf13/viper"
)

type DefaultConfig struct {
	Application ApplicationConfig `mapstructure:"app"`
	Logger      LoggerConfig      `mapstructure:"logger"`
}

type ApplicationConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type LoggerConfig struct {
	Level      int `mapstructure:"level"`
	SkipCaller int `mapstructure:"skipcaller"`
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func ReadConfigByKey(key string, cfg interface{}) error {
	if err := viper.UnmarshalKey(key, &cfg); err != nil {
		return err
	}
	return nil
}

func ReadConfig(cfg interface{}) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}
