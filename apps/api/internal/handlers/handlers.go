package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	id := uuid.New()
	_, err := h.DB.Exec(c, `INSERT INTO users (id, email, password_hash) VALUES ($1,$2,$3)`, id, req.Email, string(hash))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id.String()})
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id, hash string
	err := h.DB.QueryRow(c, `SELECT id::text, password_hash FROM users WHERE email=$1`, req.Email).Scan(&id, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte(h.JWTSecret))
	c.JSON(http.StatusOK, gin.H{"token": s})
}

func (h *Handler) CreateExercise(c *gin.Context) {
	var req models.ExerciseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	id := uuid.New()
	_, err := h.DB.Exec(c, `INSERT INTO exercises (id,user_id,name,muscle_group,notes) VALUES ($1,$2,$3,$4,$5)`, id, userID, req.Name, req.MuscleGroup, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) ListExercises(c *gin.Context) {
	userID := c.GetString("user_id")
	rows, err := h.DB.Query(c, `SELECT id::text,name,muscle_group,notes FROM exercises WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	items := []gin.H{}
	for rows.Next() {
		var id, name, mg, notes string
		_ = rows.Scan(&id, &name, &mg, &notes)
		items = append(items, gin.H{"id": id, "name": name, "muscle_group": mg, "notes": notes})
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateWorkout(c *gin.Context) {
	var req models.WorkoutCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	t, err := time.Parse(time.RFC3339, req.PerformedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "performed_at must be RFC3339"})
		return
	}
	id := uuid.New()
	_, err = h.DB.Exec(c, `INSERT INTO workout_sessions (id,user_id,performed_at) VALUES ($1,$2,$3)`, id, userID, t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) AddWorkoutSet(c *gin.Context) {
	workoutID := c.Param("id")
	var req models.SetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New()
	_, err := h.DB.Exec(c, `INSERT INTO workout_sets (id,session_id,exercise_id,set_number,reps,weight) VALUES ($1,$2,$3,$4,$5,$6)`, id, workoutID, req.ExerciseID, req.SetNumber, req.Reps, req.Weight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
