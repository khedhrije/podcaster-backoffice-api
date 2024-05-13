// Package api provides functionality for managing episodes.
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

// episodeApi is an implementation of the Episode interface.
type episodeApi struct {
	episodeAdapter port.EpisodePersister
}

// NewEpisodeApi creates a new instance of Episode.
// It takes an EpisodePersister as a dependency.
func NewEpisodeApi(episodeAdapter port.EpisodePersister) Episode {
	return &episodeApi{
		episodeAdapter: episodeAdapter,
	}
}

// Create creates a new episode.
// It takes the context and CreateEpisodeRequest, and returns an error if any.
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
		ProgramID:   req.ProgramID(),
		Position:    req.Position(),
	}
	// Call adapter
	if err := api.episodeAdapter.Create(ctx, episode); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("episode", episode).Msg("error while creating episode")
		return fmt.Errorf("error occurred while creating episode: %w", err)
	}

	return nil
}

// createEpisodeRequestValidation validates the creation request.
// It takes the context and CreateEpisodeRequest, and returns a slice of ValidationErrors.
func createEpisodeRequestValidation(ctx context.Context, req CreateEpisodeRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if req.Name() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "name", Message: "is required"})
	}
	if req.Description() == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "description", Message: "is required"})
	}
	if req.Position() <= 0 {
		vErrs = append(vErrs, model.ValidationError{Field: "position", Message: "should be greater than 0"})
	}
	return vErrs
}

// Update updates an existing episode.
// It takes the context, episode UUID, and UpdateEpisodeRequest, and returns an error if any.
func (api episodeApi) Update(ctx context.Context, uuid string, updates UpdateEpisodeRequest) error {
	// Validate request
	vErrs := updateEpisodeRequestValidation(ctx, uuid, updates)
	if len(vErrs) > 0 {
		log.Ctx(ctx).Error().Err(vErrs).Interface("updates", updates).Msg("request was not validated")
		return fmt.Errorf("request was not validated: %w", vErrs)
	}
	// Map to domain model
	episode := model.Episode{}
	if updates.Name() != "" {
		episode.Name = updates.Name()
	}
	if updates.Description() != "" {
		episode.Description = updates.Description()
	}
	if updates.ProgramID() != "" {
		episode.ProgramID = updates.ProgramID()
	}
	if updates.Position() != 0 {
		episode.Position = updates.Position()
	}
	// Call adapter
	if err := api.episodeAdapter.Update(ctx, uuid, episode); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("episode", episode).Msg("error while updating episode")
		return fmt.Errorf("error occurred while updating episode: %w", err)
	}

	return nil
}

// updateEpisodeRequestValidation validates the update request.
// It takes the context, episode UUID, and UpdateEpisodeRequest, and returns a slice of ValidationErrors.
func updateEpisodeRequestValidation(ctx context.Context, uuid string, req UpdateEpisodeRequest) model.ValidationErrors {
	var vErrs []model.ValidationError
	if uuid == "" {
		vErrs = append(vErrs, model.ValidationError{Field: "uuid", Message: "cannot be empty"})
	}
	return vErrs
}

// Find finds an episode by UUID.
// It takes the context and episode UUID, and returns an EpisodeResponse or an error.
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
		ProgramID:   episode.ProgramID,
		Position:    episode.Position,
	}
	// Return result
	return response, nil
}

// FindAll finds all episodes.
// It takes the context and returns a slice of EpisodeResponse or an error.
func (api episodeApi) FindAll(ctx context.Context) ([]*pkg.EpisodeResponse, error) {
	// Call adapter
	episodeSlice, err := api.episodeAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding episodes")
		return nil, fmt.Errorf("error occurred while finding episodes: %w", err)
	}

	// Map to response
	var response []*pkg.EpisodeResponse
	for _, episode := range episodeSlice {
		response = append(response, &pkg.EpisodeResponse{
			ID:          episode.ID,
			Name:        episode.Name,
			Description: episode.Description,
			ProgramID:   episode.ProgramID,
			Position:    episode.Position,
		})
	}
	// Return result
	return response, nil
}

// Delete deletes an episode by UUID.
// It takes the context and episode UUID, and returns an error if any.
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
	ProgramID() string
	Position() int
}

// UpdateEpisodeRequest represents the interface for updating episodes.
type UpdateEpisodeRequest interface {
	Name() string
	Description() string
	ProgramID() string
	Position() int
}
