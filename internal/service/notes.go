package service

import (
	"errors"

	"github.com/countenum404/Veksel/internal/repository"
	"github.com/countenum404/Veksel/internal/types"
)

type DefaultNotesService struct {
	Repo repository.NotesRepository
}

func NewDefaultNotesService(repo repository.NotesRepository) *DefaultNotesService {
	return &DefaultNotesService{Repo: repo}
}

func (dns *DefaultNotesService) GetNotes(userId int64) ([]types.Note, error) {
	notes, err := dns.Repo.GetNotesByUserId(userId)
	if err != nil {
		return nil, errors.New("CANNOT FIND ANY USER NOTES")
	}
	return notes, nil
}

func (dns *DefaultNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) error {
	err := dns.Repo.CreateNote(userId, note.Content, note.Header)
	if err != nil {
		return errors.New("CANNOT CREATE NOTE FOR USER")
	}
	return nil
}
