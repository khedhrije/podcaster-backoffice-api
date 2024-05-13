// Package handlers provides HTTP request handlers for managing categories.
package handlers

import (
	"github.com/khedhrije/podcaster-backoffice-api/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/api"
	"github.com/rs/zerolog/log"
)

// Category represents the interface for managing categories.
type Category interface {
	// Create returns a Gin handler function for creating a category.
	Create() gin.HandlerFunc

	// Update returns a Gin handler function for updating a category.
	Update() gin.HandlerFunc

	// Find returns a Gin handler function for finding a category by its UUID.
	Find() gin.HandlerFunc

	// FindAll returns a Gin handler function for finding all categories.
	FindAll() gin.HandlerFunc

	// Delete returns a Gin handler function for deleting a category by its UUID.
	Delete() gin.HandlerFunc

	// FindPrograms returns a Gin handler function for finding all programs associated with a category.
	FindPrograms() gin.HandlerFunc
}

// categoryHandler is an implementation of the Category interface.
type categoryHandler struct {
	api api.Category
}

// NewCategoryHandler creates a new instance of Category interface.
func NewCategoryHandler(api api.Category) Category {
	return &categoryHandler{
		api: api,
	}
}

// Create returns a Gin handler function for creating a category.
//
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @ID create-category
// @Param request body pkg.CreateCategoryRequestJSON true "create request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/categories [post]
func (handler categoryHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract body request
		var jsonRequest pkg.CreateCategoryRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to create category
		if err := handler.api.Create(c, jsonRequest); err != nil {
			log.Error().Msg("error creating category: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

// Update returns a Gin handler function for updating a category.
//
// @Summary Update category
// @Description Update category
// @Tags categories
// @ID update-category
// @Param uuid path string true "uuid"
// @Param request body pkg.UpdateCategoryRequestJSON true "update request"
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/categories/{uuid} [put]
func (handler categoryHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract category UUID from path
		categoryUUID := c.Param("uuid")

		// Extract body request
		var jsonRequest pkg.UpdateCategoryRequestJSON
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			log.Ctx(c).Error().Err(err).Msg("error binding request")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// Call API to update category
		if err := handler.api.Update(c, categoryUUID, jsonRequest); err != nil {
			log.Error().Msg("error updating category: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

// Find returns a Gin handler function for finding a category by its UUID.
//
// @Summary Find a category
// @Description Find a category
// @Tags categories
// @ID find-category
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {object} pkg.CategoryResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/categories/{uuid} [get]
func (handler categoryHandler) Find() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract category UUID from path
		categoryUUID := c.Param("uuid")

		// Call API to find category
		category, err := handler.api.Find(c, categoryUUID)
		if err != nil {
			log.Error().Msg("error finding category: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, category)
	}
}

// FindAll returns a Gin handler function for finding all categories.
//
// @Summary Find all categories
// @Description Find all categories
// @Tags categories
// @ID find-all-categories
// @Produce json
// @Success 200 {array} pkg.CategoryResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/categories [get]
func (handler categoryHandler) FindAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call API to find all categories
		categories, err := handler.api.FindAll(c)
		if err != nil {
			log.Error().Msg("error finding all categories: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, categories)
	}
}

// Delete returns a Gin handler function for deleting a category by its UUID.
//
// @Summary Delete a category
// @Description Delete a category
// @Tags categories
// @ID delete-category
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {string} string "deleted"
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/categories/{uuid} [delete]
func (handler categoryHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract category UUID from path
		categoryUUID := c.Param("uuid")

		// Call API to delete category
		if err := handler.api.Delete(c, categoryUUID); err != nil {
			log.Error().Msg("error deleting category: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, "deleted")
	}
}

// FindPrograms returns a Gin handler function for finding all programs associated with a category.
//
// @Summary Find all category's programs
// @Description Find all category's programs
// @Tags categories
// @ID find-category-programs
// @Param uuid path string true "uuid"
// @Produce json
// @Success 200 {array} pkg.ProgramResponse
// @Failure 500 {object} pkg.ErrorJSON
// @Router /private/categories/{uuid}/programs [get]
func (handler categoryHandler) FindPrograms() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract category UUID from path
		categoryUUID := c.Param("uuid")

		// Call API to find all programs associated with the category
		programs, err := handler.api.FindPrograms(c, categoryUUID)
		if err != nil {
			log.Error().Msg("error finding all category's programs: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return response
		c.JSON(http.StatusOK, programs)
	}
}
