// Package port defines the interfaces for persistence operations.
// These interfaces abstract the data access layer and provide methods for interacting with various entities in the system.
package port

import (
	"context"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
)

// WallPersister defines the interface for wall persistence operations.
type WallPersister interface {
	// Create creates a new wall in the persistence layer.
	Create(ctx context.Context, wall model.Wall) error
	// Update updates an existing wall in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Wall) error
	// Find retrieves a wall from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Wall, error)
	// FindAll retrieves all walls from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Wall, error)
	// Delete removes a wall from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// WallBlockPersister defines the interface for wall-block association persistence operations.
type WallBlockPersister interface {
	// Create creates a new wall-block association in the persistence layer.
	Create(ctx context.Context, wall model.WallBlock) error
	// Update updates an existing wall-block association in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.WallBlock) error
	// Find retrieves a wall-block association from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.WallBlock, error)
	// FindByWallID retrieves wall-block associations by wall ID.
	FindByWallID(ctx context.Context, id string) ([]*model.WallBlock, error)
	// FindByBlockID retrieves wall-block associations by block ID.
	FindByBlockID(ctx context.Context, id string) ([]*model.WallBlock, error)
	// FindByWallIDAndBlockID retrieves wall-block associations by both wall ID and block ID.
	FindByWallIDAndBlockID(ctx context.Context, wallID string, blockID string) ([]*model.WallBlock, error)
	// Delete removes a wall-block association from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// BlockPersister defines the interface for block persistence operations.
type BlockPersister interface {
	// Create creates a new block in the persistence layer.
	Create(ctx context.Context, wall model.Block) error
	// Update updates an existing block in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Block) error
	// Find retrieves a block from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Block, error)
	// FindAll retrieves all blocks from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Block, error)
	// Delete removes a block from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// BlockProgramPersister defines the interface for block-program association persistence operations.
type BlockProgramPersister interface {
	// Create creates a new block-program association in the persistence layer.
	Create(ctx context.Context, wall model.BlockProgram) error
	// Update updates an existing block-program association in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.BlockProgram) error
	// Find retrieves a block-program association from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.BlockProgram, error)
	// FindByBlockID retrieves block-program associations by block ID.
	FindByBlockID(ctx context.Context, id string) ([]*model.BlockProgram, error)
	// FindByProgramID retrieves block-program associations by program ID.
	FindByProgramID(ctx context.Context, id string) ([]*model.BlockProgram, error)
	// FindByBlockIDAndProgramID retrieves block-program associations by both block ID and program ID.
	FindByBlockIDAndProgramID(ctx context.Context, blockID string, programID string) ([]*model.BlockProgram, error)
	// Delete removes a block-program association from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// ProgramPersister defines the interface for program persistence operations.
type ProgramPersister interface {
	// Create creates a new program in the persistence layer.
	Create(ctx context.Context, wall model.Program) error
	// Update updates an existing program in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Program) error
	// Find retrieves a program from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Program, error)
	// FindAll retrieves all programs from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Program, error)
	// Delete removes a program from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// EpisodePersister defines the interface for episode persistence operations.
type EpisodePersister interface {
	// Create creates a new episode in the persistence layer.
	Create(ctx context.Context, wall model.Episode) error
	// Update updates an existing episode in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Episode) error
	// Find retrieves an episode from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Episode, error)
	// FindByProgramID retrieves episodes by program ID.
	FindByProgramID(ctx context.Context, id string) ([]*model.Episode, error)
	// FindAll retrieves all episodes from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Episode, error)
	// Delete removes an episode from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// MediaPersister defines the interface for media persistence operations.
type MediaPersister interface {
	// Create creates a new media in the persistence layer.
	Create(ctx context.Context, wall model.Media) error
	// Update updates an existing media in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Media) error
	// Find retrieves a media from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Media, error)
	// FindAll retrieves all medias from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Media, error)
	// Delete removes a media from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// TagPersister defines the interface for tag persistence operations.
type TagPersister interface {
	// Create creates a new tag in the persistence layer.
	Create(ctx context.Context, wall model.Tag) error
	// Update updates an existing tag in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Tag) error
	// Find retrieves a tag from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Tag, error)
	// FindAll retrieves all tags from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Tag, error)
	// Delete removes a tag from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// ProgramTagPersister defines the interface for program-tag association persistence operations.
type ProgramTagPersister interface {
	// Create creates a new program-tag association in the persistence layer.
	Create(ctx context.Context, wall model.ProgramTag) error
	// Update updates an existing program-tag association in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.ProgramTag) error
	// Find retrieves a program-tag association from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.ProgramTag, error)
	// FindByTagID retrieves program-tag associations by tag ID.
	FindByTagID(ctx context.Context, id string) ([]*model.ProgramTag, error)
	// FindByProgramID retrieves program-tag associations by program ID.
	FindByProgramID(ctx context.Context, id string) ([]*model.ProgramTag, error)
	// FindByTagIDAndProgramID retrieves program-tag associations by both tag ID and program ID.
	FindByTagIDAndProgramID(ctx context.Context, tagID string, programID string) ([]*model.ProgramTag, error)
	// Delete removes a program-tag association from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// CategoryPersister defines the interface for category persistence operations.
type CategoryPersister interface {
	// Create creates a new category in the persistence layer.
	Create(ctx context.Context, wall model.Category) error
	// Update updates an existing category in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.Category) error
	// Find retrieves a category from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.Category, error)
	// FindAll retrieves all categories from the persistence layer.
	FindAll(ctx context.Context) ([]*model.Category, error)
	// Delete removes a category from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}

// ProgramCategoryPersister defines the interface for program-category association persistence operations.
type ProgramCategoryPersister interface {
	// Create creates a new program-category association in the persistence layer.
	Create(ctx context.Context, wall model.ProgramCategory) error
	// Update updates an existing program-category association in the persistence layer identified by its ID.
	Update(ctx context.Context, id string, updates model.ProgramCategory) error
	// Find retrieves a program-category association from the persistence layer by its ID.
	Find(ctx context.Context, id string) (*model.ProgramCategory, error)
	// FindByCategoryID retrieves program-category associations by category ID.
	FindByCategoryID(ctx context.Context, id string) ([]*model.ProgramCategory, error)
	// FindByProgramID retrieves program-category associations by program ID.
	FindByProgramID(ctx context.Context, id string) ([]*model.ProgramCategory, error)
	// FindByCategoryIDAndProgramID retrieves program-category associations by both category ID and program ID.
	FindByCategoryIDAndProgramID(ctx context.Context, categoryID string, programID string) ([]*model.ProgramCategory, error)
	// Delete removes a program-category association from the persistence layer by its ID.
	Delete(ctx context.Context, id string) error
}
