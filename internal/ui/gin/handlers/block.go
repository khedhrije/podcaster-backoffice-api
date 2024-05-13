// Package handlers provides HTTP request handlers for managing blocks.
package handlers

import (
	"github.com/khedhrije/podcaster-backoffice-api/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Block represents the interface for managing blocks.
type Block interface {
	// Create returns a Gin handler function for creating a block.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a block.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a block by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all blocks.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a block by its UUID.
	Delete() gin.HandlerFunc

	// FindPrograms returns a Gin handler function for finding all programs of a block.
	FindPrograms() gin.HandlerFunc

	// OverwritePrograms returns a Gin handler function for overwriting the programs of a block.
	OverwritePrograms() gin.HandlerFunc
}

// blockHandler is an implementation of the Block interface.
type blockHandler struct {
	api api.Block
}

// NewBlockHandler creates a new instance of Block interface.
func NewBlockHandler(api api.Block) Block {
	return &blockHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a block.
//
// @Summary Create a new block
// @Description Create a new block
// @Tags blocks
// @ID create-block
// @Param request body pkg.CreateBlockRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/blocks [post]
func (handler blockHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest pkg.CreateBlockRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create block
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating block: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

// Update returns a Gin handler function for updating a block.
//
// @Summary Update block
// @Description Update block
// @Tags blocks
// @ID update-block
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateBlockRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/blocks/{uuid} [put]
func (handler blockHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract block UUID from path
		blockUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest pkg.UpdateBlockRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update block
		if err := handler.api.Update(c, blockUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating block: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

// Find returns a Gin handler function for finding a block by its UUID.
//
// @Summary Find a block
// @Description Find a block
// @Tags blocks
// @ID find-block
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {object} pkg.BlockResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/blocks/{uuid} [get]
func (handler blockHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract block UUID from path
		blockUUID := c.Param("uuid")

		// Call API to find block
		block, err := handler.api.Find(c, blockUUID)
		if err != nil {
			log.Error().Msg("error finding block: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, block)
	}
}

// FindAll returns a Gin handler function for finding all blocks.
//
// @Summary Find all blocks
// @Description Find all blocks
// @Tags blocks
// @ID find-all-blocks
// @Produce json
// @Success 200 {array} pkg.BlockResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/blocks [get]
func (handler blockHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all blocks
		blocks, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all blocks: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, blocks)
	}
}

// Delete returns a Gin handler function for deleting a block by its UUID.
//
// @Summary Delete a block
// @Description Delete a block
// @Tags blocks
// @ID delete-block
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "deleted"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/blocks/{uuid} [delete]
func (handler blockHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract block UUID from path
		blockUUID := c.Param("uuid")

		// Call API to delete block
		if err := handler.api.Delete(c, blockUUID); err != nil {
			log.Error().Msg("error deleting block: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}

// FindPrograms returns a Gin handler function for finding all programs of a block.
//
// @Summary Find all programs of a block
// @Description Find all programs of a block
// @Tags blocks
// @ID find-block-programs
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {array} pkg.BlockProgramsResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/blocks/{uuid}/programs [get]
func (handler blockHandler) FindPrograms() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract block UUID from path
		blockUUID := c.Param("uuid")

		// Call API to find all programs of the block
		programs, err := handler.api.FindPrograms(c, blockUUID)
		if err != nil {
			log.Error().Msg("error finding all programs of the block: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, programs)
	}
}

// OverwritePrograms returns a Gin handler function for overwriting programs.
//
// @Summary Overwrite programs of a block
// @Description Overwrite the programs of a specific block by replacing all existing programs with new ones
// @Tags blocks
// @ID overwrite-block-programs
// @Param uuid path string true "UUID of the block"
// @Param request body pkg.OverwriteProgramsRequestJSON true "List of programs' UUIDs to set"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 422 {object} pkg.ErrorJSON "Unprocessable Entity"
// @Failure 500 {object} pkg.ErrorJSON "Internal Server Error"
// @Router /private/blocks/{uuid}/programs/overwrite [put]
func (handler blockHandler) OverwritePrograms() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract block UUID from path
		blockUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest pkg.OverwriteProgramsRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to overwrite programs
		if err := handler.api.OverwritePrograms(c, blockUUID, jsonRequest); err != nil {
			log.Error().Msg("error overwriting programs: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "ok")
	}
}
