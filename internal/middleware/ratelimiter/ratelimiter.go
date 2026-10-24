package ratelimiter

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type KeyFunc func(*http.Request) string

func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	if net.ParseIP(r.RemoteAddr) != nil {
		return r.RemoteAddr
	}

	return ""
}

func RateLimiter(limiter *Limiter, keyFunc KeyFunc) func(http.Handler) http.Handler {

	if keyFunc == nil {
		keyFunc = ClientIP
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			key := strings.TrimSpace(keyFunc(r))
			if key == "" {
				http.Error(
					w,
					"Unable to identify client",
					http.StatusBadRequest,
				)
				return
			}

			decision := limiter.Allow(key)

			w.Header().Set(
				"RateLimit-Limit",
				strconv.Itoa(limiter.burst),
			)
			w.Header().Set(
				"RateLimit-Remaining",
				strconv.Itoa(decision.Remaining),
			)

			if !decision.Allowed {
				retrySeconds := int(math.Ceil(
					decision.RetryAfter.Seconds(),
				))
				if retrySeconds < 1 {
					retrySeconds = 1
				}

				w.Header().Set(
					"Retry-After",
					strconv.Itoa(retrySeconds),
				)
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = fmt.Fprintln(
					w,
					`{"error":"rate limit exceeded"}`,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
