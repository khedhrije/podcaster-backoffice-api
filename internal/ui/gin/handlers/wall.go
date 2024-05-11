// Package handlers provides HTTP request handlers for managing walls.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Wall represents the interface for managing walls.
type Wall interface {
	// Create returns a Gin handler function for creating a wall.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a wall.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a wall by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all walls.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a wall by its UUID.
	Delete() gin.HandlerFunc
}

type wallHandler struct {
	api api.Wall
}

// NewWallHandler creates a new instance of Wall interface.
func NewWallHandler(api api.Wall) Wall {
	return &wallHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a wall.
//
// @Summary Create a new wall
// @Description Create a new wall
// @Tags walls
// @ID create-wall
// @Param request body pkg.CreateWallRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls [post]
func (handler wallHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest api.CreateWallRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create wall
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Update returns a Gin handler function for updating a wall.
//
// @Summary Update wall
// @Description Update wall
// @Tags walls
// @ID update-wall
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateWallRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/{uuid} [put]
func (handler wallHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract wall UUID from path
		wallUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest api.CreateWallRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update wall
		if err := handler.api.Update(c, wallUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// Find returns a Gin handler function for finding a wall by its UUID.
//
// @Summary Find a wall
// @Description Find a wall
// @Tags walls
// @ID find-wall
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/{uuid} [get]
func (handler wallHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract wall UUID from path
		wallUUID := c.Param("uuid")

		// Call API to find wall
		wall, err := handler.api.Find(c, wallUUID)
		if err != nil {
			log.Error().Msg("error finding wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, wall)
	}
}

// FindAll returns a Gin handler function for finding all walls.
//
// @Summary Find all walls
// @Description Find all walls
// @Tags walls
// @ID find-all-walls
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/ [get]
func (handler wallHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all walls
		walls, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all walls: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, walls)
	}
}

// Delete returns a Gin handler function for deleting a wall by its UUID.
//
// @Summary Delete a wall
// @Description Delete a wall
// @Tags walls
// @ID delete-wall
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/{uuid} [delete]
func (handler wallHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract wall UUID from path
		wallUUID := c.Param("uuid")

		// Call API to delete wall
		if err := handler.api.Delete(c, wallUUID); err != nil {
			log.Error().Msg("error deleting wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}
