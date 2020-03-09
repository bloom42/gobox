package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gitlab.com/bloom42/lily/rz"
	"gitlab.com/bloom42/lily/rz/log"
	"gitlab.com/bloom42/lily/rz/rzhttp"
)

func main() {
	env := os.Getenv("GO_ENV")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.SetLogger(log.With(
		rz.Fields(
			// rz.Caller(true),
			rz.String("service", "api"), rz.String("host", "abcd.local"), rz.String("environment", env),
		),
	))

	router := chi.NewRouter()

	// replace size field name by latency and disable userAgent logging
	loggingMiddleware := rzhttp.Handler(log.Logger(), rzhttp.Duration("latency"), rzhttp.UserAgent(""))

	// here the order matters, otherwise loggingMiddleware won't see the request ID
	router.Use(requestIDMiddleware)
	router.Use(loggingMiddleware)
	router.Use(injectLoggerMiddleware(log.Logger()))

	router.Get("/", helloWorld)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("listening", rz.Err(err))
	}
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuidv4, _ := uuid.NewRandom()
		requestID := uuidv4.String()
		w.Header().Set("X-Bloom-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), rzhttp.RequestIDCtxKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func injectLoggerMiddleware(logger rz.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rid, ok := r.Context().Value(rzhttp.RequestIDCtxKey).(string); ok {
				logger = logger.With(rz.Fields(rz.String("request_id", rid)))
				ctx := logger.ToCtx(r.Context())
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	// logger := rz.FromCtx(r.Context())
	// logger.Info("hello from GET /")
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}
