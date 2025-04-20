package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dapoadedire/fem_project/internal/api"
	"github.com/dapoadedire/fem_project/internal/store"
	"github.com/dapoadedire/fem_project/migrations"
)

type Application struct {
	Logger *log.Logger
	WorkourHandler *api.WorkourHandler
	DB *sql.DB
}

func NewApplication() (*Application, error) {

	// stores will go here
	pgDB, err:= store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)



	// our handlers will go here
	workoutStore:=store.NewPostgresWorkoutStore(pgDB)

	workoutHandler:=api.NewWorkoutHandler(workoutStore, logger)


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
