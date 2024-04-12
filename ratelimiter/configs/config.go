package configs

import "github.com/spf13/viper"

type conf struct {
	DBDriver             string `mapstructure:"REDIS_PORT"`
	DBHost               string `mapstructure:"REDIS_HOST"`
	DBPassword           string `mapstructure:"REDIS_PASSWORD"`
	RateLimiterRequests  string `mapstructure:"RATELIMITER_MAX_REQUESTS"`
	RateLimiterBlockTime string `mapstructure:"RATELIMITER_BLOCK_TIME"`
	RateLimiterTokens    string `mapstructure:"RATELIMITER_CUSTOM_TOKENS"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
