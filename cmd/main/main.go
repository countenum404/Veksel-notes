package main

import (
	"net/url"
	"os"

	"github.com/countenum404/Veksel/internal/api"
	"github.com/countenum404/Veksel/internal/repository/postgres"
	"github.com/countenum404/Veksel/internal/repository/redis"
	"github.com/countenum404/Veksel/internal/service"
)

const (
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_HOST     = "DB_HOST"
	DATABASE    = "DATABASE"
)

const MAX_CONTENT_LEN = 10000
const REDIS_HOST = "REDIS_HOST"

func main() {

	SpellingUrl := url.URL{
		Scheme: "https",
		Host:   "speller.yandex.net",
		Path:   "/services/spellservice.json/checkText",
	}

	rdb := redis.NewRedisRepository(os.Getenv(REDIS_HOST), "", 0)

	cfg := map[string]string{
		"user":     os.Getenv(DB_USER),
		"host":     os.Getenv(DB_HOST),
		"database": os.Getenv(DATABASE),
		"password": os.Getenv(DB_PASSWORD),
	}

	repo, err := postgres.NewPostgresRepository(cfg)
	if err != nil {
		panic(err)
	}
	pus := postgres.NewPostgresUserRepository(repo)
	pns := postgres.NewPostgresNotesRepository(repo)

	dus := service.NewDefaultUserService(pus)
	dns, _ := service.NewSpellCheckNotesService(pns, rdb, SpellingUrl, MAX_CONTENT_LEN)

	a := api.NewApi(":4567", dns, dus)
	a.Run()
}
