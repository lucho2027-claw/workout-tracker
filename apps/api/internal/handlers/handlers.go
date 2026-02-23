package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/lucho2027/workout-tracker/apps/api/internal/db/sqlcdb"
	"github.com/lucho2027/workout-tracker/apps/api/internal/models"
)

type Handler struct {
	DB        *pgxpool.Pool
	Q         *sqlcdb.Queries
	JWTSecret string
}

func New(db *pgxpool.Pool, jwtSecret string) *Handler {
	return &Handler{DB: db, Q: sqlcdb.New(db), JWTSecret: jwtSecret}
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
	_, err = h.Q.CreateUser(r.Context(), sqlcdb.CreateUserParams{
		ID:           id,
		Email:        req.Email,
		PasswordHash: string(hash),
	})
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

	u, err := h.Q.GetUserByEmail(r.Context(), req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID.String(),
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

	uuidUser, err := uuid.Parse(userID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid subject"})
		return
	}

	id := uuid.New()
	_, err = h.Q.CreateExercise(r.Context(), sqlcdb.CreateExerciseParams{
		ID:          id,
		UserID:      uuidUser,
		Name:        req.Name,
		MuscleGroup: req.MuscleGroup,
		Notes:       req.Notes,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

func (h *Handler) ListExercises(w http.ResponseWriter, r *http.Request, userID string) {
	uuidUser, err := uuid.Parse(userID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid subject"})
		return
	}

	items, err := h.Q.ListExercisesByUser(r.Context(), uuidUser)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
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

	uuidUser, err := uuid.Parse(userID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid subject"})
		return
	}

	id := uuid.New()
	_, err = h.Q.CreateWorkoutSession(r.Context(), sqlcdb.CreateWorkoutSessionParams{
		ID:          id,
		UserID:      uuidUser,
		PerformedAt: t,
	})
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

	sessionUUID, err := uuid.Parse(workoutID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid workout id"})
		return
	}
	exerciseUUID, err := uuid.Parse(req.ExerciseID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid exercise id"})
		return
	}

	id := uuid.New()
	_, err = h.Q.CreateWorkoutSet(r.Context(), sqlcdb.CreateWorkoutSetParams{
		ID:         id,
		SessionID:  sessionUUID,
		ExerciseID: exerciseUUID,
		SetNumber:  int32(req.SetNumber),
		Reps:       int32(req.Reps),
		Weight:     req.Weight,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}
