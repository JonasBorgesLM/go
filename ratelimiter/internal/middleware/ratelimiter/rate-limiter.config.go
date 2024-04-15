package ratelimiter

import (
	"encoding/json"
	"fmt"

	"github.com/JonasBorgesLM/go/ratelimiter/configs"
	"github.com/JonasBorgesLM/go/ratelimiter/internal/infra/database"
	"github.com/JonasBorgesLM/go/ratelimiter/internal/infra/entity"
)

type RateLimitRule struct {
	MaxRequestsPerSecond  int64 `json:"max_requests"`
	BlockTimeMilliseconds int64 `json:"block_time"`
}

type RateLimiterConfig struct {
	Domain       *RateLimitRule            `json:"domain"`
	CustomTokens map[string]*RateLimitRule `json:"custom_tokens"`
	Repository   entity.RedisInterface     `json:"repository"`
}

func (rlc *RateLimiterConfig) GetRateLimitRuleForToken(token string) *RateLimitRule {
	customTokenRule, ok := rlc.CustomTokens[token]
	if !ok {
		return rlc.Domain
	}

	return customTokenRule
}

func GetRateLimiter(conf *configs.Conf) *RateLimiterConfig {
	customTokensEnvs := conf.RateLimiterTokens
	defaultMaxRequests := getDefaultMaxRequests(conf.RateLimiterRequests)
	defaultBlockTime := getDefaultBlockTime(conf.RateLimiterBlockTime)

	Domain := &RateLimitRule{
		MaxRequestsPerSecond:  defaultMaxRequests,
		BlockTimeMilliseconds: defaultBlockTime,
	}

	customTokens := getCustomTokens(customTokensEnvs, defaultMaxRequests, defaultBlockTime)

	addr := fmt.Sprintf("%s:%v", conf.DBHost, conf.DBPort)

	return &RateLimiterConfig{
		Domain:       Domain,
		CustomTokens: customTokens,
		Repository:   database.NewRedisRepository(addr, conf.DBPassword, conf.DBDataTTL),
	}
}

func getDefaultMaxRequests(maxRequests int64) int64 {
	if maxRequests == 0 {
		return 5
	}

	return maxRequests
}

func getDefaultBlockTime(blockTime int64) int64 {
	if blockTime == 0 {
		return 500
	}

	return blockTime
}

func getCustomTokens(customTokensEnvs string, defaultMaxRequests, defaultBlockTime int64) map[string]*RateLimitRule {
	var tokens []struct {
		Name                  string `json:"name"`
		MaxRequestsPerSecond  int64  `json:"max_requests_per_second"`
		BlockTimeMilliseconds int64  `json:"block_time_in_milliseconds"`
	}
	if customTokensEnvs == "" {
		customTokensEnvs = `[]`
	}
	err := json.Unmarshal([]byte(customTokensEnvs), &tokens)
	if err != nil {
		panic(err)
	}

	customTokens := make(map[string]*RateLimitRule)
	for _, tkn := range tokens {
		if tkn.MaxRequestsPerSecond == 0 {
			tkn.MaxRequestsPerSecond = defaultMaxRequests
		}

		if tkn.BlockTimeMilliseconds == 0 {
			tkn.BlockTimeMilliseconds = defaultBlockTime
		}

		customTokens[tkn.Name] = &RateLimitRule{
			MaxRequestsPerSecond:  tkn.MaxRequestsPerSecond,
			BlockTimeMilliseconds: tkn.BlockTimeMilliseconds,
		}
	}

	return customTokens
}
