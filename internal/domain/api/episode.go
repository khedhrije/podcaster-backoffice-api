// Package episode provides functionality for managing episodes.
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

// Episode represents the interface for managing episodes.
type Episode interface {
	Create(ctx context.Context, episode CreateEpisodeRequest) error
	Update(ctx context.Context, uuid string, updates UpdateEpisodeRequest) error
	Find(ctx context.Context, uuid string) (*pkg.EpisodeResponse, error)
	FindAll(ctx context.Context) ([]*pkg.EpisodeResponse, error)
	Delete(ctx context.Context, uuid string) error
}

// UpdateEpisodeRequest represents the interface for updating episodes.
type UpdateEpisodeRequest interface {
	Name() string
	Description() string
}

// episodeApi is an implementation of the Episode interface.
type episodeApi struct {
	episodeAdapter port.EpisodePersister
}

// NewEpisodeApi creates a new instance of Episode.
func NewEpisodeApi(episodeAdapter port.EpisodePersister) Episode {
	return &episodeApi{
		episodeAdapter: episodeAdapter,
	}
}

// Create creates a new episode.
func (api episodeApi) Create(ctx context.Context, req CreateEpisodeRequest) error {
	// Validate request
	vErrs := createEpisodeRequestValidation(ctx, req)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("request", req).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	episode := model.Episode{
		ID:          uuid.New().String(),
		Name:        req.Name(),
		Description: req.Description(),
	}
	// call adapter
	if err := api.episodeAdapter.Create(ctx, episode); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("episode", episode).Msg("error while creating episode")
		return fmt.Errorf("error occurred while creating episode: %w", err)
	}

	return nil
}

// createRequestValidation validates the creation request.
func createEpisodeRequestValidation(ctx context.Context, req CreateEpisodeRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	return vErrs
}

// Update updates an existing episode.
func (api episodeApi) Update(ctx context.Context, uuid string, updates UpdateEpisodeRequest) error {
	// Validate request
	vErrs := updateEpisodeRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	episode := model.Episode{
		Name:        updates.Name(),
		Description: updates.Description(),
	}
	// call adapter
	if err := api.episodeAdapter.Update(ctx, uuid, episode); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("episode", episode).Msg("error while updating episode")
		return fmt.Errorf("error occurred while updating episode: %w", err)
	}

	return nil
}

// updateRequestValidation validates the update request.
func updateEpisodeRequestValidation(ctx context.Context, uuid string, req UpdateEpisodeRequest) model.ValidationErrors {
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

// Find finds a episode by UUID.
func (api episodeApi) Find(ctx context.Context, uuid string) (*pkg.EpisodeResponse, error) {

	// Call adapter
	episode, err := api.episodeAdapter.Find(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding episode")
		return nil, fmt.Errorf("error occurred while finding episode: %w", err)
	}

	// Map to response
	response := &pkg.EpisodeResponse{
		ID:          episode.ID,
		Name:        episode.Name,
		Description: episode.Description,
	}
	// return result
	return response, nil
}

// FindAll finds all episodes.
func (api episodeApi) FindAll(ctx context.Context) ([]*pkg.EpisodeResponse, error) {
	// Call adapter
	episodeSlice, err := api.episodeAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding episode")
		return nil, fmt.Errorf("error occurred while finding episode: %w", err)
	}

	// Map to response
	var response []*pkg.EpisodeResponse
	for _, episode := range episodeSlice {
		response = append(response, &pkg.EpisodeResponse{
			ID:          episode.ID,
			Name:        episode.Name,
			Description: episode.Description,
		})
	}
	// return result
	return response, nil
}

// Delete deletes a episode by UUID.
func (api episodeApi) Delete(ctx context.Context, uuid string) error {
	if err := api.episodeAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting episode")
		return fmt.Errorf("error occurred while deleting episode: %w", err)
	}
	return nil
}

// CreateEpisodeRequest represents the interface for creating episodes.
type CreateEpisodeRequest interface {
	Name() string
	Description() string
}
