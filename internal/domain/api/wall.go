// Package wall provides functionality for managing walls.
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

// Wall represents the interface for managing walls.
type Wall interface {
	Create(ctx context.Context, wall CreateWallRequest) error
	Update(ctx context.Context, uuid string, updates UpdateWallRequest) error
	Find(ctx context.Context, uuid string) (*pkg.WallResponse, error)
	FindAll(ctx context.Context) ([]*pkg.WallResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// UpdateWallRequest represents the interface for updating walls.
type UpdateWallRequest interface {
	Name() string
	Description() string
}

// wallApi is an implementation of the Wall interface.
type wallApi struct {
	wallAdapter port.WallPersister
}

// NewWallApi creates a new instance of Wall.
func NewWallApi(wallAdapter port.WallPersister) Wall {
	return &wallApi{
		wallAdapter: wallAdapter,
	}
}

// Create creates a new wall.
func (api wallApi) Create(ctx context.Context, req CreateWallRequest) error {
	// Validate request
	vErrs := createWallRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	wall := model.Wall{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
	}
	// call adapter
	if err := api.wallAdapter.Create(ctx, wall); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("wall", wall).Msg("error while creating wall")
		return fmt.Errorf("error occurred while creating wall: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
func createWallRequestValidation(ctx context.Context, req CreateWallRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing wall.
func (api wallApi) Update(ctx context.Context, uuid string, updates UpdateWallRequest) error {
	// Validate request
	vErrs := updateWallRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	wall := model.Wall{
		Name:        updates.Name(),
		Description: updates.Description(),
	}
	// call adapter
	if err := api.wallAdapter.Update(ctx, uuid, wall); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("wall", wall).Msg("error while updating wall")
		return fmt.Errorf("error occurred while updating wall: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
func updateWallRequestValidation(ctx context.Context, uuid string, req UpdateWallRequest) model.ValidationErrors {
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

// Find finds a wall by UUID.
func (api wallApi) Find(ctx context.Context, uuid string) (*pkg.WallResponse, error) {

	// Call adapter
	wall, err := api.wallAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding wall")
		return nil, fmt.Errorf("error occurred while finding wall: %w", err)
	}

	// Map to response
	response := &pkg.WallResponse{
		ID:          wall.ID,
		Name:        wall.Name,
		Description: wall.Description,
	}
	// return result
	return response, nil
}

// FindAll finds all walls.
func (api wallApi) FindAll(ctx context.Context) ([]*pkg.WallResponse, error) {
	// Call adapter
	wallSlice, err := api.wallAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding wall")
		return nil, fmt.Errorf("error occurred while finding wall: %w", err)
	}

	// Map to response
	var response []*pkg.WallResponse
	for _, wall := range wallSlice {
		response = append(response, &pkg.WallResponse{
			ID:          wall.ID,
			Name:        wall.Name,
			Description: wall.Description,
		})
	}
	// return result
	return response, nil
}

// Delete deletes a wall by UUID.
func (api wallApi) Delete(ctx context.Context, uuid string) error {
	if err := api.wallAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting wall")
		return fmt.Errorf("error occurred while deleting wall: %w", err)
	}
	return nil
}

// CreateWallRequest represents the interface for creating walls.
type CreateWallRequest interface {
	Name() string
	Description() string
}
