package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	uc "github.com/assylzhan-a/company-task/internal/ports/usecase"
	"github.com/assylzhan-a/company-task/pkg/errors"
)

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userHandler struct {
	UserUseCase uc.UserUseCase
}

func NewUserHandler(r *chi.Mux, userUseCase uc.UserUseCase) {
	handler := &userHandler{
		UserUseCase: userUseCase,
	}
	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	if err := h.UserUseCase.Register(req.Username, req.Password); err != nil {
		errors.RespondWithError(w, errors.NewInternalServerError("Failed to register user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	token, err := h.UserUseCase.Login(req.Username, req.Password)
	if err != nil {
		errors.RespondWithError(w, errors.NewUnauthorizedError("Invalid username or password"))
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
