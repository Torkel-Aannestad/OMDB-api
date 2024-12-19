package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/database"
	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
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

	go func() {
		for {
			time.Sleep(time.Minute * 1)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Minute*3 {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: *rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				app.rateLimitExceededResponse(w, r)
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authRateLimit(next http.HandlerFunc) http.HandlerFunc {
	type client struct {
		limiter  rate.Limiter
		lastSeen time.Time
	}

	var mu sync.Mutex
	clients := make(map[string]*client)

	go func() {
		for {
			time.Sleep(time.Minute * 1)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Hour*1 {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.authLimiter.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: *rate.NewLimiter(rate.Limit(app.config.authLimiter.rps), app.config.authLimiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				app.rateLimitExceededResponse(w, r)
				mu.Unlock()
				return
			}
			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})

}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, database.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		authToken := headerParts[1]

		v := validator.New()
		database.ValidateTokenPlaintext(v, authToken)
		valid := v.Valid()
		if !valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(database.ScopeAuthentication, authToken)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		r = app.contextSetUser(r, user)

		r.UserAgent()

		next.ServeHTTP(w, r)
	})

}

func (app *application) protectedRoute(permissionCode string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		if permissionCode == "" {
			next.ServeHTTP(w, r)
			return
		}
		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Include(permissionCode) {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
