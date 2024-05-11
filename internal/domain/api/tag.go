// Package tag provides functionality for managing tags.
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

// Tag represents the interface for managing tags.
type Tag interface {
	Create(ctx context.Context, tag CreateTagRequest) error
	Update(ctx context.Context, uuid string, updates UpdateTagRequest) error
	Find(ctx context.Context, uuid string) (*pkg.TagResponse, error)
	FindAll(ctx context.Context) ([]*pkg.TagResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// UpdateTagRequest represents the interface for updating tags.
type UpdateTagRequest interface {
	Name() string
	Description() string
}

// tagApi is an implementation of the Tag interface.
type tagApi struct {
	tagAdapter port.TagPersister
}

// NewTagApi creates a new instance of Tag.
func NewTagApi(tagAdapter port.TagPersister) Tag {
	return &tagApi{
		tagAdapter: tagAdapter,
	}
}

// Create creates a new tag.
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
	// call adapter
	if err := api.tagAdapter.Create(ctx, tag); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("tag", tag).Msg("error while creating tag")
		return fmt.Errorf("error occurred while creating tag: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
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
	// call adapter
	if err := api.tagAdapter.Update(ctx, uuid, tag); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("tag", tag).Msg("error while updating tag")
		return fmt.Errorf("error occurred while updating tag: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
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
	// return result
	return response, nil
}

// FindAll finds all tags.
func (api tagApi) FindAll(ctx context.Context) ([]*pkg.TagResponse, error) {
	// Call adapter
	tagSlice, err := api.tagAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding tag")
		return nil, fmt.Errorf("error occurred while finding tag: %w", err)
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
	// return result
	return response, nil
}

// Delete deletes a tag by UUID.
func (api tagApi) Delete(ctx context.Context, uuid string) error {
	if err := api.tagAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting tag")
		return fmt.Errorf("error occurred while deleting tag: %w", err)
	}
	return nil
}

// CreateTagRequest represents the interface for creating tags.
type CreateTagRequest interface {
	Name() string
	Description() string
}
