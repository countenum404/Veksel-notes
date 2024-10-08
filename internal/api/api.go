package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/countenum404/Veksel/internal/service"
	"github.com/countenum404/Veksel/internal/types"
	"github.com/countenum404/Veksel/pkg/logger"
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

	logger.GetLogger().Info("Veksel started")

	http.ListenAndServe(a.listenAddr, router)
}

func (a *Api) BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			logger.GetLogger().Info("Forbidden", username)
			WriteJson(w, http.StatusForbidden, types.ApiError{Error: "Authentication failed: username and password are required"})
			return
		}
		user, err := a.userService.GetUser(username, password)
		if err != nil && user == nil {
			logger.GetLogger().Err("Can't authorize user", username, err)
			WriteJson(w, http.StatusUnauthorized, types.ApiError{Error: "Invalid username or password"})
			return
		}
		handler(w, r)
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
