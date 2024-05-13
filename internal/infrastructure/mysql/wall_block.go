// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// wallBlockAdapter is a struct that acts as an adapter for interacting with
// the wall_block data in the MySQL database.
type wallBlockAdapter struct {
	client *client
}

// NewWallBlockAdapter creates a new wallBlock adapter with the provided MySQL client.
// It returns an implementation of the WallBlockPersister interface.
func NewWallBlockAdapter(client *client) port.WallBlockPersister {
	return &wallBlockAdapter{
		client: client,
	}
}

// Create inserts a new wallBlock record into the database.
// It takes a context and a model.WallBlock, and returns an error if the operation fails.
func (adapter *wallBlockAdapter) Create(ctx context.Context, wallBlock model.WallBlock) error {
	const query = `
        INSERT INTO wall_block (UUID, wallUUID, blockUUID, position)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:wallUUID), UUID_TO_BIN(:blockUUID), :position)
    `
	var wallBlockDB WallBlockDB
	wallBlockDB.FromDomainModel(wallBlock)
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallBlockDB)
	return err
}

// Delete removes a wallBlock record from the database based on its UUID.
// It takes a context and the wallBlock's UUID, and returns an error if the operation fails.
func (adapter *wallBlockAdapter) Delete(ctx context.Context, wallBlockUUID string) error {
	const query = `
        DELETE FROM wall_block WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, wallBlockUUID)
	return err
}

// Update updates an existing wallBlock record in the database.
// It takes a context, the wallBlock's UUID, and the updated model.WallBlock, and returns an error if the operation fails.
func (adapter *wallBlockAdapter) Update(ctx context.Context, wallBlockUUID string, updates model.WallBlock) error {
	const query = `
        UPDATE wall_block SET 
                             wallUUID = UUID_TO_BIN(:wallUUID), 
                             blockUUID = UUID_TO_BIN(:blockUUID), 
                             position = :position
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = wallBlockUUID
	var wallBlockDB WallBlockDB
	wallBlockDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallBlockDB)
	return err
}

// Find retrieves a wallBlock record from the database by its UUID.
// It takes a context and the wallBlock's UUID, and returns a model.WallBlock and an error if the operation fails.
func (adapter *wallBlockAdapter) Find(ctx context.Context, wallBlockUUID string) (*model.WallBlock, error) {
	const query = `
        SELECT * FROM wall_block WHERE UUID = UUID_TO_BIN(?)
    `
	var wallBlockDB WallBlockDB
	if err := adapter.client.db.GetContext(ctx, &wallBlockDB, query, wallBlockUUID); err != nil {
		return nil, err
	}
	if wallBlockDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("wallBlock with ID %s not found", wallBlockUUID)
	}
	result := wallBlockDB.ToDomainModel()
	return &result, nil
}

// FindByWallID retrieves all wallBlock records from the database for a given wall ID.
// It takes a context and the wall's ID, and returns a slice of model.WallBlock and an error if the operation fails.
func (adapter *wallBlockAdapter) FindByWallID(ctx context.Context, wallID string) ([]*model.WallBlock, error) {
	const query = `
        SELECT * FROM wall_block WHERE wallUUID = UUID_TO_BIN(?)
    `
	var wallBlocksDB []*WallBlockDB
	if err := adapter.client.db.SelectContext(ctx, &wallBlocksDB, query, wallID); err != nil {
		return nil, err
	}
	var wallBlocks []*model.WallBlock
	for _, wallBlockDB := range wallBlocksDB {
		mappedWallBlock := wallBlockDB.ToDomainModel()
		wallBlocks = append(wallBlocks, &mappedWallBlock)
	}
	return wallBlocks, nil
}

// FindByBlockID retrieves all wallBlock records from the database for a given block ID.
// It takes a context and the block's ID, and returns a slice of model.WallBlock and an error if the operation fails.
func (adapter *wallBlockAdapter) FindByBlockID(ctx context.Context, blockID string) ([]*model.WallBlock, error) {
	const query = `
        SELECT * FROM wall_block WHERE blockUUID = UUID_TO_BIN(?)
    `
	var wallBlocksDB []*WallBlockDB
	if err := adapter.client.db.SelectContext(ctx, &wallBlocksDB, query, blockID); err != nil {
		return nil, err
	}
	var wallBlocks []*model.WallBlock
	for _, wallBlockDB := range wallBlocksDB {
		mappedWallBlock := wallBlockDB.ToDomainModel()
		wallBlocks = append(wallBlocks, &mappedWallBlock)
	}
	return wallBlocks, nil
}

// FindByWallIDAndBlockID retrieves all wallBlock records from the database for a given wall ID and block ID.
// It takes a context, the wall's ID, and the block's ID, and returns a slice of model.WallBlock and an error if the operation fails.
func (adapter *wallBlockAdapter) FindByWallIDAndBlockID(ctx context.Context, wallID, blockID string) ([]*model.WallBlock, error) {
	const query = `
        SELECT * FROM wall_block WHERE wallUUID = UUID_TO_BIN(?) AND blockUUID = UUID_TO_BIN(?)
    `
	var wallBlocksDB []*WallBlockDB
	if err := adapter.client.db.SelectContext(ctx, &wallBlocksDB, query, wallID, blockID); err != nil {
		return nil, err
	}
	var wallBlocks []*model.WallBlock
	for _, wallBlockDB := range wallBlocksDB {
		mappedWallBlock := wallBlockDB.ToDomainModel()
		wallBlocks = append(wallBlocks, &mappedWallBlock)
	}
	return wallBlocks, nil
}

// WallBlockDB is a struct representing the wallBlock database model.
type WallBlockDB struct {
	UUID     uuid.UUID `db:"UUID"`
	WallID   uuid.UUID `db:"wallUUID"`
	BlockID  uuid.UUID `db:"blockUUID"`
	Position int       `db:"position"`
}

// ToDomainModel converts a WallBlockDB database model to a model.WallBlock domain model.
// It returns the corresponding model.WallBlock.
func (db *WallBlockDB) ToDomainModel() model.WallBlock {
	return model.WallBlock{
		ID:       db.UUID.String(),
		WallID:   db.WallID.String(),
		BlockID:  db.BlockID.String(),
		Position: db.Position,
	}
}

// FromDomainModel converts a model.WallBlock domain model to a WallBlockDB database model.
// It sets the fields of the WallBlockDB based on the given model.WallBlock.
func (db *WallBlockDB) FromDomainModel(domain model.WallBlock) {
	db.UUID = uuid.MustParse(domain.ID)
	db.WallID = uuid.MustParse(domain.WallID)
	db.BlockID = uuid.MustParse(domain.BlockID)
	db.Position = domain.Position
}
