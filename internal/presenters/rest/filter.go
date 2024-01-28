package rest

import (
	"errors"
	"github.com/aerosystems/checkmail-service/internal/helpers"
	"github.com/aerosystems/checkmail-service/internal/models"
	RPCClient "github.com/aerosystems/checkmail-service/internal/rpc_client"
	CustomError "github.com/aerosystems/checkmail-service/pkg/custom_error"
	OAuthService "github.com/aerosystems/checkmail-service/pkg/oauth_service"
	"github.com/aerosystems/checkmail-service/pkg/validators"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
func (h *BaseHandler) CreateFilter(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*OAuthService.AccessTokenClaims)
	var requestPayload FilterCreateRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}
	if err := requestPayload.Validate(); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}
	root, _ := helpers.GetRootDomain(requestPayload.Name)
	rootDomain, _ := h.rootDomainRepo.FindByName(root)
	if rootDomain == nil {
		err := errors.New("domain does not exist")
		return h.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}

	result, err := RPCClient.GetProject(requestPayload.ProjectToken)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
	}
	if !helpers.Contains([]string{"staff"}, accessTokenClaims.UserRole) {
		if result.Token != requestPayload.ProjectToken {
			return h.ErrorResponse(c, http.StatusForbidden, "access denied", err)
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
			return h.ErrorResponse(c, http.StatusConflict, "filter for this domain name already exists", err)
		}
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not create filter", err)
	}
	return h.SuccessResponse(c, http.StatusCreated, "filter successfully created", newFilter)
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
func (h *BaseHandler) GetFilterList(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*OAuthService.AccessTokenClaims)
	var filters []models.Filter
	var userIdStr, projectToken string
	var err error
	var userUuid uuid.UUID
	userIdStr = c.QueryParam("userId")
	if userIdStr != "" {
		userUuid, err = uuid.Parse(userIdStr)
		if err != nil {
			return h.ErrorResponse(c, http.StatusBadRequest, "user id does not valid", err)
		}
	}
	projectToken = c.QueryParam("projectToken")
	switch accessTokenClaims.UserRole {
	case "business":
		if projectToken == "" {
			result, err := RPCClient.GetProjectList(uuid.MustParse(accessTokenClaims.UserUuid))
			if err != nil {
				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
			}
			if len(*result) == 0 {
				return h.ErrorResponse(c, http.StatusNotFound, "projects not found for user", nil)
			}
			for _, project := range *result {
				if projectFilters, err := h.filterRepo.FindByProjectToken(project.Token); err == nil {
					filters = append(filters, *projectFilters)
				} else {
					return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters for project", err)
				}
			}
		} else {
			result, err := RPCClient.GetProject(projectToken)
			if err != nil {
				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
			}
			if result.UserUuid != uuid.MustParse(accessTokenClaims.UserUuid) {
				return h.ErrorResponse(c, http.StatusForbidden, "access denied", err)
			}
			projectFilters, err := h.filterRepo.FindByProjectToken(result.Token)
			if err == nil {
				filters = append(filters, *projectFilters)
			} else {
				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters for project", err)
			}
			if projectFilters == nil {
				return h.ErrorResponse(c, http.StatusNotFound, "filters not found for project", nil)
			}
		}
	case "staff":
		if userIdStr != "" {
			result, err := RPCClient.GetProjectList(userUuid)
			if err != nil {
				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
			}
			if len(*result) == 0 {
				return h.ErrorResponse(c, http.StatusNotFound, "projects not found for user", nil)
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
					return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
				}
				if projectFilters, err := h.filterRepo.FindByProjectToken(result.Token); err == nil {
					filters = append(filters, *projectFilters)
				} else {
					return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters for project", err)
				}
			} else {
				allFilters, err := h.filterRepo.FindAll()
				if err != nil {
					return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
				}
				filters = *allFilters
			}
		}

	default:
		return h.ErrorResponse(c, http.StatusForbidden, "access denied", nil)
	}
	if len(filters) == 0 {
		return h.ErrorResponse(c, http.StatusNotFound, "filters not found", nil)
	}
	return h.SuccessResponse(c, http.StatusOK, "filters found", filters)
}

// UpdateFilter godoc
// @Summary Update Filter
// @Description Update Filter for Project by projectId. Roles allowed: business, staff
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
func (h *BaseHandler) UpdateFilter(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*OAuthService.AccessTokenClaims)
	filterId, err := strconv.Atoi(c.Param("filterId"))
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "filter id does not valid", err)
	}
	var requestPayload FilterUpdateRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	}

	if err := requestPayload.Validate(); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	}

	filter, err := h.filterRepo.FindById(filterId)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filter", err)
	}
	if filter == nil {
		return h.ErrorResponse(c, http.StatusNotFound, "filter not found", nil)
	}

	switch accessTokenClaims.UserRole {
	case "business":
		result, err := RPCClient.GetProject(filter.ProjectToken)
		if err != nil {
			return h.ErrorResponse(c, http.StatusInternalServerError, "could not update filter", err)
		}
		if result.Token != filter.ProjectToken {
			return h.ErrorResponse(c, http.StatusForbidden, "access denied", err)
		}
	case "staff":
		break
	default:
		return h.ErrorResponse(c, http.StatusForbidden, "access denied", nil)
	}

	filter.Type = requestPayload.Type
	filter.Coverage = requestPayload.Coverage

	if err := h.filterRepo.Update(filter); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not update filter", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "filter successfully updated", filter)
}

// DeleteFilter godoc
// @Summary Delete Filter
// @Description Delete Filter for Project by projectId. Roles allowed: business, staff
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
func (h *BaseHandler) DeleteFilter(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*OAuthService.AccessTokenClaims)
	filterId, err := strconv.Atoi(c.Param("filterId"))
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "filter id does not valid", err)
	}

	filter, err := h.filterRepo.FindById(filterId)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filter", err)
	}
	if filter == nil {
		return h.ErrorResponse(c, http.StatusNotFound, "filter not found", nil)
	}

	switch accessTokenClaims.UserRole {
	case "business":
		result, err := RPCClient.GetProject(filter.ProjectToken)
		if err != nil {
			return h.ErrorResponse(c, http.StatusInternalServerError, "could not delete filter", err)
		}
		if result.Token != filter.ProjectToken {
			return h.ErrorResponse(c, http.StatusForbidden, "access denied", err)
		}
	case "staff":
		break
	default:
		return h.ErrorResponse(c, http.StatusForbidden, "access denied", nil)
	}
	if err := h.filterRepo.Delete(filter); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not delete filter", err)
	}
	return h.SuccessResponse(c, http.StatusNoContent, "filter successfully deleted", nil)
}
