// Package handlers provides HTTP request handlers for managing tags.
package handlers

import (
	"github.com/khedhrije/podcaster-backoffice-api/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Tag represents the interface for managing tags.
type Tag interface {
	// Create returns a Gin handler function for creating a tag.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a tag.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a tag by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all tags.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a tag by its UUID.
	Delete() gin.HandlerFunc

	FindPrograms() gin.HandlerFunc
}

type tagHandler struct {
	api api.Tag
}

// NewTagHandler creates a new instance of Tag interface.
func NewTagHandler(api api.Tag) Tag {
	return &tagHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a tag.
//
// @Summary Create a new tag
// @Description Create a new tag
// @Tags tags
// @ID create-tag
// @Param request body pkg.CreateTagRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/tags [post]
func (handler tagHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest pkg.CreateTagRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create tag
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating tag: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Update returns a Gin handler function for updating a tag.
//
// @Summary Update tag
// @Description Update tag
// @Tags tags
// @ID update-tag
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateTagRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/tags/{uuid} [put]
func (handler tagHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract tag UUID from path
		tagUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest pkg.UpdateTagRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update tag
		if err := handler.api.Update(c, tagUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating tag: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Find returns a Gin handler function for finding a tag by its UUID.
//
// @Summary Find a tag
// @Description Find a tag
// @Tags tags
// @ID find-tag
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/tags/{uuid} [get]
func (handler tagHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract tag UUID from path
		tagUUID := c.Param("uuid")

		// Call API to find tag
		tag, err := handler.api.Find(c, tagUUID)
		if err != nil {
			log.Error().Msg("error finding tag: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, tag)
	}
}

// FindAll returns a Gin handler function for finding all tags.
//
// @Summary Find all tags
// @Description Find all tags
// @Tags tags
// @ID find-all-tags
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/tags/ [get]
func (handler tagHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all tags
		tags, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all tags: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, tags)
	}
}

// Delete returns a Gin handler function for deleting a tag by its UUID.
//
// @Summary Delete a tag
// @Description Delete a tag
// @Tags tags
// @ID delete-tag
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/tags/{uuid} [delete]
func (handler tagHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract tag UUID from path
		tagUUID := c.Param("uuid")

		// Call API to delete tag
		if err := handler.api.Delete(c, tagUUID); err != nil {
			log.Error().Msg("error deleting tag: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}

// FindPrograms returns a Gin handler function for finding all block's programs.
//
// @Summary Find all tag's programs
// @Description Find all tag's programs
// @Tags tags
// @ID find-tag-programs
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/tags/{uuid}/programs [get]
func (handler tagHandler) FindPrograms() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Extract tag UUID from path
		tagUUID := c.Param("uuid")

		// Call API to find all walls
		programs, err := handler.api.FindPrograms(c, tagUUID)
		if err != nil {
			log.Error().Msg("error finding all tag's programs: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, programs)
	}
}
