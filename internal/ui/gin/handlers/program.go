// Package handlers provides HTTP request handlers for managing programs.
package handlers

import (
	"github.com/khedhrije/podcaster-backoffice-api/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Program represents the interface for managing programs.
type Program interface {
	// Create returns a Gin handler function for creating a program.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a program.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a program by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all programs.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a program by its UUID.
	Delete() gin.HandlerFunc

	// FindEpisodes returns a Gin handler function for finding a program's episodes.
	FindEpisodes() gin.HandlerFunc

	// FindTags returns a Gin handler function for finding a program's tags.
	FindTags() gin.HandlerFunc

	// FindCategories returns a Gin handler function for finding a program's categories.
	FindCategories() gin.HandlerFunc

	// OverwriteCategories returns a Gin handler function for overwriting categories of a program.
	OverwriteCategories() gin.HandlerFunc

	// OverwriteTags returns a Gin handler function for overwriting tags of a program.
	OverwriteTags() gin.HandlerFunc
}

type programHandler struct {
	api api.Program
}

// NewProgramHandler creates a new instance of Program interface.
func NewProgramHandler(api api.Program) Program {
	return &programHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a program.
//
// @Summary Create a new program
// @Description Create a new program
// @Tags programs
// @ID create-program
// @Param request body pkg.CreateProgramRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs [post]
func (handler programHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest pkg.CreateProgramRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create program
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating program: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

// Update returns a Gin handler function for updating a program.
//
// @Summary Update program
// @Description Update program
// @Tags programs
// @ID update-program
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateProgramRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid} [put]
func (handler programHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest pkg.UpdateProgramRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update program
		if err := handler.api.Update(c, programUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating program: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

// Find returns a Gin handler function for finding a program by its UUID.
//
// @Summary Find a program
// @Description Find a program
// @Tags programs
// @ID find-program
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {object} pkg.ProgramResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid} [get]
func (handler programHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program
		program, err := handler.api.Find(c, programUUID)
		if err != nil {
			log.Error().Msg("error finding program: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, program)
	}
}

// FindAll returns a Gin handler function for finding all programs.
//
// @Summary Find all programs
// @Description Find all programs
// @Tags programs
// @ID find-all-programs
// @Produce json
// @Success 200 {array} pkg.ProgramResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs [get]
func (handler programHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all programs
		programs, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all programs: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, programs)
	}
}

// Delete returns a Gin handler function for deleting a program by its UUID.
//
// @Summary Delete a program
// @Description Delete a program
// @Tags programs
// @ID delete-program
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "deleted"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid} [delete]
func (handler programHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to delete program
		if err := handler.api.Delete(c, programUUID); err != nil {
			log.Error().Msg("error deleting program: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}

// FindEpisodes returns a Gin handler function for finding a program's episodes.
//
// @Summary Find a program's episodes
// @Description Find a program's episodes
// @Tags programs
// @ID find-program-episodes
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {array} pkg.EpisodeResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid}/episodes [get]
func (handler programHandler) FindEpisodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program's episodes
		episodes, err := handler.api.FindEpisodes(c, programUUID)
		if err != nil {
			log.Error().Msg("error finding program's episodes: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, episodes)
	}
}

// FindTags returns a Gin handler function for finding a program's tags.
//
// @Summary Find a program's tags
// @Description Find a program's tags
// @Tags programs
// @ID find-program-tags
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {array} pkg.TagResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid}/tags [get]
func (handler programHandler) FindTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program's tags
		tags, err := handler.api.FindTags(c, programUUID)
		if err != nil {
			log.Error().Msg("error finding program's tags: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, tags)
	}
}

// FindCategories returns a Gin handler function for finding a program's categories.
//
// @Summary Find a program's categories
// @Description Find a program's categories
// @Tags programs
// @ID find-program-categories
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {array} pkg.CategoryResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid}/categories [get]
func (handler programHandler) FindCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program's categories
		categories, err := handler.api.FindCats(c, programUUID)
		if err != nil {
			log.Error().Msg("error finding program's categories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, categories)
	}
}

// OverwriteCategories returns a Gin handler function for overwriting categories.
//
// @Summary Overwrite categories of a program
// @Description Overwrite the categories of a specific program by replacing all existing categories with new ones
// @Tags programs
// @ID overwrite-program-categories
// @Param uuid path string true "UUID of the program"
// @Param request body []string true "List of categories' UUIDs to set"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 422 {object} pkg.ErrorJSON "Unprocessable Entity"
// @Failure 500 {object} pkg.ErrorJSON "Internal Server Error"
// @Router /private/programs/{uuid}/categories/overwrite [put]
func (handler programHandler) OverwriteCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest []string
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to overwrite categories
		if err := handler.api.OverwriteCategories(c, programUUID, jsonRequest); err != nil {
			log.Error().Msg("error overwriting program's categories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "ok")
	}
}

// OverwriteTags returns a Gin handler function for overwriting tags.
//
// @Summary Overwrite tags of a program
// @Description Overwrite the tags of a specific program by replacing all existing tags with new ones
// @Tags programs
// @ID overwrite-program-tags
// @Param uuid path string true "UUID of the program"
// @Param request body []string true "List of tags UUIDs to set"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 422 {object} pkg.ErrorJSON "Unprocessable Entity"
// @Failure 500 {object} pkg.ErrorJSON "Internal Server Error"
// @Router /private/programs/{uuid}/tags/overwrite [put]
func (handler programHandler) OverwriteTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest []string
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to overwrite tags
		if err := handler.api.OverwriteTags(c, programUUID, jsonRequest); err != nil {
			log.Error().Msg("error overwriting program's tags: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "ok")
	}
}
