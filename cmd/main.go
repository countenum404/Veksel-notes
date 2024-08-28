package main

import (
	"os"

	"github.com/countenum404/Veksel/internal/api"
	"github.com/countenum404/Veksel/internal/repository/postgres"
	"github.com/countenum404/Veksel/internal/service"
)

const (
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_HOST     = "DB_HOST"
	DATABASE    = "DATABASE"
)

func main() {

	hardcodedPostgresCfg := map[string]string{
		"user":     os.Getenv(DB_USER),
		"host":     os.Getenv(DB_HOST),
		"database": os.Getenv(DATABASE),
		"password": os.Getenv(DB_PASSWORD),
	}

	repo, err := postgres.NewPostgresRepository(hardcodedPostgresCfg)
	if err != nil {
		panic(err)
	}
	pus := postgres.NewPostgresUserRepository(repo)
	pns := postgres.NewPostgresNotesRepository(repo)

	dus := service.NewDefaultUserService(pus)
	dns := service.NewDefaultNotesService(pns)

	a := api.NewApi(":4567", dns, dus)
	a.Run()
}
