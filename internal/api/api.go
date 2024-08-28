package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/countenum404/Veksel/internal/service"
	"github.com/countenum404/Veksel/internal/types"
	"github.com/gorilla/mux"
)

const JSON = "application/json"

type Api struct {
	notesService service.NotesService
	userService  service.UserService
	listenAddr   string
}

func NewApi(addr string, ns service.NotesService, us service.UserService) *Api {
	return &Api{listenAddr: addr, notesService: ns, userService: us}
}

func (a *Api) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/notes", a.BasicAuthMiddleware(a.notesHandler))

	log.Println("Veksel started")

	http.ListenAndServe(a.listenAddr, router)
}

func (a *Api) BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, ok := r.BasicAuth()
		if ok {
			handler(w, r)
			return
		}
		WriteJson(w, http.StatusUnauthorized, types.ApiError{Error: "You are unauthorized"})
	}
}

type HttpMethodHandler struct {
	Methods map[string]func(http.ResponseWriter, *http.Request) error
}

func (m *HttpMethodHandler) Call(method string, w http.ResponseWriter, r *http.Request) error {
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

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", JSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
