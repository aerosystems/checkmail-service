package HttpServer

import (
	"github.com/aerosystems/checkmail-service/internal/models"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) setupRoutes() {
	// Private routes Basic Auth
	docsGroup := s.echo.Group("/docs")
	docsGroup.GET("/*", echoSwagger.WrapHandler)

	// Auth X-Api-Key and reCAPTCHA implemented on API Gateway
	s.echo.POST("/v1/data/inspect", s.inspectHandler.Inspect)

	// Protected with reCAPTCHA on API Gateway
	s.echo.POST("/v1/domains/count", s.domainHandler.Count)
	s.echo.POST("/v1/reviews", s.reviewHandler.CreateReview)

	// Private routes OAuth 2.0: check roles [customer, staff]. Auth implemented on API Gateway
	s.echo.GET("/v1/filters", s.filterHandler.GetFilterList, s.AuthTokenMiddleware(models.CustomerRole, models.StaffRole))
	s.echo.POST("/v1/filters", s.filterHandler.CreateFilter, s.AuthTokenMiddleware(models.CustomerRole, models.StaffRole))
	s.echo.PUT("/v1/filters/:filterId", s.filterHandler.UpdateFilter, s.AuthTokenMiddleware(models.CustomerRole, models.StaffRole))
	s.echo.DELETE("/v1/filters/:filterId", s.filterHandler.DeleteFilter, s.AuthTokenMiddleware(models.CustomerRole, models.StaffRole))

	// Private routes OAuth 2.0: check roles [staff]. Auth implemented on API Gateway
	s.echo.GET("/v1/domains/:domainName", s.domainHandler.GetDomain, s.AuthTokenMiddleware(models.StaffRole))
	s.echo.POST("/v1/domains", s.domainHandler.CreateDomain, s.AuthTokenMiddleware(models.StaffRole))
	s.echo.PATCH("/v1/domains/:domainName", s.domainHandler.UpdateDomain, s.AuthTokenMiddleware(models.StaffRole))
	s.echo.DELETE("/v1/domains/:domainName", s.domainHandler.DeleteDomain, s.AuthTokenMiddleware(models.StaffRole))
}
