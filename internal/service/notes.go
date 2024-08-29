package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/countenum404/Veksel/internal/repository"
	"github.com/countenum404/Veksel/internal/types"
)

type DefaultNotesService struct {
	Repo repository.NotesRepository
}

func NewDefaultNotesService(repo repository.NotesRepository) (*DefaultNotesService, error) {
	return &DefaultNotesService{Repo: repo}, nil
}

func (dns *DefaultNotesService) GetNotes(userId int64) ([]types.Note, error) {
	notes, err := dns.Repo.GetNotesByUserId(userId)
	if err != nil {
		return nil, errors.New("CANNOT FIND ANY USER NOTES")
	}
	return notes, nil
}

func (dns *DefaultNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) (*types.SpellResult, error) {
	err := dns.Repo.CreateNote(userId, note.Header, note.Content)
	if err != nil {
		return nil, errors.New("CANNOT CREATE NOTE FOR USER")
	}
	return nil, nil
}

type SpellCheckNotesService struct {
	DefaultNotesService
	SpellerUrl string
	MaxLen     int
}

func NewSpellCheckNotesService(repo repository.NotesRepository, spellerUrl string, maxLen int) (*SpellCheckNotesService, error) {
	if spellerUrl == "" || maxLen <= 0 {
		return nil, errors.New("invalid parameters")
	}

	defaultNotesService, err := NewDefaultNotesService(repo)
	if err != nil {
		return nil, err
	}

	service := &SpellCheckNotesService{
		DefaultNotesService: *defaultNotesService,
		SpellerUrl:          spellerUrl,
		MaxLen:              maxLen,
	}

	service.SpellerUrl = spellerUrl
	service.MaxLen = maxLen

	return service, nil
}

func (ns *SpellCheckNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) (*types.SpellResult, error) {
	result, err := ns.spellCheck(note.Content)
	if err != nil {
		return nil, err
	}

	if _, err := ns.DefaultNotesService.CreateNote(note, userId); err != nil {
		return nil, err
	}

	return result, nil
}

func (ns *SpellCheckNotesService) spellCheck(text string) (*types.SpellResult, error) {
	if len(text) > ns.MaxLen {
		return nil, errors.New("NOTE IS TOO LARGE")
	}

	var url strings.Builder
	url.WriteString(ns.SpellerUrl)
	url.WriteString(strings.Join(strings.Split(text, " "), "+"))

	var spellResult types.SpellResult
	req, err := http.Get(url.String())

	json.NewDecoder(req.Body).Decode(&spellResult)
	if err != nil {
		log.Println(err)
	}
	req.Body.Close()

	return &spellResult, nil
}
