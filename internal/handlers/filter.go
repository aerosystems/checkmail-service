package handlers

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	RPCClient "github.com/aerosystems/checkmail-service/internal/rpc_client"
	AuthService "github.com/aerosystems/checkmail-service/pkg/auth_service"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type FilterCreateRequest struct {
	Name         string `json:"name" example:"gmail.com"`
	Type         string `json:"type" example:"whitelist"`
	Coverage     string `json:"coverage" example:"equals"`
	ProjectToken string `json:"projectToken" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
}

func (cr *FilterCreateRequest) Validate() *CustomError.Error {
	if err := validators.ValidateDomainTypes(cr.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainCoverage(cr.Coverage); err != nil {
		return err
	}
	if err := validators.ValidateDomainName(cr.Name); err != nil {
		return err
	}
	return nil
}

type FilterUpdateRequest struct {
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

func (ur *FilterUpdateRequest) Validate() *CustomError.Error {
	if err := validators.ValidateDomainTypes(ur.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainCoverage(ur.Coverage); err != nil {
		return err
	}
	return nil
}

// CreateFilter godoc
// @Summary Create Filter
// @Description Create Filter for Project. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param filter body FilterCreateRequest true "raw request body"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filters [post]
func (h *BaseHandler) CreateFilter(w http.ResponseWriter, r *http.Request) {
	accessTokenClaims := r.Context().Value(helpers.ContextKey("accessTokenClaimsKey")).(*AuthService.AccessTokenClaims)
	var requestPayload FilterCreateRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}
	if err := requestPayload.Validate(); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}
	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, err.Error(), err))
		return
	}

	result, err := RPCClient.GetProject(requestPayload.ProjectToken)
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400202, "projectToken does not exist", err))
		return
	}
	if !helpers.Contains([]string{"staff"}, accessTokenClaims.UserRole) {
		if result.Token != requestPayload.ProjectToken {
			_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", err))
			return
		}
	}

	newFilter := models.Filter{
		Name:         requestPayload.Name,
		Type:         requestPayload.Type,
		Coverage:     requestPayload.Coverage,
		ProjectToken: requestPayload.ProjectToken,
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

// GetFilterList godoc
// @Summary Get Filter List
// @Description Get Filter List for all user Projects. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param userId query int false "user id"
// @Param projectToken query string false "project token"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filters [get]
func (h *BaseHandler) GetFilterList(w http.ResponseWriter, r *http.Request) {
	accessTokenClaims := r.Context().Value(helpers.ContextKey("accessTokenClaimsKey")).(*AuthService.AccessTokenClaims)
	var filters []models.Filter
	var userIdStr, projectToken string
	var err error
	var userId int
	userIdStr = r.URL.Query().Get("userId")
	if userIdStr != "" {
		userId, err = strconv.Atoi(userIdStr)
		if err != nil {
			_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, "user id does not valid", err))
			return
		}
	}
	projectToken = r.URL.Query().Get("projectToken")
	switch accessTokenClaims.UserRole {
	case "business":
		if projectToken == "" {
			result, err := RPCClient.GetProjectList(accessTokenClaims.UserId)
			if err != nil {
				_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not find filters", err))
				return
			}
			if len(*result) == 0 {
				_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404205, "projects not found for user", nil))
				return
			}
			for _, project := range *result {
				if projectFilters, err := h.filterRepo.FindByProjectToken(project.Token); err == nil {
					filters = append(filters, *projectFilters)
				} else {
					_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not find filters for project", err))
					return
				}
			}
		} else {
			result, err := RPCClient.GetProject(projectToken)
			if err != nil {
				_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not find filters", err))
				return
			}
			if result.UserId != accessTokenClaims.UserId {
				_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", err))
				return
			}
			if projectFilters, err := h.filterRepo.FindByProjectToken(result.Token); err == nil {
				filters = append(filters, *projectFilters)
			} else {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404205, "filters not found for project", err))
					return
				}
				_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not find filters for project", err))
				return
			}
		}
	case "staff":
		if userIdStr != "" {
			result, err := RPCClient.GetProjectList(userId)
			if err != nil {
				_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not find filters", err))
				return
			}
			if len(*result) == 0 {
				_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404205, "projects not found for user", nil))
				return
			}
			for _, project := range *result {
				if projectToken != "" {
					if project.Token != projectToken {
						continue
					}
				}
				projectFilters, err := h.filterRepo.FindByProjectToken(project.Token)
				if err != nil {
					continue
				}
				h.log.Info(projectFilters)
				filters = append(filters, *projectFilters)
			}
		} else {
			if projectToken != "" {
				result, err := RPCClient.GetProject(projectToken)
				if err != nil {
					_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500202, "could not find filters", err))
					return
				}
				if projectFilters, err := h.filterRepo.FindByProjectToken(result.Token); err == nil {
					filters = append(filters, *projectFilters)
				} else {
					_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not find filters for project", err))
					return
				}
			} else {
				allFilters, err := h.filterRepo.FindAll()
				if err != nil {
					_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not find filters", err))
					return
				}
				filters = *allFilters
			}
		}

	default:
		_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", nil))
		return
	}
	if len(filters) == 0 {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404205, "filters not found", nil))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("filters found", filters))
	return
}

// UpdateFilter godoc
// @Summary Update Filter
// @Description Update Filter for Project by ID. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int true "filter id"
// @Param filter body FilterUpdateRequest true "raw request body"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filters/{filterId} [put]
func (h *BaseHandler) UpdateFilter(w http.ResponseWriter, r *http.Request) {
	accessTokenClaims := r.Context().Value(helpers.ContextKey("accessTokenClaimsKey")).(*AuthService.AccessTokenClaims)
	filterId, err := strconv.Atoi(chi.URLParam(r, "filterId"))
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, "filter id does not valid", err))
		return
	}
	var requestPayload FilterUpdateRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422201, "could not read request body", err))
		return
	}

	if err := requestPayload.Validate(); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(err.Code, err.Message, err.Error()))
		return
	}

	filter, err := h.filterRepo.FindByID(filterId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404205, "filter not found", err))
		return
	}
	if err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not find filter", err))
		return
	}

	switch accessTokenClaims.UserRole {
	case "business":
		result, err := RPCClient.GetProject(filter.ProjectToken)
		if err != nil {
			_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not update filter", err))
			return
		}
		if result.Token != filter.ProjectToken {
			_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", err))
			return
		}
	case "staff":
		break
	default:
		_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", nil))
		return
	}

	filter.Type = requestPayload.Type
	filter.Coverage = requestPayload.Coverage

	if err := h.filterRepo.Update(filter); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not update filter", err))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("filter updated", filter))
	return
}

// DeleteFilter godoc
// @Summary Delete Filter
// @Description Delete Filter for Project by ID. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int true "filter id"
// @Success 204 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/filters/{filterId} [delete]
func (h *BaseHandler) DeleteFilter(w http.ResponseWriter, r *http.Request) {
	accessTokenClaims := r.Context().Value(helpers.ContextKey("accessTokenClaimsKey")).(*AuthService.AccessTokenClaims)
	filterId, err := strconv.Atoi(chi.URLParam(r, "filterId"))
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400205, "filter id does not valid", err))
		return
	}

	filter, err := h.filterRepo.FindByID(filterId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404205, "filter not found", err))
		return
	}
	if err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not find filter", err))
		return
	}

	switch accessTokenClaims.UserRole {
	case "business":
		result, err := RPCClient.GetProject(filter.ProjectToken)
		if err != nil {
			_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not delete filter", err))
			return
		}
		if result.Token != filter.ProjectToken {
			_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", err))
			return
		}
	case "staff":
		break
	default:
		_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403201, "access denied", nil))
		return
	}
	if err := h.filterRepo.Delete(filter); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500205, "could not delete filter", err))
		return
	}
	_ = WriteResponse(w, http.StatusNoContent, NewResponsePayload("filter deleted", nil))
}
