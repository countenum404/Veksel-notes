package postgres

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/countenum404/Veksel/pkg/logger"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(config map[string]string) (*PostgresRepository, error) {
	dataSource := NewDataSourceString(
		"postgres", config["host"], config["database"], config["user"], config["password"], "disable",
	)

	logger.GetLogger().Info(dataSource)
	var dbpointer *sql.DB
	for {
		db, err := sql.Open("postgres", dataSource)
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

func NewDataSourceString(proto, host, path, user, password, sslmode string) string {
	const SSLMODE = "sslmode"
	var v = make(url.Values)
	v.Add(SSLMODE, sslmode)

	var u = url.URL{
		Scheme:   proto,
		Host:     host,
		Path:     path,
		User:     url.UserPassword(user, password),
		RawQuery: v.Encode(),
	}
	return u.String()
}
