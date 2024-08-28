package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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
	log.Println(opts)
	var dbpointer *sql.DB
	for {
		db, err := sql.Open("postgres", opts)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err == nil {
			log.Println("Connected to db")
			dbpointer = db
			break
		}
		log.Println("Waiting for db connection")
		time.Sleep(time.Second)
		continue
	}
	return &PostgresRepository{db: dbpointer}, nil
}
