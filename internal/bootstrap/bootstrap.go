package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/configuration"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/khedhrije/podcaster-backoffice-api/internal/infrastructure/mysql"
	"github.com/khedhrije/podcaster-backoffice-api/internal/ui/gin/handlers"
	"github.com/khedhrije/podcaster-backoffice-api/internal/ui/gin/router"
	"github.com/rs/zerolog/log"
)

// Bootstrap struct encapsulates the configuration settings and the HTTP router necessary for the application to run.
type Bootstrap struct {
	Config *configuration.AppConfig // Application configuration settings
	Router *gin.Engine              // HTTP router for handling web requests
}

// InitBootstrap initializes the bootstrap process and returns a Bootstrap instance.
// It serves as a public entry point for the initialization process.
func InitBootstrap() Bootstrap {
	return initBootstrap()
}

// initBootstrap sets up the Bootstrap struct by initializing configurations, database connections, middleware, and routes.
// It panics if the configuration is not set, ensuring that the application does not run with nil configurations.
func initBootstrap() Bootstrap {
	if configuration.Config == nil {
		log.Panic().Msg("configuration is nil")
	}

	app := Bootstrap{}
	app.Config = configuration.Config

	// Initialize MySQL client using application configuration
	mysqlClient := mysql.NewClient(app.Config)

	// Initialize MySQL adapters for different domain models, setting up the data access layer
	wallAdapter := mysql.NewWallAdapter(mysqlClient)
	wallBlockAdapter := mysql.NewWallBlockAdapter(mysqlClient)
	blockAdapter := mysql.NewBlockAdapter(mysqlClient)
	blockProgramAdapter := mysql.NewBlockProgramAdapter(mysqlClient)
	programAdapter := mysql.NewProgramAdapter(mysqlClient)
	episodeAdapter := mysql.NewEpisodeAdapter(mysqlClient)
	mediaAdapter := mysql.NewMediaAdapter(mysqlClient)
	tagAdapter := mysql.NewTagAdapter(mysqlClient)
	programTagAdapter := mysql.NewProgramTagAdapter(mysqlClient)
	categoryAdapter := mysql.NewCategoryAdapter(mysqlClient)
	programCategoryAdapter := mysql.NewProgramCategoryAdapter(mysqlClient)

	// Initialize APIs for different domain models, enabling business logic operations
	wallApi := api.NewWallApi(wallAdapter, wallBlockAdapter, blockAdapter)
	blockApi := api.NewBlockApi(blockAdapter, blockProgramAdapter, programAdapter)
	programApi := api.NewProgramApi(programAdapter, episodeAdapter, programTagAdapter, tagAdapter, programCategoryAdapter, categoryAdapter)
	episodeApi := api.NewEpisodeApi(episodeAdapter)
	mediaApi := api.NewMediaApi(mediaAdapter)
	tagApi := api.NewTagApi(tagAdapter, programTagAdapter, programAdapter)
	catApi := api.NewCategoryApi(categoryAdapter, programCategoryAdapter, programAdapter)

	// Initialize handlers for different APIs, setting up the presentation layer
	wallHandler := handlers.NewWallHandler(wallApi)
	blockHandler := handlers.NewBlockHandler(blockApi)
	programHandler := handlers.NewProgramHandler(programApi)
	episodeHandler := handlers.NewEpisodeHandler(episodeApi)
	mediaHandler := handlers.NewMediaHandler(mediaApi)
	tagHandler := handlers.NewTagHandler(tagApi)
	catHandler := handlers.NewCategoryHandler(catApi)

	// Create the router with the initialized handlers, configuring the request handling
	r := router.CreateRouter(
		wallHandler,
		blockHandler,
		programHandler,
		episodeHandler,
		mediaHandler,
		tagHandler,
		catHandler,
	)
	app.Router = r
	return app
}

// Run starts the application by running the HTTP server on the configured host address and port.
// It logs a fatal error if the server cannot be started, ensuring that the failure is captured and reported.
func (b Bootstrap) Run() {
	dsn := fmt.Sprintf("%s:%d", b.Config.HostAddress, b.Config.HostPort)
	if errRun := b.Router.Run(dsn); errRun != nil {
		log.Fatal().Msg("error during service instantiation")
	}
}
