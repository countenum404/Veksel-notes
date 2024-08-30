package api

import (
	"encoding/json"
	"net/http"

	"github.com/countenum404/Veksel/internal/types"
	"github.com/countenum404/Veksel/pkg/logger"
)

func (a *Api) notesHandler(w http.ResponseWriter, r *http.Request) {
	methods := map[string]func(http.ResponseWriter, *http.Request) error{
		"GET":  a.handleGetNotes,
		"POST": a.handlePostNote,
	}
	mh := HttpMethodHandler{Methods: methods}
	if err := mh.Call(r.Method, w, r); err != nil {
		WriteJson(w, http.StatusBadRequest, types.ApiError{Error: err.Error()})
	}
}

func (a *Api) handlePostNote(w http.ResponseWriter, r *http.Request) error {
	username, password, _ := r.BasicAuth()
	user, _ := a.userService.GetUser(username, password)

	createNoteRequest := &types.CreateNoteRequest{}
	if err := json.NewDecoder(r.Body).Decode(createNoteRequest); err != nil {
		return err
	}

	spells, err := a.notesService.CreateNote(createNoteRequest, user.ID)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, types.SpellResponse{NoteRequest: *createNoteRequest, Spells: *spells})
	logger.GetLogger().Info("Note", createNoteRequest.Header, "created for User id:", user.ID, "noteLen:", len(createNoteRequest.Content))
	return nil
}

func (a *Api) handleGetNotes(w http.ResponseWriter, r *http.Request) error {
	username, password, _ := r.BasicAuth()
	user, _ := a.userService.GetUser(username, password)

	notes, err := a.notesService.GetNotes(user.ID)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, notes)
	logger.GetLogger().Info("Notes sent to User id:", user.ID, "firstname:", user.Fisrtname, "lastname:", user.Lastname)
	return nil
}
