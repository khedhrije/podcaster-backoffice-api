// Package api provides functionality for managing programs.
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

// CreateProgramRequest represents the interface for creating programs.
type CreateProgramRequest interface {
	Name() string
	Description() string
}

// UpdateProgramRequest represents the interface for updating programs.
type UpdateProgramRequest interface {
	Name() string
	Description() string
}

// Program represents the interface for managing programs.
type Program interface {
	Create(ctx context.Context, program CreateProgramRequest) error
	Update(ctx context.Context, uuid string, updates UpdateProgramRequest) error
	Find(ctx context.Context, uuid string) (*pkg.ProgramResponse, error)
	FindAll(ctx context.Context) ([]*pkg.ProgramResponse, error)
	Delete(ctx context.Context, uuid string) error
	FindEpisodes(ctx context.Context, uuid string) ([]*pkg.EpisodeResponse, error)
	FindTags(ctx context.Context, uuid string) ([]*pkg.TagResponse, error)
	FindCats(ctx context.Context, uuid string) ([]*pkg.CategoryResponse, error)
	OverwriteCategories(ctx context.Context, programID string, cats []string) error
	OverwriteTags(ctx context.Context, programID string, tags []string) error
}

// programApi is an implementation of the Program interface.
type programApi struct {
	programAdapter    port.ProgramPersister
	episodeAdapter    port.EpisodePersister
	programTagAdapter port.ProgramTagPersister
	tagAdapter        port.TagPersister
	programCatAdapter port.ProgramCategoryPersister
	catAdapter        port.CategoryPersister
}

// NewProgramApi creates a new instance of Program.
// It takes adapters for program, episode, tag, and category persistence as dependencies.
func NewProgramApi(
	programAdapter port.ProgramPersister,
	episodeAdapter port.EpisodePersister,
	programTagAdapter port.ProgramTagPersister,
	tagAdapter port.TagPersister,
	programCatAdapter port.ProgramCategoryPersister,
	catAdapter port.CategoryPersister,
) Program {
	return programApi{
		programAdapter:    programAdapter,
		episodeAdapter:    episodeAdapter,
		programTagAdapter: programTagAdapter,
		tagAdapter:        tagAdapter,
		programCatAdapter: programCatAdapter,
		catAdapter:        catAdapter,
	}
}

// Create creates a new program.
// It takes the context and CreateProgramRequest, and returns an error if any.
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
	// Call adapter
	if err := api.programAdapter.Create(ctx, program); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("program", program).Msg("error while creating program")
		return fmt.Errorf("error occurred while creating program: %w", err)
	}

	return nil
}

// createProgramRequestValidation validates the creation request.
// It takes the context and CreateProgramRequest, and returns a slice of ValidationErrors.
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
// It takes the context, program UUID, and UpdateProgramRequest, and returns an error if any.
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
	// Call adapter
	if err := api.programAdapter.Update(ctx, uuid, program); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("program", program).Msg("error while updating program")
		return fmt.Errorf("error occurred while updating program: %w", err)
	}

	return nil
}

// updateProgramRequestValidation validates the update request.
// It takes the context, program UUID, and UpdateProgramRequest, and returns a slice of ValidationErrors.
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
// It takes the context and program UUID, and returns a ProgramResponse or an error.
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
	// Return result
	return response, nil
}

// FindAll finds all programs.
// It takes the context and returns a slice of ProgramResponse or an error.
func (api programApi) FindAll(ctx context.Context) ([]*pkg.ProgramResponse, error) {
	// Call adapter
	programSlice, err := api.programAdapter.FindAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error while finding programs")
		return nil, fmt.Errorf("error occurred while finding programs: %w", err)
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
	// Return result
	return response, nil
}

// Delete deletes a program by UUID.
// It takes the context and program UUID, and returns an error if any.
func (api programApi) Delete(ctx context.Context, uuid string) error {
	if err := api.programAdapter.Delete(ctx, uuid); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while deleting program")
		return fmt.Errorf("error occurred while deleting program: %w", err)
	}
	return nil
}

// FindEpisodes finds a program's episodes.
// It takes the context and program UUID, and returns a slice of EpisodeResponse or an error.
func (api programApi) FindEpisodes(ctx context.Context, uuid string) ([]*pkg.EpisodeResponse, error) {
	// Call adapter
	episodes, err := api.episodeAdapter.FindByProgramID(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding program episodes")
		return nil, fmt.Errorf("error occurred while finding program episodes: %w", err)
	}

	// Map to response
	var response []*pkg.EpisodeResponse
	for _, episode := range episodes {
		response = append(response, &pkg.EpisodeResponse{
			ID:          episode.ID,
			Name:        episode.Name,
			Description: episode.Description,
		})
	}

	// Return result
	return response, nil
}

// FindTags finds a program's tags.
// It takes the context and program UUID, and returns a slice of TagResponse or an error.
func (api programApi) FindTags(ctx context.Context, uuid string) ([]*pkg.TagResponse, error) {
	// Call adapter
	associations, err := api.programTagAdapter.FindByProgramID(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding program tags")
		return nil, fmt.Errorf("error occurred while finding program's tags: %w", err)
	}

	// Map to response
	var response []*pkg.TagResponse
	for _, association := range associations {
		tag, err := api.tagAdapter.Find(ctx, association.TagID)
		if err != nil {
			return nil, err
		}
		response = append(response, &pkg.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
		})
	}

	// Return result
	return response, nil
}

// FindCats finds a program's categories.
// It takes the context and program UUID, and returns a slice of CategoryResponse or an error.
func (api programApi) FindCats(ctx context.Context, uuid string) ([]*pkg.CategoryResponse, error) {
	// Call adapter
	associations, err := api.programCatAdapter.FindByProgramID(ctx, uuid)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("uuid", uuid).Msg("error while finding program categories")
		return nil, fmt.Errorf("error occurred while finding program's categories: %w", err)
	}

	// Map to response
	var response []*pkg.CategoryResponse
	for _, association := range associations {
		cat, err := api.catAdapter.Find(ctx, association.CategoryID)
		if err != nil {
			return nil, err
		}
		response = append(response, &pkg.CategoryResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
		})
	}

	// Return result
	return response, nil
}

// OverwriteCategories overwrites the categories associated with a program.
// It takes the context, program ID, and a slice of category IDs, and returns an error if any.
func (api programApi) OverwriteCategories(ctx context.Context, programID string, catIDs []string) error {
	// Find all existing associations by programID
	associations, err := api.programCatAdapter.FindByProgramID(ctx, programID)
	if err != nil {
		return err
	}

	// Remove all existing associations for the program
	for _, association := range associations {
		if err = api.programCatAdapter.Delete(ctx, association.ID); err != nil {
			return err
		}
	}

	// Create all new associations
	for _, catID := range catIDs {
		programCategory := model.ProgramCategory{
			ID:         uuid.New().String(),
			ProgramID:  programID,
			CategoryID: catID,
		}

		// Call adapter to create program category
		if err := api.programCatAdapter.Create(ctx, programCategory); err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("programCategory", programCategory).Msg("error while creating program category")
			return fmt.Errorf("error occurred while creating program category: %w", err)
		}
	}

	return nil
}

// OverwriteTags overwrites the tags associated with a program.
// It takes the context, program ID, and a slice of tag IDs, and returns an error if any.
func (api programApi) OverwriteTags(ctx context.Context, programID string, tagIDs []string) error {
	// Find all existing associations by programID
	associations, err := api.programTagAdapter.FindByProgramID(ctx, programID)
	if err != nil {
		return err
	}

	// Remove all existing associations for the program
	for _, association := range associations {
		if err = api.programTagAdapter.Delete(ctx, association.ID); err != nil {
			return err
		}
	}

	// Create all new associations
	for _, tagID := range tagIDs {
		programTag := model.ProgramTag{
			ID:        uuid.New().String(),
			ProgramID: programID,
			TagID:     tagID,
		}

		// Call adapter to create program tag
		if err := api.programTagAdapter.Create(ctx, programTag); err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("programTag", programTag).Msg("error while creating program tag")
			return fmt.Errorf("error occurred while creating program tag: %w", err)
		}
	}

	return nil
}
