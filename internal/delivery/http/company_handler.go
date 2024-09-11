package http

import (
	"encoding/json"
	"github.com/assylzhan-a/company-task/internal/auth"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	uc "github.com/assylzhan-a/company-task/internal/ports/usecase"
	"github.com/assylzhan-a/company-task/pkg/errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type companyHandler struct {
	companyUseCase uc.CompanyUseCase
}

func NewCompanyHandler(r *chi.Mux, useCase uc.CompanyUseCase) {
	handler := &companyHandler{
		companyUseCase: useCase,
	}
	// Public route (getter)
	r.Get("/v1/companies/{id}", handler.Get)

	// Protected routes
	r.Route("/v1/companies", func(r chi.Router) {
		r.Use(auth.JWTAuth)
		r.Post("/", handler.Create)
		r.Patch("/{id}", handler.Patch)
		r.Delete("/{id}", handler.Delete)
	})
}

func (h *companyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var company entity.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	if err := company.Validate(); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.companyUseCase.Create(r.Context(), &company); err != nil {
		errors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

func (h *companyHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid company ID"))
		return
	}

	var patchCompany entity.PatchCompany
	if err := json.NewDecoder(r.Body).Decode(&patchCompany); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid request payload"))
		return
	}

	if err := patchCompany.Validate(); err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.companyUseCase.Patch(r.Context(), id, &patchCompany); err != nil {
		errors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Company updated successfully"})
}

func (h *companyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid company ID"))
		return
	}

	if err := h.companyUseCase.Delete(r.Context(), id); err != nil {
		errors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *companyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		errors.RespondWithError(w, errors.NewBadRequestError("Invalid company ID"))
		return
	}

	company, err := h.companyUseCase.GetByID(r.Context(), id)
	if err != nil {
		errors.RespondWithError(w, err)
		return
	}

	json.NewEncoder(w).Encode(company)
}
