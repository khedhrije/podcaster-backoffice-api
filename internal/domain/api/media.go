// Package api provides functionality for managing medias.
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

// mediaApi is an implementation of the Media interface.
type mediaApi struct {
	mediaAdapter port.MediaPersister
}

// NewMediaApi creates a new instance of Media.
// It takes a MediaPersister as a dependency.
func NewMediaApi(mediaAdapter port.MediaPersister) Media {
	return &mediaApi{
		mediaAdapter: mediaAdapter,
	}
}

// Create creates a new media.
// It takes the context and CreateMediaRequest, and returns an error if any.
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
		EpisodeID:  req.EpisodeID(),
	}
	// Call adapter
	if err := api.mediaAdapter.Create(ctx, media); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("media", media).Msg("error while creating media")
		return fmt.Errorf("error occurred while creating media: %w", err)
	}

	return nil
}

// createMediaRequestValidation validates the creation request.
// It takes the context and CreateMediaRequest, and returns a slice of ValidationErrors.
func createMediaRequestValidation(ctx context.Context, req CreateMediaRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.DirectLink() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "directLink", Message: "is required"})
	}
	if req.Kind() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "kind", Message: "is required"})
	}
	if req.EpisodeID() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "episodeID", Message: "cannot be empty"})
	}
	return vErrs
}

// Update updates an existing media.
// It takes the context, media UUID, and UpdateMediaRequest, and returns an error if any.
func (api mediaApi) Update(ctx context.Context, uuid string, updates UpdateMediaRequest) error {
	// Validate request
	vErrs := updateMediaRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	media := model.Media{}
	if updates.DirectLink() != "" {
		media.DirectLink = updates.DirectLink()
	}
	if updates.Kind() != "" {
		media.Kind = updates.Kind()
	}
	if updates.EpisodeID() != "" {
		media.EpisodeID = updates.EpisodeID()
	}

	// Call adapter
	if err := api.mediaAdapter.Update(ctx, uuid, media); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("media", media).Msg("error while updating media")
		return fmt.Errorf("error occurred while updating media: %w", err)
	}

	return nil
}

// updateMediaRequestValidation validates the update request.
// It takes the context, media UUID, and UpdateMediaRequest, and returns a slice of ValidationErrors.
func updateMediaRequestValidation(ctx context.Context, uuid string, req UpdateMediaRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if uuid == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "uuid", Message: "cannot be empty"})
	}
	return vErrs
}

// Find finds a media by UUID.
// It takes the context and media UUID, and returns a MediaResponse or an error.
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
		EpisodeID:  media.EpisodeID,
	}
	// Return result
	return response, nil
}

// FindAll finds all medias.
// It takes the context and returns a slice of MediaResponse or an error.
func (api mediaApi) FindAll(ctx context.Context) ([]*pkg.MediaResponse, error) {
	// Call adapter
	mediaSlice, err := api.mediaAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding medias")
		return nil, fmt.Errorf("error occurred while finding medias: %w", err)
	}

	// Map to response
	var response []*pkg.MediaResponse
	for _, media := range mediaSlice {
		response = append(response, &pkg.MediaResponse{
			ID:         media.ID,
			DirectLink: media.DirectLink,
			Kind:       media.Kind,
			EpisodeID:  media.EpisodeID,
		})
	}
	// Return result
	return response, nil
}

// Delete deletes a media by UUID.
// It takes the context and media UUID, and returns an error if any.
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
	EpisodeID() string
}

// UpdateMediaRequest represents the interface for updating medias.
type UpdateMediaRequest interface {
	DirectLink() string
	Kind() string
	EpisodeID() string
}
