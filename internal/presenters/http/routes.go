package HttpServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

func (s *Server) setupRoutes() {
	// Auth X-Api-Key
	s.echo.POST("/v1/data/inspect", s.checkHandler.Inspect, s.apiKeyAuthMiddleware.Auth())

	// PubSub
	s.echo.POST("/v1/access", s.accessHandler.CreateAccess)

	// Protected with reCAPTCHA on API Gateway
	s.echo.POST("/v1/domains/count", s.domainHandler.Count)
	s.echo.POST("/v1/reviews", s.reviewHandler.CreateReview)

	// Private routes OAuth 2.0: check roles [customer, staff]. Auth implemented on API Gateway
	s.echo.GET("/v1/filters", s.filterHandler.GetFilterList, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))
	s.echo.POST("/v1/filters", s.filterHandler.CreateFilter, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))
	s.echo.PUT("/v1/filters/:filterId", s.filterHandler.UpdateFilter, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))
	s.echo.DELETE("/v1/filters/:filterId", s.filterHandler.DeleteFilter, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))

	// Private routes OAuth 2.0: check roles [staff]. Auth implemented on API Gateway
	s.echo.GET("/v1/domains/:domainName", s.domainHandler.GetDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
	s.echo.POST("/v1/domains", s.domainHandler.CreateDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
	s.echo.PATCH("/v1/domains/:domainName", s.domainHandler.UpdateDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
	s.echo.DELETE("/v1/domains/:domainName", s.domainHandler.DeleteDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
}
