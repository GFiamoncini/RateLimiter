package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/GFiamoncini/RateLimiter/limiter"
)

type RateLimiterMiddleware struct {
	limiter       limiter.RateLimiter
	limitPerIP    int
	limitPerToken int
}

func NewRateLimiterMiddleware(l limiter.RateLimiter, limitPerIP, limitPerToken int) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter:       l,
		limitPerIP:    limitPerIP,
		limitPerToken: limitPerToken,
	}
}

func (m *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		apiKey := r.Header.Get("API_KEY")
		key := ""

		if apiKey != "" {
			key = "token:" + apiKey
			allowed, ttl := m.limiter.Allow(ctx, key, m.limitPerToken)
			if !allowed {
				http.Error(w, "Você atingiu o número máximo de solicitações para esse Token", http.StatusTooManyRequests)
				w.Header().Set("Tente novamente mais tarde ou troque sua autenticação.", ttl.String())
				return
			}
		} else {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			key = "ip:" + ip
			allowed, ttl := m.limiter.Allow(ctx, key, m.limitPerIP)
			if !allowed {
				http.Error(w, "Você atingiu o número máximo de Solicitações para este IP", http.StatusTooManyRequests)
				w.Header().Set("Tente novamente mais tarde...", ttl.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
