package main

import (
	"github.com/khedhrije/podcaster-backoffice-api/internal/bootstrap"
)

// @title           podcaster-backoffice-api
// @version         1.0.0
// @description     This is the documentation for the podcaster-backoffice-api.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.email   khedhri.je@gmail.com
// @host      		localhost:8080
//
// @securityDefinitions.apikey Bearer-JWT
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a valid JWT token.
//
// @securityDefinitions.apikey Bearer-APIKey
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a valid API key.
func main() {
	// Initialize the bootstrap process, which sets up the application
	bootstrap.InitBootstrap().Run()
}
