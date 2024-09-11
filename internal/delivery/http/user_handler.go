package http

import (
	"encoding/json"
	"net/http"

	uc "github.com/assylzhan-a/company-task/internal/ports/usecase"
	"github.com/assylzhan-a/company-task/pkg/errors"
)

type UserHandler struct {
	userUseCase uc.UserUseCase
}

func NewUserHandler(userUseCase uc.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	if err := h.userUseCase.Register(req.Username, req.Password); err != nil {
		errors.RespondWithError(w, errors.NewInternalServerError("Failed to register user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	token, err := h.userUseCase.Login(req.Username, req.Password)
	if err != nil {
		errors.RespondWithError(w, errors.NewUnauthorizedError("Invalid username or password"))
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
