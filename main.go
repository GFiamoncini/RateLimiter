package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GFiamoncini/RateLimiter/config"
	"github.com/GFiamoncini/RateLimiter/limiter"
	"github.com/GFiamoncini/RateLimiter/middleware"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	rLimiter := limiter.NewRedisLimiter(cfg.RedisAddr, cfg.RedisPass)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rLimiter, cfg.LimitPerIP, cfg.LimitPerToken)

	r := mux.NewRouter()
	r.Use(rateLimiterMiddleware.Handle)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Bem Vindo")
	})

	log.Println("Server Iniciado na Porta :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
