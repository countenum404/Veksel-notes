package repository

import (
	"github.com/countenum404/Veksel/internal/types"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	GetUser(username string) (*types.User, error)
}

type NotesRepository interface {
	GetNotesByUserId(userId int64) ([]types.Note, error)
	CreateNote(userId int64, header string, content string) error
}
