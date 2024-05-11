// Package handlers provides HTTP request handlers for managing episodes.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Episode represents the interface for managing episodes.
type Episode interface {
	// Create returns a Gin handler function for creating a episode.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a episode.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a episode by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all episodes.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a episode by its UUID.
	Delete() gin.HandlerFunc
}

type episodeHandler struct {
	api api.Episode
}

// NewEpisodeHandler creates a new instance of Episode interface.
func NewEpisodeHandler(api api.Episode) Episode {
	return &episodeHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a episode.
//
// @Summary Create a new episode
// @Description Create a new episode
// @Tags episodes
// @ID create-episode
// @Param request body pkg.CreateEpisodeRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/episodes [post]
func (handler episodeHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest api.CreateEpisodeRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create episode
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating episode: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Update returns a Gin handler function for updating a episode.
//
// @Summary Update episode
// @Description Update episode
// @Tags episodes
// @ID update-episode
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateEpisodeRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/episodes/{uuid} [put]
func (handler episodeHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract episode UUID from path
		episodeUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest api.CreateEpisodeRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update episode
		if err := handler.api.Update(c, episodeUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating episode: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Find returns a Gin handler function for finding a episode by its UUID.
//
// @Summary Find a episode
// @Description Find a episode
// @Tags episodes
// @ID find-episode
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/episodes/{uuid} [get]
func (handler episodeHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract episode UUID from path
		episodeUUID := c.Param("uuid")

		// Call API to find episode
		episode, err := handler.api.Find(c, episodeUUID)
		if err != nil {
			log.Error().Msg("error finding episode: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, episode)
	}
}

// FindAll returns a Gin handler function for finding all episodes.
//
// @Summary Find all episodes
// @Description Find all episodes
// @Tags episodes
// @ID find-all-episodes
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/episodes/ [get]
func (handler episodeHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all episodes
		episodes, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all episodes: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, episodes)
	}
}

// Delete returns a Gin handler function for deleting a episode by its UUID.
//
// @Summary Delete a episode
// @Description Delete a episode
// @Tags episodes
// @ID delete-episode
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/episodes/{uuid} [delete]
func (handler episodeHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract episode UUID from path
		episodeUUID := c.Param("uuid")

		// Call API to delete episode
		if err := handler.api.Delete(c, episodeUUID); err != nil {
			log.Error().Msg("error deleting episode: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}
