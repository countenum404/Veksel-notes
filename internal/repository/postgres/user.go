package postgres

import (
	"github.com/countenum404/Veksel/internal/types"
)

type PostgresUserRepository struct {
	*PostgresRepository
}

func NewPostgresUserRepository(repo *PostgresRepository) *PostgresUserRepository {
	return &PostgresUserRepository{PostgresRepository: repo}
}

func (pus *PostgresUserRepository) GetUser(username string) (*types.User, error) {
	row := pus.db.QueryRow("SELECT * FROM users WHERE username=$1", username)
	existingUser := types.User{}
	row.Scan(&existingUser.ID, &existingUser.Fisrtname, &existingUser.Lastname, &existingUser.Username, &existingUser.Password)
	return &existingUser, nil
}
