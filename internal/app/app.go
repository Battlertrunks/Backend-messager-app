package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Battlertrunks/cmd/api"
	"github.com/Battlertrunks/database"
	"github.com/Battlertrunks/internal/store"
)

// type Application struct {
// 	UserHandler
// }

type Application struct {
	Logger 	*log.Logger
	UserHandler	*api.UserHandler
	db				  *sql.DB
}

// Constructor
func NewApplication() (*Application, error) {
	Sqlite, err := store.Open()
	if err != nil {
		return nil, err
	}

	if err := store.MigrateFS(Sqlite, database.FS, "."); err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	userStore := store.NewSqliteUserStore(Sqlite)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		UserHandler: userHandler,
		Logger: logger,
		db: Sqlite,
	}

	return app, nil
}