package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/common/validators"
	"github.com/labstack/echo/v4"
)

type UpdateFilterRequest struct {
	Type     string `json:"type" example:"whitelist"`
	Coverage string `json:"coverage" example:"equals"`
}

func (ur *UpdateFilterRequest) Validate() error {
	if err := validators.ValidateDomainTypes(ur.Type); err != nil {
		return err
	}
	if err := validators.ValidateDomainCoverage(ur.Coverage); err != nil {
		return err
	}
	return nil
}

// UpdateFilter TODO: refactor this

// UpdateFilter godoc
// @Summary Update Filter
// @Description Update Filter for ProjectRPCPayload by projectId. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int true "filter id"
// @Param filter body UpdateFilterRequest true "raw request body"
// @Success 200 {object} Filter
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/filters/{filterId} [put]
func (fh FilterHandler) UpdateFilter(c echo.Context) error {
	return nil
	//accessTokenClaims, _ := c.Get("accessTokenClaims").(*OAuthService.AccessTokenClaims)
	//filterId, err := strconv.Atoi(c.Param("filterId"))
	//if err != nil {
	//	return h.ErrorResponse(c, http.StatusBadRequest, "filter id does not valid", err)
	//}
	//var requestPayload UpdateFilterRequest
	//if err := c.Bind(&requestPayload); err != nil {
	//	return h.ErrorResponse(c, http.StatusUnprocessableEntity, "could not read request body", err)
	//}
	//
	//if err := requestPayload.Validate(); err != nil {
	//	return h.ErrorResponse(c, http.StatusBadRequest, err.Message, err.Error())
	//}
	//
	//filter, err := h.filterRepo.FindById(filterId)
	//if err != nil {
	//	return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filter", err)
	//}
	//if filter == nil {
	//	return h.ErrorResponse(c, http.StatusNotFound, "filter not found", nil)
	//}
	//
	//switch accessTokenClaims.UserRole {
	//case "business":
	//	result, err := RpcClient.GetProject(filter.ProjectToken)
	//	if err != nil {
	//		return h.ErrorResponse(c, http.StatusInternalServerError, "could not update filter", err)
	//	}
	//	if result.Token != filter.ProjectToken {
	//		return h.ErrorResponse(c, http.StatusForbidden, "access denied", err)
	//	}
	//case "staff":
	//	break
	//default:
	//	return h.ErrorResponse(c, http.StatusForbidden, "access denied", nil)
	//}
	//
	//filter.Type = requestPayload.Type
	//filter.Match = requestPayload.Match
	//
	//if err := h.filterRepo.Update(filter); err != nil {
	//	return h.ErrorResponse(c, http.StatusInternalServerError, "could not update filter", err)
	//}
	//return h.SuccessResponse(c, http.StatusOK, "filter successfully updated", filter)
}
