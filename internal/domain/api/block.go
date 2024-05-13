// Package block provides functionality for managing blocks.
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

// CreateBlockRequest represents the interface for creating blocks.
type CreateBlockRequest interface {
	Name() string
	Description() string
	Kind() string
}

// UpdateBlockRequest represents the interface for updating blocks.
type UpdateBlockRequest interface {
	Name() string
	Description() string
	Kind() string
}

// Block represents the interface for managing blocks.
type Block interface {
	Create(ctx context.Context, block CreateBlockRequest) error
	Update(ctx context.Context, uuid string, updates UpdateBlockRequest) error
	Find(ctx context.Context, uuid string) (*pkg.BlockResponse, error)
	FindAll(ctx context.Context) ([]*pkg.BlockResponse, error)
	Delete(ctx context.Context, uuid string) error
	FindPrograms(ctx context.Context, uuid string) ([]*pkg.BlockProgramsResponse, error)
	OverwritePrograms(ctx context.Context, blockID string, req OverwriteProgramsRequest) error
}

// blockApi is an implementation of the Block interface.
type blockApi struct {
	blockAdapter        port.BlockPersister
	blockProgramAdapter port.BlockProgramPersister
	programAdapter      port.ProgramPersister
}

// NewBlockApi creates a new instance of Block.
func NewBlockApi(blockAdapter port.BlockPersister, blockProgramAdapter port.BlockProgramPersister, programAdapter port.ProgramPersister) Block {
	return &blockApi{
		blockAdapter:        blockAdapter,
		blockProgramAdapter: blockProgramAdapter,
		programAdapter:      programAdapter,
	}
}

// Create creates a new block.
func (api blockApi) Create(ctx context.Context, req CreateBlockRequest) error {
	// Validate request
	vErrs := createBlockRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	block := model.Block{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
		Kind:        req.Kind(),
	}
	// call adapter
	if err := api.blockAdapter.Create(ctx, block); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("block", block).Msg("error while creating block")
		return fmt.Errorf("error occurred while creating block: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
func createBlockRequestValidation(ctx context.Context, req CreateBlockRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing block.
func (api blockApi) Update(ctx context.Context, uuid string, updates UpdateBlockRequest) error {
	// Validate request
	vErrs := updateBlockRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	block := model.Block{}
	if updates.Name() != "" {
		block.Name = updates.Name()
	}
	if updates.Description() != "" {
		block.Description = updates.Description()
	}
	if updates.Kind() != "" {
		block.Kind = updates.Kind()
	}

	// call adapter
	if err := api.blockAdapter.Update(ctx, uuid, block); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("block", block).Msg("error while updating block")
		return fmt.Errorf("error occurred while updating block: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
func updateBlockRequestValidation(ctx context.Context, uuid string, req UpdateBlockRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if uuid == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "uuid", Message: "cannot be empty"})
	}
	return vErrs
}

// Find finds a block by UUID.
func (api blockApi) Find(ctx context.Context, uuid string) (*pkg.BlockResponse, error) {

	// Call adapter
	block, err := api.blockAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding block")
		return nil, fmt.Errorf("error occurred while finding block: %w", err)
	}

	// Map to response
	response := &pkg.BlockResponse{
		ID:          block.ID,
		Name:        block.Name,
		Description: block.Description,
		Kind:        block.Kind,
	}
	// return result
	return response, nil
}

// FindAll finds all blocks.
func (api blockApi) FindAll(ctx context.Context) ([]*pkg.BlockResponse, error) {
	// Call adapter
	blockSlice, err := api.blockAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding block")
		return nil, fmt.Errorf("error occurred while finding block: %w", err)
	}

	// Map to response
	var response []*pkg.BlockResponse
	for _, block := range blockSlice {
		response = append(response, &pkg.BlockResponse{
			ID:          block.ID,
			Name:        block.Name,
			Description: block.Description,
			Kind:        block.Kind,
		})
	}
	// return result
	return response, nil
}

// Delete deletes a block by UUID.
func (api blockApi) Delete(ctx context.Context, uuid string) error {
	if err := api.blockAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting block")
		return fmt.Errorf("error occurred while deleting block: %w", err)
	}
	return nil
}

// FindPrograms finds programs associated with a block.
func (api blockApi) FindPrograms(ctx context.Context, uuid string) ([]*pkg.BlockProgramsResponse, error) {

	associations, err := api.blockProgramAdapter.FindByBlockID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var response []*pkg.BlockProgramsResponse
	for _, association := range associations {
		program, err := api.programAdapter.Find(ctx, association.ProgramID)
		if err != nil {
			return nil, err
		}
		blockProgram := &pkg.BlockProgramsResponse{
			ProgramResponse: pkg.ProgramResponse{
				ID:          program.ID,
				Name:        program.Name,
				Description: program.Description,
			},
			Position: association.Position,
		}

		response = append(response, blockProgram)
	}

	return response, nil
}

func (api blockApi) OverwritePrograms(ctx context.Context, blockID string, req OverwriteProgramsRequest) error {

	// Find all existing associations by wallID
	associations, err := api.blockProgramAdapter.FindByBlockID(ctx, blockID)
	if err != nil {
		return err
	}

	// Remove all existing association for the wall
	for _, association := range associations {
		err = api.blockProgramAdapter.Delete(ctx, association.ID)
		if err != nil {
			return err
		}
	}

	// Create all new associations
	for programID, position := range req.OrderedPrograms() {
		wallBlock := model.BlockProgram{
			ID:        uuid.New().String(),
			BlockID:   blockID,
			ProgramID: programID,
			Position:  position,
		}

		// Call adapter to create wall
		if err := api.blockProgramAdapter.Create(ctx, wallBlock); err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("wallBlock", wallBlock).Msg("error while creating wall")
			return fmt.Errorf("error occurred while creating wall: %w", err)
		}
	}

	return nil
}

type OverwriteProgramsRequest interface {
	OrderedPrograms() map[string]int
}
