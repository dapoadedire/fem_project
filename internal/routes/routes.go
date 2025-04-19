package routes

import (
	"github.com/dapoadedire/fem_project/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkourHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkourHandler.HandleCreateWorkout)
	return r
}
