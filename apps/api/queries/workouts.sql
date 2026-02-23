-- name: CreateWorkoutSession :one
INSERT INTO workout_sessions (id, user_id, performed_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, performed_at, created_at;

-- name: CreateWorkoutSet :one
INSERT INTO workout_sets (id, session_id, exercise_id, set_number, reps, weight)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, session_id, exercise_id, set_number, reps, weight, created_at;
