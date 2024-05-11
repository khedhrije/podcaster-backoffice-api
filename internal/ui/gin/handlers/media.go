// Package handlers provides HTTP request handlers for managing medias.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Media represents the interface for managing medias.
type Media interface {
	// Create returns a Gin handler function for creating a media.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a media.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a media by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all medias.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a media by its UUID.
	Delete() gin.HandlerFunc
}

type mediaHandler struct {
	api api.Media
}

// NewMediaHandler creates a new instance of Media interface.
func NewMediaHandler(api api.Media) Media {
	return &mediaHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a media.
//
// @Summary Create a new media
// @Description Create a new media
// @Tags medias
// @ID create-media
// @Param request body pkg.CreateMediaRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/medias [post]
func (handler mediaHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest api.CreateMediaRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create media
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating media: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Update returns a Gin handler function for updating a media.
//
// @Summary Update media
// @Description Update media
// @Tags medias
// @ID update-media
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateMediaRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/medias/{uuid} [put]
func (handler mediaHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract media UUID from path
		mediaUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest api.CreateMediaRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update media
		if err := handler.api.Update(c, mediaUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating media: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Find returns a Gin handler function for finding a media by its UUID.
//
// @Summary Find a media
// @Description Find a media
// @Tags medias
// @ID find-media
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/medias/{uuid} [get]
func (handler mediaHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract media UUID from path
		mediaUUID := c.Param("uuid")

		// Call API to find media
		media, err := handler.api.Find(c, mediaUUID)
		if err != nil {
			log.Error().Msg("error finding media: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, media)
	}
}

// FindAll returns a Gin handler function for finding all medias.
//
// @Summary Find all medias
// @Description Find all medias
// @Tags medias
// @ID find-all-medias
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/medias/ [get]
func (handler mediaHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all medias
		medias, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all medias: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, medias)
	}
}

// Delete returns a Gin handler function for deleting a media by its UUID.
//
// @Summary Delete a media
// @Description Delete a media
// @Tags medias
// @ID delete-media
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/medias/{uuid} [delete]
func (handler mediaHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract media UUID from path
		mediaUUID := c.Param("uuid")

		// Call API to delete media
		if err := handler.api.Delete(c, mediaUUID); err != nil {
			log.Error().Msg("error deleting media: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}
