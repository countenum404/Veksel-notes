package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const JSON = "application/json"

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", JSON)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type MethodHandler struct {
	Methods map[string]func(http.ResponseWriter, *http.Request) error
}

func (m *MethodHandler) Call(method string, w http.ResponseWriter, r *http.Request) error {
	if handler, exists := m.Methods[r.Method]; exists {
		err := handler(w, r)
		if err != nil {
			return err
		}
	} else {
		return errors.New("METHOD NOT SUPPORTED")
	}
	return nil
}

type Api struct {
	listenAddr string
}

func NewApi(addr string) *Api {
	return &Api{listenAddr: addr}
}

func (a *Api) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", a.accountHandler)
	router.HandleFunc("/notes", a.notesHandler)

	log.Println("Veksel started")

	http.ListenAndServe(a.listenAddr, router)
}

func (a *Api) accountHandler(w http.ResponseWriter, r *http.Request) {
	methods := map[string]func(http.ResponseWriter, *http.Request) error{
		"GET": a.handleGetAccount,
	}
	mh := MethodHandler{Methods: methods}
	if err := mh.Call(r.Method, w, r); err != nil {
		WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}
}

func (a *Api) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *Api) notesHandler(w http.ResponseWriter, r *http.Request) {
	methods := map[string]func(http.ResponseWriter, *http.Request) error{
		"GET":  a.handleGetNote,
		"POST": a.handleCreateNote,
	}
	if handler, exists := methods[r.Method]; exists {
		err := handler(w, r)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	} else {
		WriteJson(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not supported"})
	}
}

func (a *Api) handleCreateNote(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *Api) handleGetNote(w http.ResponseWriter, r *http.Request) error {
	return nil
}
