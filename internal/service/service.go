package service

import (
	"github.com/countenum404/Veksel/internal/types"
)

type UserService interface {
	Authenticate(username, password string) (*types.User, error)
}

type NotesService interface {
	GetNotes(userId int64) ([]types.Note, error)
	CreateNote(note *types.CreateNoteRequest, userId int64) error
}
