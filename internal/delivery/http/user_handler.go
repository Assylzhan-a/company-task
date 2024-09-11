package http

import (
	"encoding/json"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	uc "github.com/assylzhan-a/company-task/internal/ports/usecase"
	"github.com/assylzhan-a/company-task/pkg/errors"
	"github.com/go-chi/chi/v5"
	"net/http"
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
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	if err := h.UserUseCase.Register(ctx, req.Username, req.Password); err != nil {
		switch err {
		case entity.ErrEmptyUsername, entity.ErrEmptyPassword:
			errors.RespondWithError(w, errors.NewBadRequestError(err.Error()))
		case entity.ErrUsernameTaken:
			errors.RespondWithError(w, errors.NewBadRequestError("Username is already taken"))
		default:
			errors.RespondWithError(w, errors.NewInternalServerError("Failed to register user"))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	token, err := h.UserUseCase.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch err {
		case entity.ErrEmptyUsername, entity.ErrEmptyPassword, entity.ErrInvalidCredentials:
			errors.RespondWithError(w, errors.NewUnauthorizedError(err.Error()))
		default:
			errors.RespondWithError(w, errors.NewInternalServerError("Failed to log in"))
		}
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
