package ratelimiter

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/JonasBorgesLM/go/ratelimiter/configs"
	"github.com/JonasBorgesLM/go/ratelimiter/internal/infra/entity"
	"github.com/redis/go-redis/v9"
)

type RateLimiterCheckFunction func(ctx context.Context, key string, storage entity.RedisInterface, rateConfig *RateLimitRule) (*time.Time, error)

func NewRateLimiter(conf *configs.Conf) func(next http.Handler) http.Handler {
	rate := GetRateLimiter(conf)

	return func(next http.Handler) http.Handler {
		return rateLimiter(rate, next, checkRateLimit)
	}
}

func rateLimiter(config *RateLimiterConfig, next http.Handler, checkFunc RateLimiterCheckFunction) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var blockTime *time.Time
		var err error

		token := r.Header.Get("API_KEY")
		if token != "" {
			tokenConfig := config.GetRateLimitRuleForToken(token)
			_, _ = checkFunc(r.Context(), token, config.Repository, tokenConfig)
		}

		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		blockTime, err = checkFunc(r.Context(), host, config.Repository, config.Domain)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if blockTime != nil {
			http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkRateLimit(
	ctx context.Context,
	key string,
	storage entity.RedisInterface,
	rateConfig *RateLimitRule,
) (*time.Time, error) {
	if key == "" {
		return nil, nil
	}

	block, err := storage.GetBlock(ctx, key)
	if err != nil {
		return nil, handleRateLimitError(err)
	}

	if shouldBlockAccess(block) {
		return block, nil
	}

	success, err := storage.IncrementAccessCount(ctx, key, rateConfig.MaxRequestsPerSecond)
	if err != nil {
		return nil, handleRateLimitError(err)
	}

	if !success {
		block, err := storage.AddBlock(ctx, key, rateConfig.BlockTimeMilliseconds)
		if err != nil {
			return nil, handleRateLimitError(err)
		}
		return block, nil
	}

	return nil, nil
}

func shouldBlockAccess(block *time.Time) bool {
	return block != nil
}

func handleRateLimitError(err error) error {
	if err == redis.Nil {
		return nil
	}

	return err
}
