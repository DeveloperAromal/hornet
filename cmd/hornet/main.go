package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DeveloperAromal/hornet/internal/auth"
	banner "github.com/DeveloperAromal/hornet/internal/cli"
	"github.com/DeveloperAromal/hornet/internal/middleware"
	"github.com/DeveloperAromal/hornet/internal/proxy"
	"github.com/DeveloperAromal/hornet/internal/ratelimiter"
	"github.com/DeveloperAromal/hornet/internal/requestguard"
)

func main() {
	banner.Banner()

	limiter, err := ratelimiter.NewLimiter(ratelimiter.Config{
		RatePerSecond: 10,
		Burst:         20,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer limiter.Close()

	verifier, err := auth.NewJwtVerifier(
		os.Getenv("HORNET_JWT_SECRET"),
		"hornet-auth",
	)
	if err != nil {
		log.Fatal(err)
	}

	p := proxy.NewProxy()

	rp, err := p.NewReverseProxy("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	guard := requestguard.Limits{
		MaxURLLength: 8 * 1024,
		MaxBodySize:  10 * 1024 * 1024,
	}

	handler := middleware.Pipeline(
		http.HandlerFunc(p.ProxyRequestHandler(rp)),
		ratelimiter.RateLimiter(limiter, nil),
		guard.RequestGuard,
		verifier.Authentication,
	)

	server := &http.Server{
		Addr:              ":9000",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("HORNET listening on :9000")

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
