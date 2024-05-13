// Package api provides functionality for managing tags.
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

// CreateTagRequest represents the interface for creating tags.
type CreateTagRequest interface {
	Name() string
	Description() string
}

// UpdateTagRequest represents the interface for updating tags.
type UpdateTagRequest interface {
	Name() string
	Description() string
}

// Tag represents the interface for managing tags.
type Tag interface {
	Create(ctx context.Context, tag CreateTagRequest) error
	Update(ctx context.Context, uuid string, updates UpdateTagRequest) error
	Find(ctx context.Context, uuid string) (*pkg.TagResponse, error)
	FindAll(ctx context.Context) ([]*pkg.TagResponse, error)
	Delete(ctx context.Context, uuid string) error
	FindPrograms(ctx context.Context, uuid string) ([]*pkg.ProgramResponse, error)
}

// tagApi is an implementation of the Tag interface.
type tagApi struct {
	tagAdapter        port.TagPersister
	programTagAdapter port.ProgramTagPersister
	programAdapter    port.ProgramPersister
}

// NewTagApi creates a new instance of Tag.
// It takes adapters for tag, program-tag, and program persistence as dependencies.
func NewTagApi(tagAdapter port.TagPersister, programTagAdapter port.ProgramTagPersister, programAdapter port.ProgramPersister) Tag {
	return &tagApi{
		tagAdapter:        tagAdapter,
		programTagAdapter: programTagAdapter,
		programAdapter:    programAdapter,
	}
}

// Create creates a new tag.
// It takes the context and CreateTagRequest, and returns an error if any.
func (api tagApi) Create(ctx context.Context, req CreateTagRequest) error {
	// Validate request
	vErrs := createTagRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	tag := model.Tag{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
	}
	// Call adapter
	if err := api.tagAdapter.Create(ctx, tag); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("tag", tag).Msg("error while creating tag")
		return fmt.Errorf("error occurred while creating tag: %w", err)
	}

	return nil
}

// createTagRequestValidation validates the creation request.
// It takes the context and CreateTagRequest, and returns a slice of ValidationErrors.
func createTagRequestValidation(ctx context.Context, req CreateTagRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing tag.
// It takes the context, tag UUID, and UpdateTagRequest, and returns an error if any.
func (api tagApi) Update(ctx context.Context, uuid string, updates UpdateTagRequest) error {
	// Validate request
	vErrs := updateTagRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	tag := model.Tag{
		Name:        updates.Name(),
		Description: updates.Description(),
	}
	// Call adapter
	if err := api.tagAdapter.Update(ctx, uuid, tag); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("tag", tag).Msg("error while updating tag")
		return fmt.Errorf("error occurred while updating tag: %w", err)
	}

	return nil
}

// updateTagRequestValidation validates the update request.
// It takes the context, tag UUID, and UpdateTagRequest, and returns a slice of ValidationErrors.
func updateTagRequestValidation(ctx context.Context, uuid string, req UpdateTagRequest) model.ValidationErrors {
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

// Find finds a tag by UUID.
// It takes the context and tag UUID, and returns a TagResponse or an error.
func (api tagApi) Find(ctx context.Context, uuid string) (*pkg.TagResponse, error) {
	// Call adapter
	tag, err := api.tagAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding tag")
		return nil, fmt.Errorf("error occurred while finding tag: %w", err)
	}

	// Map to response
	response := &pkg.TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Description: tag.Description,
	}
	// Return result
	return response, nil
}

// FindAll finds all tags.
// It takes the context and returns a slice of TagResponse or an error.
func (api tagApi) FindAll(ctx context.Context) ([]*pkg.TagResponse, error) {
	// Call adapter
	tagSlice, err := api.tagAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding tags")
		return nil, fmt.Errorf("error occurred while finding tags: %w", err)
	}

	// Map to response
	var response []*pkg.TagResponse
	for _, tag := range tagSlice {
		response = append(response, &pkg.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
		})
	}
	// Return result
	return response, nil
}

// Delete deletes a tag by UUID.
// It takes the context and tag UUID, and returns an error if any.
func (api tagApi) Delete(ctx context.Context, uuid string) error {
	if err := api.tagAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting tag")
		return fmt.Errorf("error occurred while deleting tag: %w", err)
	}
	return nil
}

// FindPrograms finds programs associated with a tag.
// It takes the context and tag UUID, and returns a slice of ProgramResponse or an error.
func (api tagApi) FindPrograms(ctx context.Context, uuid string) ([]*pkg.ProgramResponse, error) {
	associations, err := api.programTagAdapter.FindByTagID(ctx, uuid)
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
