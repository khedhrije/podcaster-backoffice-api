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

// CreateCategoryRequest represents the interface for creating categories.
type CreateCategoryRequest interface {
	Name() string
	Description() string
	ParentID() string
}

// UpdateCategoryRequest represents the interface for updating categories.
type UpdateCategoryRequest interface {
	Name() string
	Description() string
	ParentID() string
}

// Category represents the interface for managing categories.
type Category interface {
	Create(ctx context.Context, category CreateCategoryRequest) error
	Update(ctx context.Context, uuid string, updates UpdateCategoryRequest) error
	Find(ctx context.Context, uuid string) (*pkg.CategoryResponse, error)
	FindAll(ctx context.Context) ([]*pkg.CategoryResponse, error)
	Delete(ctx context.Context, uuid string) error
	FindPrograms(ctx context.Context, uuid string) ([]*pkg.ProgramResponse, error)
}

// categoryApi is an implementation of the Category interface.
type categoryApi struct {
	categoryAdapter   port.CategoryPersister
	programCatAdapter port.ProgramCategoryPersister
	programAdapter    port.ProgramPersister
}

// NewCategoryApi creates a new instance of Category.
func NewCategoryApi(categoryAdapter port.CategoryPersister, programCatAdapter port.ProgramCategoryPersister, programAdapter port.ProgramPersister) Category {
	return &categoryApi{
		categoryAdapter:   categoryAdapter,
		programCatAdapter: programCatAdapter,
		programAdapter:    programAdapter,
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
		Parent: &model.Category{
			ID: req.ParentID(),
		},
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
		log.
			Ctx(ctx).
			Error().
			Err(vErrs).
			Interface("updates", updates).
			Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}

	// Map to domain model
	category := model.Category{}
	if updates.Name() != "" {
		category.Name = updates.Name()
	}
	if updates.Description() != "" {
		category.Description = updates.Description()
	}
	if updates.ParentID() != "" {
		category.Parent = &model.Category{
			ID: updates.ParentID(),
		}
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

	parentID := ""
	if category.Parent != nil {
		parentID = category.Parent.ID
	}
	// Map to response
	response := &pkg.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    parentID,
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
		parentID := ""
		if category.Parent != nil {
			parentID = category.Parent.ID
		}
		response = append(response, &pkg.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			ParentID:    parentID,
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

// FindPrograms finds programs associated with a cat.
func (api categoryApi) FindPrograms(ctx context.Context, uuid string) ([]*pkg.ProgramResponse, error) {

	associations, err := api.programCatAdapter.FindByCategoryID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var response []*pkg.ProgramResponse
	for _, association := range associations {
		program, err := api.programAdapter.Find(ctx, association.ProgramID)
		if err != nil {
			return nil, err
		}

		response = append(response, &pkg.ProgramResponse{
			ID:          program.ID,
			Name:        program.Name,
			Description: program.Description,
		})

	}

	return response, nil
}
