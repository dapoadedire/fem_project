package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"database/sql"

	"github.com/dapoadedire/fem_project/internal/api"
	"github.com/dapoadedire/fem_project/internal/store"
)

type Application struct {
	Logger *log.Logger
	WorkourHandler *api.WorkourHandler
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores will go here
	pgDB, err:= store.Open()
	if err != nil {
		return nil, err
	}


	// our handlers will go here

	workoutHandler:=api.NewWorkoutHandler()


	app := &Application{
		Logger: logger,
		WorkourHandler: workoutHandler,
		DB: pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
