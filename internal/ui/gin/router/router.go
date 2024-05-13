// Package router sets up the routes for the API using the Gin framework.
package router

import (
	"github.com/gin-gonic/gin"
	spec "github.com/khedhrije/podcaster-backoffice-api/deployments/swagger"
	"github.com/khedhrije/podcaster-backoffice-api/internal/configuration"
	"github.com/khedhrije/podcaster-backoffice-api/internal/ui/gin/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
)

// CreateRouter sets up and returns a new Gin router with the defined routes.
func CreateRouter(wall handlers.Wall, block handlers.Block, program handlers.Program, episode handlers.Episode, media handlers.Media, tag handlers.Tag, category handlers.Category) *gin.Engine {
	// Initialize a new Gin router without any middleware by default.
	r := gin.New()

	// Configure Swagger documentation URL based on the environment.
	if configuration.Config.Env == "dev" {
		spec.SwaggerInfo.Host = configuration.Config.DocsAddress
	} else {
		spec.SwaggerInfo.Host = removeHTTPS(configuration.Config.DocsAddress)
	}

	// Set up the route for Swagger documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Define private routes that require authentication.
	private := r.Group("/private")
	{
		// Routes for managing walls.
		walls := private.Group("/walls")
		{
			walls.POST("", wall.Create())
			walls.PUT("/:uuid", wall.Update())
			walls.GET("/:uuid", wall.Find())
			walls.GET("", wall.FindAll())
			walls.DELETE("/:uuid", wall.Delete())
			walls.GET("/:uuid/blocks", wall.FindBlocks())
			walls.PUT("/:uuid/blocks/overwrite", wall.OverwriteBlocks())
		}

		// Routes for managing blocks.
		blocks := private.Group("/blocks")
		{
			blocks.POST("", block.Create())
			blocks.PUT("/:uuid", block.Update())
			blocks.GET("/:uuid", block.Find())
			blocks.GET("", block.FindAll())
			blocks.DELETE("/:uuid", block.Delete())
			blocks.GET("/:uuid/programs", block.FindPrograms())
			blocks.PUT("/:uuid/programs/overwrite", block.OverwritePrograms())
		}

		// Routes for managing programs.
		programs := private.Group("/programs")
		{
			programs.POST("", program.Create())
			programs.PUT("/:uuid", program.Update())
			programs.GET("/:uuid", program.Find())
			programs.GET("", program.FindAll())
			programs.DELETE("/:uuid", program.Delete())
			programs.GET("/:uuid/episodes", program.FindEpisodes())
			programs.GET("/:uuid/tags", program.FindTags())
			programs.GET("/:uuid/categories", program.FindCategories())
			programs.PUT("/:uuid/tags/overwrite", program.OverwriteTags())
			programs.PUT("/:uuid/categories/overwrite", program.OverwriteCategories())
		}

		// Routes for managing episodes.
		episodes := private.Group("/episodes")
		{
			episodes.POST("", episode.Create())
			episodes.PUT("/:uuid", episode.Update())
			episodes.GET("/:uuid", episode.Find())
			episodes.GET("", episode.FindAll())
			episodes.DELETE("/:uuid", episode.Delete())
		}

		// Routes for managing media.
		mediaRoutes := private.Group("/medias")
		{
			mediaRoutes.POST("", media.Create())
			mediaRoutes.PUT("/:uuid", media.Update())
			mediaRoutes.GET("/:uuid", media.Find())
			mediaRoutes.GET("", media.FindAll())
			mediaRoutes.DELETE("/:uuid", media.Delete())
		}

		// Routes for managing tags.
		tags := private.Group("/tags")
		{
			tags.POST("", tag.Create())
			tags.PUT("/:uuid", tag.Update())
			tags.GET("/:uuid", tag.Find())
			tags.GET("", tag.FindAll())
			tags.DELETE("/:uuid", tag.Delete())
			tags.GET("/:uuid/programs", tag.FindPrograms())
		}

		// Routes for managing categories.
		categories := private.Group("/categories")
		{
			categories.POST("", category.Create())
			categories.PUT("/:uuid", category.Update())
			categories.GET("/:uuid", category.Find())
			categories.GET("", category.FindAll())
			categories.DELETE("/:uuid", category.Delete())
			categories.GET("/:uuid/programs", category.FindPrograms())
		}
	}

	// Return the configured router.
	return r
}

// removeHTTPS removes the "https://" prefix from a URL string.
func removeHTTPS(url string) string {
	// Check if the URL starts with "https://"
	if strings.HasPrefix(url, "https://") {
		// Remove "https://" from the URL
		return strings.TrimPrefix(url, "https://")
	}
	// If the URL doesn't start with "https://", return it as is
	return url
}
