// Package program provides functionality for managing programs.
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

// Program represents the interface for managing programs.
type Program interface {
	Create(ctx context.Context, program CreateProgramRequest) error
	Update(ctx context.Context, uuid string, updates UpdateProgramRequest) error
	Find(ctx context.Context, uuid string) (*pkg.ProgramResponse, error)
	FindAll(ctx context.Context) ([]*pkg.ProgramResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// UpdateProgramRequest represents the interface for updating programs.
type UpdateProgramRequest interface {
	Name() string
	Description() string
}

// programApi is an implementation of the Program interface.
type programApi struct {
	programAdapter port.ProgramPersister
}

// NewProgramApi creates a new instance of Program.
func NewProgramApi(programAdapter port.ProgramPersister) Program {
	return programApi{
		programAdapter: programAdapter,
	}
}

// Create creates a new program.
func (api programApi) Create(ctx context.Context, req CreateProgramRequest) error {
	// Validate request
	vErrs := createProgramRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	program := model.Program{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
	}
	// call adapter
	if err := api.programAdapter.Create(ctx, program); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("program", program).Msg("error while creating program")
		return fmt.Errorf("error occurred while creating program: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
func createProgramRequestValidation(ctx context.Context, req CreateProgramRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing program.
func (api programApi) Update(ctx context.Context, uuid string, updates UpdateProgramRequest) error {
	// Validate request
	vErrs := updateProgramRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	program := model.Program{
		Name:        updates.Name(),
		Description: updates.Description(),
	}
	// call adapter
	if err := api.programAdapter.Update(ctx, uuid, program); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("program", program).Msg("error while updating program")
		return fmt.Errorf("error occurred while updating program: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
func updateProgramRequestValidation(ctx context.Context, uuid string, req UpdateProgramRequest) model.ValidationErrors {
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

// Find finds a program by UUID.
func (api programApi) Find(ctx context.Context, uuid string) (*pkg.ProgramResponse, error) {

	// Call adapter
	program, err := api.programAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding program")
		return nil, fmt.Errorf("error occurred while finding program: %w", err)
	}

	// Map to response
	response := &pkg.ProgramResponse{
		ID:          program.ID,
		Name:        program.Name,
		Description: program.Description,
	}
	// return result
	return response, nil
}

// FindAll finds all programs.
func (api programApi) FindAll(ctx context.Context) ([]*pkg.ProgramResponse, error) {
	// Call adapter
	programSlice, err := api.programAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding program")
		return nil, fmt.Errorf("error occurred while finding program: %w", err)
	}

	// Map to response
	var response []*pkg.ProgramResponse
	for _, program := range programSlice {
		response = append(response, &pkg.ProgramResponse{
			ID:          program.ID,
			Name:        program.Name,
			Description: program.Description,
		})
	}
	// return result
	return response, nil
}

// Delete deletes a program by UUID.
func (api programApi) Delete(ctx context.Context, uuid string) error {
	if err := api.programAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting program")
		return fmt.Errorf("error occurred while deleting program: %w", err)
	}
	return nil
}

// CreateProgramRequest represents the interface for creating programs.
type CreateProgramRequest interface {
	Name() string
	Description() string
}
