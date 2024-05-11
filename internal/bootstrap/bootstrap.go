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

// Bootstrap struct holds all the components necessary for the application to run.
// It contains configuration settings and the main application router.
type Bootstrap struct {
	Config *configuration.AppConfig // Configuration settings for the application
	Router *gin.Engine              // HTTP router for handling web requests
}

// InitBootstrap initializes the bootstrap process.
// It's a public wrapper around the private initBootstrap function for initializing the Bootstrap struct.
func InitBootstrap() Bootstrap {
	return initBootstrap()
}

// initBootstrap handles the actual initialization logic for the Bootstrap struct.
// It includes setting up configurations, database connections, middleware, and routes.
// It panics if the configuration is not set.
func initBootstrap() Bootstrap {
	if configuration.Config == nil {
		log.Panic().Msg("configuration is nil")
	}

	app := Bootstrap{}
	app.Config = configuration.Config

	// Initialize MySQL client
	mysqlClient := mysql.NewClient(app.Config)

	// Initialize MySQL adapters for different domain models
	wallAdapter := mysql.NewWallAdapter(mysqlClient)
	blockAdapter := mysql.NewBlockAdapter(mysqlClient)
	programAdapter := mysql.NewProgramAdapter(mysqlClient)
	episodeAdapter := mysql.NewEpisodeAdapter(mysqlClient)
	mediaAdapter := mysql.NewMediaAdapter(mysqlClient)
	tagAdapter := mysql.NewTagAdapter(mysqlClient)
	catagoryAdapter := mysql.NewCategoryAdapter(mysqlClient)

	// Init apis
	wallApi := api.NewWallApi(wallAdapter)
	blockApi := api.NewBlockApi(blockAdapter)
	programApi := api.NewProgramApi(programAdapter)
	episodeApi := api.NewEpisodeApi(episodeAdapter)
	mediaApi := api.NewMediaApi(mediaAdapter)
	tagApi := api.NewTagApi(tagAdapter)
	catApi := api.NewCategoryApi(catagoryAdapter)

	// Init handlers
	wallHandler := handlers.NewWallHandler(wallApi)
	blockHandler := handlers.NewBlockHandler(blockApi)
	programHandler := handlers.NewProgramHandler(programApi)
	episodeHandler := handlers.NewEpisodeHandler(episodeApi)
	mediaHanlder := handlers.NewMediaHandler(mediaApi)
	tagHanlder := handlers.NewTagHandler(tagApi)
	catHanlder := handlers.NewCategoryHandler(catApi)

	r := router.CreateRouter(
		wallHandler,
		blockHandler,
		programHandler,
		episodeHandler,
		mediaHanlder,
		tagHanlder,
		catHanlder,
	)
	app.Router = r
	return app
}

// Run starts the application by running the HTTP server.
// It logs a fatal error if the server cannot be started.
func (b Bootstrap) Run() {
	dsn := fmt.Sprintf("%s:%d", b.Config.HostAddress, b.Config.HostPort)
	if errRun := b.Router.Run(dsn); errRun != nil {
		log.Fatal().Msg("error during service instantiation")
	}
}
