package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
		return nil, errors.New("CANNOT FIND ANY USER NOTES")
	}
	return notes, nil
}

func (dns *DefaultNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) (*types.SpellResult, error) {
	err := dns.repo.CreateNote(userId, note.Header, note.Content)
	if err != nil {
		return nil, errors.New("CANNOT CREATE NOTE FOR USER")
	}
	return nil, nil
}

type SpellCheckNotesService struct {
	DefaultNotesService
	redisRepo  *redis.RedisRepository
	spellerUrl string
	maxLen     int
}

func NewSpellCheckNotesService(repo repository.NotesRepository, rdb *redis.RedisRepository, spellerUrl string, maxLen int) (*SpellCheckNotesService, error) {
	if spellerUrl == "" || maxLen <= 0 {
		return nil, errors.New("invalid parameters")
	}

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
		return nil, errors.New("CANNOT FIND ANY USER NOTES")
	}

	if err := ns.redisRepo.PutNotes(userId, notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func (ns *SpellCheckNotesService) CreateNote(note *types.CreateNoteRequest, userId int64) (*types.SpellResult, error) {
	result, err := ns.spellCheck(note.Content)
	if err != nil {
		return nil, err
	}

	if _, err := ns.DefaultNotesService.CreateNote(note, userId); err != nil {
		return nil, err
	}
	go func() {
		notes, err := ns.DefaultNotesService.GetNotes(userId)
		if err != nil {
			logger.GetLogger().Err(err)
		}
		ns.redisRepo.PutNotes(userId, notes)
	}()
	return result, nil
}

func (ns *SpellCheckNotesService) spellCheck(text string) (*types.SpellResult, error) {
	if len(text) > ns.maxLen {
		return nil, errors.New("NOTE IS TOO LARGE")
	}

	var url strings.Builder
	url.WriteString(ns.spellerUrl)
	url.WriteString(strings.Join(strings.Split(text, " "), "+"))

	var spellResult types.SpellResult
	req, err := http.Get(url.String())

	json.NewDecoder(req.Body).Decode(&spellResult)
	if err != nil {
		logger.GetLogger().Err(err)
	}
	req.Body.Close()

	return &spellResult, nil
}
