package handlers

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type FilterRequest struct {
	Name     string `json:"name" example:"gmail.com"`
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

func (r *FilterRequest) Validate() *CustomError.Error {
	if err := validators.ValidateDomainTypes(r.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainCoverage(r.Coverage); err != nil {
		return err
	}
	if err := validators.ValidateDomainName(r.Name); err != nil {
		return err
	}
	return nil
}

func (h *BaseHandler) CreateFilter(w http.ResponseWriter, r *http.Request) {
	var requestPayload FilterRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}

	if err := requestPayload.Validate(); err != nil {
		WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}

	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	newFilter := models.Filter{
		Name:     requestPayload.Name,
		Type:     requestPayload.Type,
		Coverage: requestPayload.Coverage,
	}

	if err := h.filterRepo.Create(&newFilter); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409201, "filter for this domain name already exists", err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not create filter", err))
		return
	}

	_ = WriteResponse(w, http.StatusCreated, NewResponsePayload("filter created", newFilter))
	return
}
