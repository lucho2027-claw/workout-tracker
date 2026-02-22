package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/lucho2027/workout-tracker/apps/api/internal/config"
	"github.com/lucho2027/workout-tracker/apps/api/internal/db"
	"github.com/lucho2027/workout-tracker/apps/api/internal/handlers"
	"github.com/lucho2027/workout-tracker/apps/api/internal/middleware"
)

func main() {
	cfg := config.Load()
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	h := handlers.New(pool, cfg.JWTSecret)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)

	auth := r.Group("/")
	auth.Use(middleware.Auth(cfg.JWTSecret))
	auth.GET("/exercises", h.ListExercises)
	auth.POST("/exercises", h.CreateExercise)
	auth.POST("/workouts", h.CreateWorkout)
	auth.POST("/workouts/:id/sets", h.AddWorkoutSet)

	log.Printf("api listening on :%s", cfg.Port)
	_ = r.Run(":" + cfg.Port)
}
