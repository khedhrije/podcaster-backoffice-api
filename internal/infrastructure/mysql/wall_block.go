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
func NewWallBlockAdapter(client *client) port.WallBlockPersister {
	return &wallBlockAdapter{
		client: client,
	}
}

// Create inserts a new wallBlock record into the database.
func (adapter *wallBlockAdapter) Create(ctx context.Context, wallBlock model.WallBlock) error {
	// SQL query to insert a new wallBlock record
	const query = `
        INSERT INTO wall_block (UUID, wallUUID, blockUUID, position)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:wallUUID), UUID_TO_BIN(:blockUUID), :position)
    `
	// Convert domain model to database model
	var wallBlockDB WallBlockDB
	wallBlockDB.FromDomainModel(wallBlock)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallBlockDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a wallBlock record from the database based on its UUID.
func (adapter *wallBlockAdapter) Delete(ctx context.Context, wallBlockUUID string) error {
	// SQL query to delete a wallBlock record by UUID
	const query = `
        DELETE FROM wall_block WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, wallBlockUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing wallBlock record in the database.
func (adapter *wallBlockAdapter) Update(ctx context.Context, wallBlockUUID string, updates model.WallBlock) error {
	// SQL query to update a wallBlock record
	const query = `
        UPDATE wall_block SET 
                             wallUUID = UUID_TO_BIN(:wallUUID), 
                             blockUUID = UUID_TO_BIN(:blockUUID), 
                             position = :position
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = wallBlockUUID
	// Convert domain model to database model
	var wallBlockDB WallBlockDB
	wallBlockDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallBlockDB)
	if err != nil {
		return err
	}
	return nil
}

// Find retrieves a wallBlock record from the database by its UUID.
func (adapter *wallBlockAdapter) Find(ctx context.Context, wallBlockUUID string) (*model.WallBlock, error) {
	// SQL query to select a wallBlock record by UUID
	const query = `
        SELECT * FROM wall_block WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve result
	var wallBlockDB WallBlockDB
	if err := adapter.client.db.GetContext(ctx, &wallBlockDB, query, wallBlockUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if wallBlockDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("wallBlock with ID %s not found", wallBlockUUID)
	}
	// Convert database model to domain model
	result := wallBlockDB.ToDomainModel()
	return &result, nil
}

// FindByWallID retrieves all wallBlock records from the database for a given wall ID.
func (adapter *wallBlockAdapter) FindByWallID(ctx context.Context, wallID string) ([]*model.WallBlock, error) {
	// SQL query to select wallBlock records by wall ID
	const query = `
        SELECT * FROM wall_block WHERE wallUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var wallBlocksDB []*WallBlockDB
	if err := adapter.client.db.SelectContext(ctx, &wallBlocksDB, query, wallID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var wallBlocks []*model.WallBlock
	for _, wallBlockDB := range wallBlocksDB {
		mappedWallBlock := wallBlockDB.ToDomainModel()
		wallBlocks = append(wallBlocks, &mappedWallBlock)
	}
	return wallBlocks, nil
}

// FindByBlockID retrieves all wallBlock records from the database for a given block ID.
func (adapter *wallBlockAdapter) FindByBlockID(ctx context.Context, blockID string) ([]*model.WallBlock, error) {
	// SQL query to select wallBlock records by block ID
	const query = `
        SELECT * FROM wall_block WHERE blockUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var wallBlocksDB []*WallBlockDB
	if err := adapter.client.db.SelectContext(ctx, &wallBlocksDB, query, blockID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var wallBlocks []*model.WallBlock
	for _, wallBlockDB := range wallBlocksDB {
		mappedWallBlock := wallBlockDB.ToDomainModel()
		wallBlocks = append(wallBlocks, &mappedWallBlock)
	}
	return wallBlocks, nil
}

// FindByWallIDAndBlockID retrieves all wallBlock records from the database for a given wall ID and block ID.
func (adapter *wallBlockAdapter) FindByWallIDAndBlockID(ctx context.Context, wallID, blockID string) ([]*model.WallBlock, error) {
	// SQL query to select wallBlock records by wall ID and block ID
	const query = `
        SELECT * FROM wall_block WHERE wallUUID = UUID_TO_BIN(?) AND blockUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var wallBlocksDB []*WallBlockDB
	if err := adapter.client.db.SelectContext(ctx, &wallBlocksDB, query, wallID, blockID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
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
func (db *WallBlockDB) ToDomainModel() model.WallBlock {
	return model.WallBlock{
		ID:       db.UUID.String(),
		WallID:   db.WallID.String(),
		BlockID:  db.BlockID.String(),
		Position: db.Position,
	}
}

// FromDomainModel converts a model.WallBlock domain model to a WallBlockDB database model.
func (db *WallBlockDB) FromDomainModel(domain model.WallBlock) {
	db.UUID = uuid.MustParse(domain.ID)
	db.WallID = uuid.MustParse(domain.WallID)
	db.BlockID = uuid.MustParse(domain.BlockID)
	db.Position = domain.Position
}
