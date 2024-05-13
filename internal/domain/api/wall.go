// Package api provides functionality for managing walls.
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
	FindBlocks(ctx context.Context, uuid string) ([]*pkg.WallBlocksResponse, error)
	OverwriteBlocks(ctx context.Context, wallID string, req OverwriteBlocksRequest) error
}

// wallApi is an implementation of the Wall interface.
type wallApi struct {
	wallAdapter      port.WallPersister
	wallBlockAdapter port.WallBlockPersister
	blockAdapter     port.BlockPersister
}

// NewWallApi creates a new instance of Wall.
// It takes wallAdapter, wallBlockAdapter, and blockAdapter as dependencies.
func NewWallApi(wallAdapter port.WallPersister, wallBlockAdapter port.WallBlockPersister, blockAdapter port.BlockPersister) Wall {
	return &wallApi{
		wallAdapter:      wallAdapter,
		wallBlockAdapter: wallBlockAdapter,
		blockAdapter:     blockAdapter,
	}
}

// FindBlocks finds blocks associated with a wall.
// It takes the context and the wall UUID, and returns a slice of WallBlocksResponse or an error.
func (api wallApi) FindBlocks(ctx context.Context, uuid string) ([]*pkg.WallBlocksResponse, error) {
	// Find associations by wall ID
	associations, err := api.wallBlockAdapter.FindByWallID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var response []*pkg.WallBlocksResponse
	for _, association := range associations {
		// Find block by ID
		block, err := api.blockAdapter.Find(ctx, association.BlockID)
		if err != nil {
			return nil, err
		}

		// Map block information to response
		wallBlocks := &pkg.WallBlocksResponse{
			BlockResponse: pkg.BlockResponse{
				ID:          block.ID,
				Name:        block.Name,
				Kind:        block.Kind,
				Description: block.Description,
			},
			Position: association.Position,
		}

		response = append(response, wallBlocks)
	}

	return response, nil
}

// Create creates a new wall.
// It takes the context and CreateWallRequest, and returns an error if any.
func (api wallApi) Create(ctx context.Context, req CreateWallRequest) error {
	// Validate request
	vErrs := createWallRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}

	// Map request to domain model
	wall := model.Wall{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
	}

	// Call adapter to create wall
	if err := api.wallAdapter.Create(ctx, wall); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("wall", wall).Msg("error while creating wall")
		return fmt.Errorf("error occurred while creating wall: %w", err)
	}

	return nil
}

// createWallRequestValidation validates the creation request.
// It takes the context and CreateWallRequest, and returns a slice of ValidationErrors.
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
// It takes the context, wall UUID, and UpdateWallRequest, and returns an error if any.
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

	// Call adapter to update wall
	if err := api.wallAdapter.Update(ctx, uuid, wall); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("wall", wall).Msg("error while updating wall")
		return fmt.Errorf("error occurred while updating wall: %w", err)
	}

	return nil
}

// updateWallRequestValidation validates the update request.
// It takes the context, wall UUID, and UpdateWallRequest, and returns a slice of ValidationErrors.
func updateWallRequestValidation(ctx context.Context, uuid string, req UpdateWallRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if uuid == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "uuid", Message: "cannot be empty"})
	}
	return vErrs
}

// Find finds a wall by UUID.
// It takes the context and wall UUID, and returns a WallResponse or an error.
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

	// Return result
	return response, nil
}

// FindAll finds all walls.
// It takes the context and returns a slice of WallResponse or an error.
func (api wallApi) FindAll(ctx context.Context) ([]*pkg.WallResponse, error) {
	// Call adapter
	wallSlice, err := api.wallAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding walls")
		return nil, fmt.Errorf("error occurred while finding walls: %w", err)
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

	// Return result
	return response, nil
}

// Delete deletes a wall by UUID.
// It takes the context and wall UUID, and returns an error if any.
func (api wallApi) Delete(ctx context.Context, uuid string) error {
	if err := api.wallAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting wall")
		return fmt.Errorf("error occurred while deleting wall: %w", err)
	}
	return nil
}

// OverwriteBlocks overwrites the blocks associated with a wall.
// It takes the context, wall ID, and OverwriteBlocksRequest, and returns an error if any.
func (api wallApi) OverwriteBlocks(ctx context.Context, wallID string, req OverwriteBlocksRequest) error {
	// Find all existing associations by wall ID
	associations, err := api.wallBlockAdapter.FindByWallID(ctx, wallID)
	if err != nil {
		return err
	}

	// Remove all existing associations for the wall
	for _, association := range associations {
		if err = api.wallBlockAdapter.Delete(ctx, association.ID); err != nil {
			return err
		}
	}

	// Create all new associations
	for blockID, position := range req.OrderedBlocks() {
		wallBlock := model.WallBlock{
			ID:       uuid.New().String(),
			WallID:   wallID,
			BlockID:  blockID,
			Position: position,
		}

		// Call adapter to create wall block
		if err := api.wallBlockAdapter.Create(ctx, wallBlock); err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("wallBlock", wallBlock).Msg("error while creating wall block")
			return fmt.Errorf("error occurred while creating wall block: %w", err)
		}
	}

	return nil
}

// OverwriteBlocksRequest represents the interface for overwriting wallBlock associations.
type OverwriteBlocksRequest interface {
	OrderedBlocks() map[string]int
}

// CreateWallRequest represents the interface for creating walls.
type CreateWallRequest interface {
	Name() string
	Description() string
}

// UpdateWallRequest represents the interface for updating walls.
type UpdateWallRequest interface {
	Name() string
	Description() string
}
