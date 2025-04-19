package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dapoadedire/fem_project/internal/api"
)

type Application struct {
	Logger *log.Logger
	WorkourHandler *api.WorkourHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores will go here


	// our handlers will go here

	workoutHandler:=api.NewWorkoutHandler()


	app := &Application{
		Logger: logger,
		WorkourHandler: workoutHandler,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
