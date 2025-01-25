package HTTPServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
)

func (s *Server) setupRoutes() {
	s.echo.POST("/v1/data/inspect", s.checkHandler.Inspect, s.apiKeyAuthMiddleware.Auth())

	s.echo.POST("/v1/access", s.accessHandler.CreateAccess)

	s.echo.POST("/v1/domains/count", s.domainHandler.Count)

	s.echo.POST("/v1/reviews", s.reviewHandler.CreateReview)

	s.echo.GET("/v1/filters", s.filterHandler.GetFilterList, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))
	s.echo.POST("/v1/filters", s.filterHandler.CreateFilter, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))
	s.echo.DELETE("/v1/filters/:domain_name", s.filterHandler.DeleteFilter, s.firebaseAuthMiddleware.RoleBasedAuth(models.CustomerRole, models.StaffRole))

	s.echo.GET("/v1/domains/:domain_name", s.domainHandler.GetDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
	s.echo.POST("/v1/domains", s.domainHandler.CreateDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
	s.echo.PATCH("/v1/domains/:domain_name", s.domainHandler.UpdateDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
	s.echo.DELETE("/v1/domains/:domain_name", s.domainHandler.DeleteDomain, s.firebaseAuthMiddleware.RoleBasedAuth(models.StaffRole))
}
