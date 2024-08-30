package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/countenum404/Veksel/internal/repository"
	"github.com/countenum404/Veksel/internal/repository/redis"
	"github.com/countenum404/Veksel/internal/types"
	"github.com/countenum404/Veksel/pkg/logger"
)

type DefaultNotesService struct {
	repo repository.NotesRepository
}

func NewDefaultNotesService(repo repository.NotesRepository) (*DefaultNotesService, error) {
	return &DefaultNotesService{repo: repo}, nil
}

func (dns *DefaultNotesService) GetNotes(userId int64) ([]types.Note, error) {
	notes, err := dns.repo.GetNotesByUserId(userId)
	if err != nil {
		return nil, errors.New("cannot find any user notes")
	}
	return notes, nil
}

func (dns *DefaultNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) (*types.SpellResult, error) {
	err := dns.repo.CreateNote(userId, note.Header, note.Content)
	if err != nil {
		return nil, errors.New("cannot create note for user")
	}
	return nil, nil
}

type SpellCheckNotesService struct {
	DefaultNotesService
	redisRepo  *redis.RedisRepository
	spellerUrl url.URL
	maxLen     int
}

func NewSpellCheckNotesService(repo repository.NotesRepository, rdb *redis.RedisRepository, spellerUrl url.URL, maxLen int) (*SpellCheckNotesService, error) {

	defaultNotesService, err := NewDefaultNotesService(repo)
	if err != nil {
		return nil, err
	}

	service := &SpellCheckNotesService{
		DefaultNotesService: *defaultNotesService,
		redisRepo:           rdb,
		spellerUrl:          spellerUrl,
		maxLen:              maxLen,
	}

	service.spellerUrl = spellerUrl
	service.maxLen = maxLen

	return service, nil
}

func (ns *SpellCheckNotesService) GetNotes(userId int64) ([]types.Note, error) {
	notes, err := ns.redisRepo.GetNotesByUserId(userId)
	if err != nil {
		logger.GetLogger().Err(err)
	} else {
		return notes, nil
	}

	notes, err = ns.repo.GetNotesByUserId(userId)
	if err != nil {
		return nil, errors.New("cannot find any user notes")
	}

	if err := ns.redisRepo.PutNotes(userId, notes); err != nil {
		logger.GetLogger().Err(err)
	}

	return notes, nil
}

func (ns *SpellCheckNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) (*types.SpellResult, error) {
	result, err := ns.spellCheck(note.Content)
	if err != nil {
		logger.GetLogger().Err(err)
		return nil, err
	}

	if _, err := ns.DefaultNotesService.CreateNote(note, userId); err != nil {
		logger.GetLogger().Err(err)
		return nil, err
	}
	if err := ns.redisRepo.DeleteNotes(userId); err != nil {
		logger.GetLogger().Err(err)
	}
	return result, nil
}

func (ns *SpellCheckNotesService) spellCheck(text string) (*types.SpellResult, error) {
	if len(text) > ns.maxLen {
		return nil, errors.New("note is too large")
	}

	var spellResult types.SpellResult
	req, err := http.Get(ns.buildSpellerUrl(text))
	if err != nil {
		logger.GetLogger().Err(err)
	}
	json.NewDecoder(req.Body).Decode(&spellResult)
	req.Body.Close()
	return &spellResult, nil
}

func (ns *SpellCheckNotesService) buildSpellerUrl(text string) string {
	v := make(url.Values)
	v.Add("text", text)

	su := ns.spellerUrl
	su.RawQuery = v.Encode()

	return su.String()
}
