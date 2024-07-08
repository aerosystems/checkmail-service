package filter

import "github.com/labstack/echo/v4"

// GetFilter TODO: refactor this

// GetFilterList godoc
// @Summary Get Filter List
// @Description Get Filter List for all user Projects. Roles allowed: business, staff
// @Tags Filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param userId query int false "user id"
// @Param projectToken query string false "project token"
// @Success 200 {object} Filter
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/filters [get]
func (fh Handler) GetFilterList(c echo.Context) error {
	return nil
	//accessTokenClaims, _ := c.Get("accessTokenClaims").(*OAuthService.AccessTokenClaims)
	//var filters []models.Filter
	//var userIdStr, projectToken string
	//var err error
	//var userUuid uuid.UUID
	//userIdStr = c.QueryParam("userId")
	//if userIdStr != "" {
	//	userUuid, err = uuid.Parse(userIdStr)
	//	if err != nil {
	//		return h.ErrorResponse(c, http.StatusBadRequest, "user id does not valid", err)
	//	}
	//}
	//projectToken = c.QueryParam("projectToken")
	//switch accessTokenClaims.UserRole {
	//case "business":
	//	if projectToken == "" {
	//		result, err := RpcClient.GetProjectList(uuid.MustParse(accessTokenClaims.UserUuid))
	//		if err != nil {
	//			return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
	//		}
	//		if len(*result) == 0 {
	//			return h.ErrorResponse(c, http.StatusNotFound, "projects not found for user", nil)
	//		}
	//		for _, project := range *result {
	//			if projectFilters, err := h.filterRepo.FindByProjectToken(project.Token); err == nil {
	//				filters = append(filters, *projectFilters)
	//			} else {
	//				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters for project", err)
	//			}
	//		}
	//	} else {
	//		result, err := RpcClient.GetProject(projectToken)
	//		if err != nil {
	//			return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
	//		}
	//		if result.UserUuid != uuid.MustParse(accessTokenClaims.UserUuid) {
	//			return h.ErrorResponse(c, http.StatusForbidden, "access denied", err)
	//		}
	//		projectFilters, err := h.filterRepo.FindByProjectToken(result.Token)
	//		if err == nil {
	//			filters = append(filters, *projectFilters)
	//		} else {
	//			return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters for project", err)
	//		}
	//		if projectFilters == nil {
	//			return h.ErrorResponse(c, http.StatusNotFound, "filters not found for project", nil)
	//		}
	//	}
	//case "staff":
	//	if userIdStr != "" {
	//		result, err := RpcClient.GetProjectList(userUuid)
	//		if err != nil {
	//			return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
	//		}
	//		if len(*result) == 0 {
	//			return h.ErrorResponse(c, http.StatusNotFound, "projects not found for user", nil)
	//		}
	//		for _, project := range *result {
	//			if projectToken != "" {
	//				if project.Token != projectToken {
	//					continue
	//				}
	//			}
	//			projectFilters, err := h.filterRepo.FindByProjectToken(project.Token)
	//			if err != nil {
	//				continue
	//			}
	//			h.log.Info(projectFilters)
	//			filters = append(filters, *projectFilters)
	//		}
	//	} else {
	//		if projectToken != "" {
	//			result, err := RpcClient.GetProject(projectToken)
	//			if err != nil {
	//				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
	//			}
	//			if projectFilters, err := h.filterRepo.FindByProjectToken(result.Token); err == nil {
	//				filters = append(filters, *projectFilters)
	//			} else {
	//				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters for project", err)
	//			}
	//		} else {
	//			allFilters, err := h.filterRepo.FindAll()
	//			if err != nil {
	//				return h.ErrorResponse(c, http.StatusInternalServerError, "could not find filters", err)
	//			}
	//			filters = *allFilters
	//		}
	//	}
	//
	//default:
	//	return h.ErrorResponse(c, http.StatusForbidden, "access denied", nil)
	//}
	//if len(filters) == 0 {
	//	return h.ErrorResponse(c, http.StatusNotFound, "filters not found", nil)
	//}
	//return h.SuccessResponse(c, http.StatusOK, "filters found", filters)
}
