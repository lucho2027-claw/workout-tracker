package models

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ExerciseCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	MuscleGroup string `json:"muscle_group"`
	Notes       string `json:"notes"`
}

type WorkoutCreateRequest struct {
	PerformedAt string `json:"performed_at" binding:"required"`
}

type SetCreateRequest struct {
	ExerciseID string  `json:"exercise_id" binding:"required,uuid"`
	SetNumber  int     `json:"set_number" binding:"required,min=1"`
	Reps       int     `json:"reps" binding:"required,min=1"`
	Weight     float64 `json:"weight" binding:"required,min=0"`
}
