package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	DBPort               int    `mapstructure:"REDIS_PORT"`
	DBHost               string `mapstructure:"REDIS_HOST"`
	DBPassword           string `mapstructure:"REDIS_PASSWORD"`
	DBDataTTL            int    `mapstructure:"REDIS_DATATTL"`
	RateLimiterRequests  int64  `mapstructure:"RATELIMITER_MAX_REQUESTS"`
	RateLimiterBlockTime int64  `mapstructure:"RATELIMITER_BLOCK_TIME"`
	RateLimiterTokens    string `mapstructure:"RATELIMITER_CUSTOM_TOKENS"`
	ServerPort           int    `mapstructure:"SERVER_PORT"`
}

func LoadConfig(path string) (*Conf, error) {
	var conf *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return conf, err
}
