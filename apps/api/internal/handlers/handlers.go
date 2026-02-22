package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/lucho2027/workout-tracker/apps/api/internal/models"
)

type Handler struct {
	DB        *pgxpool.Pool
	JWTSecret string
}

func New(db *pgxpool.Pool, jwtSecret string) *Handler {
	return &Handler{DB: db, JWTSecret: jwtSecret}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSON(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := decodeJSON(r, &req); err != nil || req.Email == "" || len(req.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "password hashing failed"})
		return
	}
	id := uuid.New()
	_, err = h.DB.Exec(r.Context(), `INSERT INTO users (id, email, password_hash) VALUES ($1,$2,$3)`, id, req.Email, string(hash))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email already in use"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := decodeJSON(r, &req); err != nil || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	var id, hash string
	err := h.DB.QueryRow(r.Context(), `SELECT id::text, password_hash FROM users WHERE email=$1`, req.Email).Scan(&id, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": signed})
}

func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request, userID string) {
	var req models.ExerciseCreateRequest
	if err := decodeJSON(r, &req); err != nil || req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	id := uuid.New()
	_, err := h.DB.Exec(r.Context(), `INSERT INTO exercises (id,user_id,name,muscle_group,notes) VALUES ($1,$2,$3,$4,$5)`, id, userID, req.Name, req.MuscleGroup, req.Notes)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

func (h *Handler) ListExercises(w http.ResponseWriter, r *http.Request, userID string) {
	rows, err := h.DB.Query(r.Context(), `SELECT id::text,name,muscle_group,notes FROM exercises WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := []map[string]string{}
	for rows.Next() {
		var id, name, mg, notes string
		_ = rows.Scan(&id, &name, &mg, &notes)
		items = append(items, map[string]string{"id": id, "name": name, "muscle_group": mg, "notes": notes})
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request, userID string) {
	var req models.WorkoutCreateRequest
	if err := decodeJSON(r, &req); err != nil || req.PerformedAt == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	t, err := time.Parse(time.RFC3339, req.PerformedAt)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "performed_at must be RFC3339"})
		return
	}

	id := uuid.New()
	_, err = h.DB.Exec(r.Context(), `INSERT INTO workout_sessions (id,user_id,performed_at) VALUES ($1,$2,$3)`, id, userID, t)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

func (h *Handler) AddWorkoutSet(w http.ResponseWriter, r *http.Request, workoutID string) {
	var req models.SetCreateRequest
	if err := decodeJSON(r, &req); err != nil || req.ExerciseID == "" || req.SetNumber < 1 || req.Reps < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	id := uuid.New()
	_, err := h.DB.Exec(r.Context(), `INSERT INTO workout_sets (id,session_id,exercise_id,set_number,reps,weight) VALUES ($1,$2,$3,$4,$5,$6)`, id, workoutID, req.ExerciseID, req.SetNumber, req.Reps, req.Weight)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}
