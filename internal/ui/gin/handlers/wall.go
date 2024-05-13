// Package handlers provides HTTP request handlers for managing walls.
package handlers

import (
	"github.com/khedhrije/podcaster-backoffice-api/pkg"
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

	// FindBlocks returns a Gin handler function for finding all blocks associated with a wall.
	FindBlocks() gin.HandlerFunc

	// OverwriteBlocks returns a Gin handler function for overwriting the blocks of a wall.
	OverwriteBlocks() gin.HandlerFunc
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
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonRequest pkg.CreateWallRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "ok")
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
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallUUID := c.Param("uuid")

		var jsonRequest pkg.UpdateWallRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if err := handler.api.Update(c, wallUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "ok")
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
// @Success 200 {object} pkg.WallResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/{uuid} [get]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallUUID := c.Param("uuid")

		wall, err := handler.api.Find(c, wallUUID)
		if err != nil {
			log.Error().Msg("error finding wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
// @Success 200 {array} pkg.WallResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls [get]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		walls, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all walls: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
// @Success 200 {string} string "deleted"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/{uuid} [delete]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallUUID := c.Param("uuid")

		if err := handler.api.Delete(c, wallUUID); err != nil {
			log.Error().Msg("error deleting wall: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "deleted")
	}
}

// FindBlocks returns a Gin handler function for finding all blocks associated with a wall.
//
// @Summary Find all wall's blocks
// @Description Find all wall's blocks
// @Tags walls
// @ID find-wall-block
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {array} pkg.BlockResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/walls/{uuid}/blocks [get]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) FindBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallUUID := c.Param("uuid")

		blocks, err := handler.api.FindBlocks(c, wallUUID)
		if err != nil {
			log.Error().Msg("error finding all wall's blocks: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, blocks)
	}
}

// OverwriteBlocks returns a Gin handler function for overwriting the blocks of a wall.
//
// @Summary Overwrite blocks of a wall
// @Description Overwrite the blocks of a specific wall by replacing all existing blocks with new ones
// @Tags walls
// @ID overwrite-wall-blocks
// @Param uuid path string true "UUID of the wall"
// @Param request body pkg.OverwriteBlocksRequestJSON true "List of blocks' UUIDs to set"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 422 {object} pkg.ErrorJSON "Unprocessable Entity"
// @Failure 500 {object} pkg.ErrorJSON "Internal Server Error"
// @Router /private/walls/{uuid}/blocks/overwrite [put]
//
// @Security Bearer-APIKey || Bearer-JWT
func (handler wallHandler) OverwriteBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallUUID := c.Param("uuid")

		var jsonRequest pkg.OverwriteBlocksRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if err := handler.api.OverwriteBlocks(c, wallUUID, jsonRequest); err != nil {
			log.Error().Msg("error overwriting blocks: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "ok")
	}
}
