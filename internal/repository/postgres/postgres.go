package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/countenum404/Veksel/pkg/logger"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(config map[string]string) (*PostgresRepository, error) {
	opts := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config["user"],
		config["password"],
		config["host"],
		config["database"])
	logger.GetLogger().Info(opts)
	var dbpointer *sql.DB
	for {
		db, err := sql.Open("postgres", opts)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err == nil {
			logger.GetLogger().Info("Connected to db")
			dbpointer = db
			break
		}
		logger.GetLogger().Info("Waiting for db connection")
		time.Sleep(time.Second)
		continue
	}
	return &PostgresRepository{db: dbpointer}, nil
}
