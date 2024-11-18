package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  rate.Limiter
		lastSeen time.Time
	}

	var mu sync.Mutex
	clients := make(map[string]*client)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: *rate.NewLimiter(2, 4),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			mu.Unlock()
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
