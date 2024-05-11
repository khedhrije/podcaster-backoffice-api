// Package media provides functionality for managing medias.
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

// Media represents the interface for managing medias.
type Media interface {
	Create(ctx context.Context, media CreateMediaRequest) error
	Update(ctx context.Context, uuid string, updates UpdateMediaRequest) error
	Find(ctx context.Context, uuid string) (*pkg.MediaResponse, error)
	FindAll(ctx context.Context) ([]*pkg.MediaResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// UpdateMediaRequest represents the interface for updating medias.
type UpdateMediaRequest interface {
	DirectLink() string
	Kind() string
}

// mediaApi is an implementation of the Media interface.
type mediaApi struct {
	mediaAdapter port.MediaPersister
}

// New creates a new instance of Media.
func NewMediaApi(mediaAdapter port.MediaPersister) Media {
	return &mediaApi{
		mediaAdapter: mediaAdapter,
	}
}

// Create creates a new media.
func (api mediaApi) Create(ctx context.Context, req CreateMediaRequest) error {
	// Validate request
	vErrs := createMediaRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	media := model.Media{
		ID:         uuid.New().String(),
		DirectLink: req.DirectLink(),
		Kind:       req.Kind(),
	}
	// call adapter
	if err := api.mediaAdapter.Create(ctx, media); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("media", media).Msg("error while creating media")
		return fmt.Errorf("error occurred while creating media: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
func createMediaRequestValidation(ctx context.Context, req CreateMediaRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.DirectLink() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "directLink", Message: "is required"})
	}
	if req.Kind() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "kind", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing media.
func (api mediaApi) Update(ctx context.Context, uuid string, updates UpdateMediaRequest) error {
	// Validate request
	vErrs := updateMediaRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	media := model.Media{
		DirectLink: updates.DirectLink(),
		Kind:       updates.Kind(),
	}
	// call adapter
	if err := api.mediaAdapter.Update(ctx, uuid, media); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("media", media).Msg("error while updating media")
		return fmt.Errorf("error occurred while updating media: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
func updateMediaRequestValidation(ctx context.Context, uuid string, req UpdateMediaRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if uuid == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "uuid", Message: "cannot be empty"})
	}
	if req.DirectLink() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "directLink", Message: "cannot be empty"})
	}
	if req.Kind() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "kind", Message: "cannot be empty"})
	}
	return vErrs
}

// Find finds a media by UUID.
func (api mediaApi) Find(ctx context.Context, uuid string) (*pkg.MediaResponse, error) {

	// Call adapter
	media, err := api.mediaAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding media")
		return nil, fmt.Errorf("error occurred while finding media: %w", err)
	}

	// Map to response
	response := &pkg.MediaResponse{
		ID:         media.ID,
		DirectLink: media.DirectLink,
		Kind:       media.Kind,
	}
	// return result
	return response, nil
}

// FindAll finds all medias.
func (api mediaApi) FindAll(ctx context.Context) ([]*pkg.MediaResponse, error) {
	// Call adapter
	mediaSlice, err := api.mediaAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding media")
		return nil, fmt.Errorf("error occurred while finding media: %w", err)
	}

	// Map to response
	var response []*pkg.MediaResponse
	for _, media := range mediaSlice {
		response = append(response, &pkg.MediaResponse{
			ID:         media.ID,
			DirectLink: media.DirectLink,
			Kind:       media.Kind,
		})
	}
	// return result
	return response, nil
}

// Delete deletes a media by UUID.
func (api mediaApi) Delete(ctx context.Context, uuid string) error {
	if err := api.mediaAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting media")
		return fmt.Errorf("error occurred while deleting media: %w", err)
	}
	return nil
}

// CreateMediaRequest represents the interface for creating medias.
type CreateMediaRequest interface {
	DirectLink() string
	Kind() string
}
