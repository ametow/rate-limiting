package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/ametow/rate-limiting/limiter"
)

func getClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	return host
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	ipLimiterMap := sync.Map{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		ipLimiter, _ := ipLimiterMap.LoadOrStore(ip, limiter.NewTokenBucketLimiter(2.0, 10))

		if !ipLimiter.(limiter.RateLimiter).Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("hi there\n"))
	})

	// handler := rateLimiterMiddleware(mux)
	handler := rateLimiterMiddleware(mux)

	err := http.ListenAndServe("127.0.0.1:8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
