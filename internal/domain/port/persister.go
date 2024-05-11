package port

import (
	"context"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
)

type WallPersister interface {
	Create(ctx context.Context, wall model.Wall) error
	Update(ctx context.Context, id string, updates model.Wall) error
	Find(ctx context.Context, id string) (*model.Wall, error)
	FindAll(ctx context.Context) ([]*model.Wall, error)
	Delete(ctx context.Context, id string) error
}

type WallBlockPersister interface {
	Create(ctx context.Context, wall model.WallBlock) error
	Update(ctx context.Context, id string, updates model.WallBlock) error
	Find(ctx context.Context, id string) (*model.WallBlock, error)
	FindByWallID(ctx context.Context, id string) ([]*model.WallBlock, error)
	FindByBlockID(ctx context.Context, id string) ([]*model.WallBlock, error)
	FindByWallIDAndBlockID(ctx context.Context, wallID string, blockID string) ([]*model.WallBlock, error)
	Delete(ctx context.Context, id string) error
}

type BlockPersister interface {
	Create(ctx context.Context, wall model.Block) error
	Update(ctx context.Context, id string, updates model.Block) error
	Find(ctx context.Context, id string) (*model.Block, error)
	FindAll(ctx context.Context) ([]*model.Block, error)
	Delete(ctx context.Context, id string) error
}

type BlockProgramPersister interface {
	Create(ctx context.Context, wall model.BlockProgram) error
	Update(ctx context.Context, id string, updates model.BlockProgram) error
	Find(ctx context.Context, id string) (*model.BlockProgram, error)
	FindByBlockID(ctx context.Context, id string) ([]*model.BlockProgram, error)
	FindByProgramID(ctx context.Context, id string) ([]*model.BlockProgram, error)
	FindByBlockIDAndProgramID(ctx context.Context, blockID string, programID string) ([]*model.BlockProgram, error)
	Delete(ctx context.Context, id string) error
}

type ProgramPersister interface {
	Create(ctx context.Context, wall model.Program) error
	Update(ctx context.Context, id string, updates model.Program) error
	Find(ctx context.Context, id string) (*model.Program, error)
	FindAll(ctx context.Context) ([]*model.Program, error)
	Delete(ctx context.Context, id string) error
}

type EpisodePersister interface {
	Create(ctx context.Context, wall model.Episode) error
	Update(ctx context.Context, id string, updates model.Episode) error
	Find(ctx context.Context, id string) (*model.Episode, error)
	FindByProgramID(ctx context.Context, id string) ([]*model.Episode, error)
	FindAll(ctx context.Context) ([]*model.Episode, error)
	Delete(ctx context.Context, id string) error
}

type MediaPersister interface {
	Create(ctx context.Context, wall model.Media) error
	Update(ctx context.Context, id string, updates model.Media) error
	Find(ctx context.Context, id string) (*model.Media, error)
	FindAll(ctx context.Context) ([]*model.Media, error)
	Delete(ctx context.Context, id string) error
}

type TagPersister interface {
	Create(ctx context.Context, wall model.Tag) error
	Update(ctx context.Context, id string, updates model.Tag) error
	Find(ctx context.Context, id string) (*model.Tag, error)
	FindAll(ctx context.Context) ([]*model.Tag, error)
	Delete(ctx context.Context, id string) error
}

type ProgramTagPersister interface {
	Create(ctx context.Context, wall model.ProgramTag) error
	Update(ctx context.Context, id string, updates model.ProgramTag) error
	Find(ctx context.Context, id string) (*model.ProgramTag, error)
	FindByTagID(ctx context.Context, id string) ([]*model.ProgramTag, error)
	FindByProgramID(ctx context.Context, id string) ([]*model.ProgramTag, error)
	FindByTagIDAndProgramID(ctx context.Context, tagID string, programID string) ([]*model.ProgramTag, error)
	Delete(ctx context.Context, id string) error
}

type CategoryPersister interface {
	Create(ctx context.Context, wall model.Category) error
	Update(ctx context.Context, id string, updates model.Category) error
	Find(ctx context.Context, id string) (*model.Category, error)
	FindAll(ctx context.Context) ([]*model.Category, error)
	Delete(ctx context.Context, id string) error
}

type ProgramCategoryPersister interface {
	Create(ctx context.Context, wall model.ProgramCategory) error
	Update(ctx context.Context, id string, updates model.ProgramCategory) error
	Find(ctx context.Context, id string) (*model.ProgramCategory, error)
	FindByCategoryID(ctx context.Context, id string) ([]*model.ProgramCategory, error)
	FindByProgramID(ctx context.Context, id string) ([]*model.ProgramCategory, error)
	FindByCategoryIDAndProgramID(ctx context.Context, categoryID string, programID string) ([]*model.ProgramCategory, error)
	Delete(ctx context.Context, id string) error
}
