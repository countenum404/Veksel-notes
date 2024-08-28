package postgres

import (
	"github.com/countenum404/Veksel/internal/types"
)

type PostgresNotesRepository struct {
	*PostgresRepository
}

func NewPostgresNotesRepository(repo *PostgresRepository) *PostgresNotesRepository {
	return &PostgresNotesRepository{PostgresRepository: repo}
}

func (pns *PostgresNotesRepository) GetNotesByUserId(userId int64) ([]types.Note, error) {
	query := "SELECT * FROM notes WHERE user_id=$1"
	rows, err := pns.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []types.Note{}
	for rows.Next() {
		note := new(types.Note)
		var uid int
		err := rows.Scan(&note.ID, &note.Header, &note.Content, &uid)
		if err != nil {
			return nil, err
		}
		notes = append(notes, *note)
	}
	return notes, nil
}

func (pns *PostgresNotesRepository) CreateNote(userId int64, header string, content string) error {
	query := "INSERT INTO notes (header, content, user_id) VALUES ($1, $2, $3)"
	if _, err := pns.db.Exec(query, header, content, userId); err != nil {
		return err
	}
	return nil
}
