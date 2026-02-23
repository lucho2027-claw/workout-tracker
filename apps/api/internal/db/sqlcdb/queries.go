package sqlcdb

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, `
INSERT INTO users (id, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, password_hash, created_at
`, arg.ID, arg.Email, arg.PasswordHash)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, `
SELECT id, email, password_hash, created_at
FROM users
WHERE email = $1
`, email)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

type Exercise struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	MuscleGroup string    `json:"muscle_group"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateExerciseParams struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	MuscleGroup string
	Notes       string
}

func (q *Queries) CreateExercise(ctx context.Context, arg CreateExerciseParams) (Exercise, error) {
	row := q.db.QueryRow(ctx, `
INSERT INTO exercises (id, user_id, name, muscle_group, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, muscle_group, notes, created_at
`, arg.ID, arg.UserID, arg.Name, arg.MuscleGroup, arg.Notes)
	var e Exercise
	err := row.Scan(&e.ID, &e.UserID, &e.Name, &e.MuscleGroup, &e.Notes, &e.CreatedAt)
	return e, err
}

func (q *Queries) ListExercisesByUser(ctx context.Context, userID uuid.UUID) ([]Exercise, error) {
	rows, err := q.db.Query(ctx, `
SELECT id, user_id, name, muscle_group, notes, created_at
FROM exercises
WHERE user_id = $1
ORDER BY created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Exercise, 0)
	for rows.Next() {
		var e Exercise
		if err := rows.Scan(&e.ID, &e.UserID, &e.Name, &e.MuscleGroup, &e.Notes, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

type WorkoutSession struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	PerformedAt time.Time `json:"performed_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateWorkoutSessionParams struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	PerformedAt time.Time
}

func (q *Queries) CreateWorkoutSession(ctx context.Context, arg CreateWorkoutSessionParams) (WorkoutSession, error) {
	row := q.db.QueryRow(ctx, `
INSERT INTO workout_sessions (id, user_id, performed_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, performed_at, created_at
`, arg.ID, arg.UserID, arg.PerformedAt)
	var s WorkoutSession
	err := row.Scan(&s.ID, &s.UserID, &s.PerformedAt, &s.CreatedAt)
	return s, err
}

type WorkoutSet struct {
	ID         uuid.UUID `json:"id"`
	SessionID  uuid.UUID `json:"session_id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
	SetNumber  int32     `json:"set_number"`
	Reps       int32     `json:"reps"`
	Weight     float64   `json:"weight"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateWorkoutSetParams struct {
	ID         uuid.UUID
	SessionID  uuid.UUID
	ExerciseID uuid.UUID
	SetNumber  int32
	Reps       int32
	Weight     float64
}

func (q *Queries) CreateWorkoutSet(ctx context.Context, arg CreateWorkoutSetParams) (WorkoutSet, error) {
	row := q.db.QueryRow(ctx, `
INSERT INTO workout_sets (id, session_id, exercise_id, set_number, reps, weight)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, session_id, exercise_id, set_number, reps, weight::float8, created_at
`, arg.ID, arg.SessionID, arg.ExerciseID, arg.SetNumber, arg.Reps, arg.Weight)
	var ws WorkoutSet
	err := row.Scan(&ws.ID, &ws.SessionID, &ws.ExerciseID, &ws.SetNumber, &ws.Reps, &ws.Weight, &ws.CreatedAt)
	return ws, err
}
