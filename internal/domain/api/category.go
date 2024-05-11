// Package category provides functionality for managing categories.
package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
	"github.com/khedhrije/podcaster-backoffice-api/pkg"
	"github.com/rs/zerolog/log"
)

// Category represents the interface for managing categories.
type Category interface {
	Create(ctx context.Context, category CreateCategoryRequest) error
	Update(ctx context.Context, uuid string, updates UpdateCategoryRequest) error
	Find(ctx context.Context, uuid string) (*pkg.CategoryResponse, error)
	FindAll(ctx context.Context) ([]*pkg.CategoryResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// UpdateCategoryRequest represents the interface for updating categories.
type UpdateCategoryRequest interface {
	Name() string
	Description() string
}

// categoryApi is an implementation of the Category interface.
type categoryApi struct {
	categoryAdapter port.CategoryPersister
}

// NewCategoryApi creates a new instance of Category.
func NewCategoryApi(categoryAdapter port.CategoryPersister) Category {
	return &categoryApi{
		categoryAdapter: categoryAdapter,
	}
}

// Create creates a new category.
func (api categoryApi) Create(ctx context.Context, req CreateCategoryRequest) error {
	// Validate request
	vErrs := createCategoryRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	category := model.Category{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
	}
	// call adapter
	if err := api.categoryAdapter.Create(ctx, category); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("category", category).Msg("error while creating category")
		return fmt.Errorf("error occurred while creating category: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
func createCategoryRequestValidation(ctx context.Context, req CreateCategoryRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing category.
func (api categoryApi) Update(ctx context.Context, uuid string, updates UpdateCategoryRequest) error {
	// Validate request
	vErrs := updateCategoryRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	category := model.Category{
		Name:        updates.Name(),
		Description: updates.Description(),
	}
	// call adapter
	if err := api.categoryAdapter.Update(ctx, uuid, category); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("category", category).Msg("error while updating category")
		return fmt.Errorf("error occurred while updating category: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
func updateCategoryRequestValidation(ctx context.Context, uuid string, req UpdateCategoryRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if uuid == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "uuid", Message: "cannot be empty"})
	}
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "cannot be empty"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "cannot be empty"})
	}
	return vErrs
}

// Find finds a category by UUID.
func (api categoryApi) Find(ctx context.Context, uuid string) (*pkg.CategoryResponse, error) {

	// Call adapter
	category, err := api.categoryAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding category")
		return nil, fmt.Errorf("error occurred while finding category: %w", err)
	}

	// Map to response
	response := &pkg.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	// return result
	return response, nil
}

// FindAll finds all categories.
func (api categoryApi) FindAll(ctx context.Context) ([]*pkg.CategoryResponse, error) {
	// Call adapter
	categorieslice, err := api.categoryAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding category")
		return nil, fmt.Errorf("error occurred while finding category: %w", err)
	}

	// Map to response
	var response []*pkg.CategoryResponse
	for _, category := range categorieslice {
		response = append(response, &pkg.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	// return result
	return response, nil
}

// Delete deletes a category by UUID.
func (api categoryApi) Delete(ctx context.Context, uuid string) error {
	if err := api.categoryAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting category")
		return fmt.Errorf("error occurred while deleting category: %w", err)
	}
	return nil
}

// CreateCategoryRequest represents the interface for creating categories.
type CreateCategoryRequest interface {
	Name() string
	Description() string
}
