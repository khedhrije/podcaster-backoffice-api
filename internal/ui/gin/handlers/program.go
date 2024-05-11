// Package handlers provides HTTP request handlers for managing programs.
package handlers

import (
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

	FindEpisodes() gin.HandlerFunc
	FindTags() gin.HandlerFunc
	FindCategories() gin.HandlerFunc
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
		var jsonRequest api.CreateProgramRequest
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
		var jsonRequest api.CreateProgramRequest
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
// @Success 200 {string} string "ok"
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
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/ [get]
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
// @Success 200 {string} string "ok"
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
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid}/episodes [get]
func (handler programHandler) FindEpisodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program
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
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid}/tags [get]
func (handler programHandler) FindTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program
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
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/programs/{uuid}/categories [get]
func (handler programHandler) FindCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract program UUID from path
		programUUID := c.Param("uuid")

		// Call API to find program
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
