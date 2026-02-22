package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lucho2027/workout-tracker/apps/api/internal/config"
	"github.com/lucho2027/workout-tracker/apps/api/internal/db"
	"github.com/lucho2027/workout-tracker/apps/api/internal/handlers"
)

func main() {
	cfg := config.Load()
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	h := handlers.New(pool, cfg.JWTSecret)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
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

	log.Printf("api listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

type authedHandler func(http.ResponseWriter, *http.Request, string)

func withAuth(jwtSecret string, next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid claims"}`, http.StatusUnauthorized)
			return
		}
		userID, _ := claims["sub"].(string)
		if userID == "" {
			http.Error(w, `{"error":"invalid subject"}`, http.StatusUnauthorized)
			return
		}
		next(w, r, userID)
	}
}
