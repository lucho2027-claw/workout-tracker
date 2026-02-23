-- name: CreateExercise :one
INSERT INTO exercises (id, user_id, name, muscle_group, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, muscle_group, notes, created_at;

-- name: ListExercisesByUser :many
SELECT id, user_id, name, muscle_group, notes, created_at
FROM exercises
WHERE user_id = $1
ORDER BY created_at DESC;
