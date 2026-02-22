package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lucho2027/workout-tracker/apps/api/internal/config"
	"github.com/lucho2027/workout-tracker/apps/api/internal/db"
	"github.com/lucho2027/workout-tracker/apps/api/internal/handlers"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "change-me" && cfg.AppEnv != "development" {
		log.Fatal("JWT_SECRET must be set in non-development environments")
	}

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	h := handlers.New(pool, cfg.JWTSecret)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{"ok": false, "error": "database unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)

	mux.HandleFunc("/exercises", withAuth(cfg.JWTSecret, func(w http.ResponseWriter, r *http.Request, userID string) {
		switch r.Method {
		case http.MethodGet:
			h.ListExercises(w, r, userID)
		case http.MethodPost:
			h.CreateExercise(w, r, userID)
		default:
			http.NotFound(w, r)
		}
	}))

	mux.HandleFunc("POST /workouts", withAuth(cfg.JWTSecret, func(w http.ResponseWriter, r *http.Request, userID string) {
		h.CreateWorkout(w, r, userID)
	}))

	mux.HandleFunc("/workouts/", withAuth(cfg.JWTSecret, func(w http.ResponseWriter, r *http.Request, userID string) {
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/sets") {
			http.NotFound(w, r)
			return
		}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 3 || parts[0] != "workouts" || parts[2] != "sets" {
			http.NotFound(w, r)
			return
		}
		workoutID := parts[1]
		h.AddWorkoutSet(w, r, workoutID)
	}))

	handler := withMiddleware(mux, cfg.AllowedOrigin)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("api listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

type authedHandler func(http.ResponseWriter, *http.Request, string)

func withAuth(jwtSecret string, next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
			return
		}
		userID, _ := claims["sub"].(string)
		if userID == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid subject"})
			return
		}
		next(w, r, userID)
	}
}

func withMiddleware(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		origin := r.Header.Get("Origin")
		if origin != "" && origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
