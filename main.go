package main

import (
	"fmt"
	"github.com/khedhrije/podcaster-backoffice-api/internal/bootstrap"
	"github.com/khedhrije/podcaster-backoffice-api/internal/configuration"
	"github.com/rs/zerolog/log"
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
	log.Info().
		Interface("app", configuration.Config.Name).
		Interface("env", configuration.Config.Env).
		Interface("address", fmt.Sprintf("%s:%v", configuration.Config.HostAddress, configuration.Config.HostPort)).
		Msg("app is ready")
	bootstrap.InitBootstrap().Run()
}
