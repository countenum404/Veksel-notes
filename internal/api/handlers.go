package api

import (
	"encoding/json"
	"net/http"

	"github.com/countenum404/Veksel/internal/types"
)

func (a *Api) accountHandler(w http.ResponseWriter, r *http.Request) {
	methods := map[string]func(http.ResponseWriter, *http.Request) error{
		"GET": a.handleGetAccount,
	}
	mh := HttpMethodHandler{Methods: methods}
	if err := mh.Call(r.Method, w, r); err != nil {
		WriteJson(w, http.StatusBadRequest, types.ApiError{Error: err.Error()})
	}
}

func (a *Api) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *Api) notesHandler(w http.ResponseWriter, r *http.Request) {
	methods := map[string]func(http.ResponseWriter, *http.Request) error{
		"GET":  a.handleGetNotes,
		"POST": a.handleCreateNote,
	}
	mh := HttpMethodHandler{Methods: methods}
	if err := mh.Call(r.Method, w, r); err != nil {
		WriteJson(w, http.StatusBadRequest, types.ApiError{Error: err.Error()})
	}
}

func (a *Api) handleCreateNote(w http.ResponseWriter, r *http.Request) error {
	createNoteRequest := &types.CreateNoteRequest{}
	if err := json.NewDecoder(r.Body).Decode(createNoteRequest); err != nil {
		return err
	}
	username, password, _ := r.BasicAuth()
	user, err := a.userService.Authenticate(username, password)
	if err != nil {
		return err
	}
	if err := a.notesService.CreateNote(createNoteRequest, user.ID); err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, createNoteRequest)
	return nil
}

func (a *Api) handleGetNotes(w http.ResponseWriter, r *http.Request) error {
	username, password, _ := r.BasicAuth()
	user, err := a.userService.Authenticate(username, password)
	if err != nil {
		return err
	}
	notes, err := a.notesService.GetNotes(user.ID)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, notes)
	return nil
}
