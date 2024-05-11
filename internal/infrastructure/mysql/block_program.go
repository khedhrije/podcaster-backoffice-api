package mysql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// blockProgramAdapter is a struct that acts as an adapter for interacting with
// the block_program data in the MySQL database.
type blockProgramAdapter struct {
	client *client
}

// NewBlockProgramAdapter creates a new blockProgram adapter with the provided MySQL client.
func NewBlockProgramAdapter(client *client) port.BlockProgramPersister {
	return &blockProgramAdapter{
		client: client,
	}
}

// Create inserts a new blockProgram record into the database.
func (adapter *blockProgramAdapter) Create(ctx context.Context, blockProgram model.BlockProgram) error {
	// SQL query to insert a new blockProgram record
	const query = `
        INSERT INTO block_program (UUID, blockUUID, programUUID, position)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:blockUUID), UUID_TO_BIN(:programUUID), :position)
    `
	// Convert domain model to database model
	var blockProgramDB BlockProgramDB
	blockProgramDB.FromDomainModel(blockProgram)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockProgramDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a blockProgram record from the database based on its UUID.
func (adapter *blockProgramAdapter) Delete(ctx context.Context, blockProgramUUID string) error {
	// SQL query to delete a blockProgram record by UUID
	const query = `
        DELETE FROM block_program WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, blockProgramUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing blockProgram record in the database.
func (adapter *blockProgramAdapter) Update(ctx context.Context, blockProgramUUID string, updates model.BlockProgram) error {
	// SQL query to update a blockProgram record
	const query = `
        UPDATE block_program SET 
                             blockUUID = UUID_TO_BIN(:blockUUID), 
                             programUUID = UUID_TO_BIN(:programUUID), 
                             position = :position
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = blockProgramUUID
	// Convert domain model to database model
	var blockProgramDB BlockProgramDB
	blockProgramDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockProgramDB)
	if err != nil {
		return err
	}
	return nil
}

// Find retrieves a blockProgram record from the database by its UUID.
func (adapter *blockProgramAdapter) Find(ctx context.Context, blockProgramUUID string) (*model.BlockProgram, error) {
	// SQL query to select a blockProgram record by UUID
	const query = `
        SELECT * FROM block_program WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve result
	var blockProgramDB BlockProgramDB
	if err := adapter.client.db.GetContext(ctx, &blockProgramDB, query, blockProgramUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if blockProgramDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("blockProgram with ID %s not found", blockProgramUUID)
	}
	// Convert database model to domain model
	result := blockProgramDB.ToDomainModel()
	return &result, nil
}

// FindByBlockID retrieves all blockProgram records from the database for a given block ID.
func (adapter *blockProgramAdapter) FindByBlockID(ctx context.Context, blockID string) ([]*model.BlockProgram, error) {
	// SQL query to select blockProgram records by block ID
	const query = `
        SELECT * FROM block_program WHERE blockUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var blockProgramsDB []*BlockProgramDB
	if err := adapter.client.db.SelectContext(ctx, &blockProgramsDB, query, blockID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var blockPrograms []*model.BlockProgram
	for _, blockProgramDB := range blockProgramsDB {
		mappedBlockProgram := blockProgramDB.ToDomainModel()
		blockPrograms = append(blockPrograms, &mappedBlockProgram)
	}
	return blockPrograms, nil
}

// FindByProgramID retrieves all blockProgram records from the database for a given program ID.
func (adapter *blockProgramAdapter) FindByProgramID(ctx context.Context, programID string) ([]*model.BlockProgram, error) {
	// SQL query to select blockProgram records by program ID
	const query = `
        SELECT * FROM block_program WHERE programUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var blockProgramsDB []*BlockProgramDB
	if err := adapter.client.db.SelectContext(ctx, &blockProgramsDB, query, programID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var blockPrograms []*model.BlockProgram
	for _, blockProgramDB := range blockProgramsDB {
		mappedBlockProgram := blockProgramDB.ToDomainModel()
		blockPrograms = append(blockPrograms, &mappedBlockProgram)
	}
	return blockPrograms, nil
}

// FindByBlockIDAndProgramID retrieves all blockProgram records from the database for a given block ID and program ID.
func (adapter *blockProgramAdapter) FindByBlockIDAndProgramID(ctx context.Context, blockID, programID string) ([]*model.BlockProgram, error) {
	// SQL query to select blockProgram records by block ID and program ID
	const query = `
        SELECT * FROM block_program WHERE blockUUID = UUID_TO_BIN(?) AND programUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var blockProgramsDB []*BlockProgramDB
	if err := adapter.client.db.SelectContext(ctx, &blockProgramsDB, query, blockID, programID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var blockPrograms []*model.BlockProgram
	for _, blockProgramDB := range blockProgramsDB {
		mappedBlockProgram := blockProgramDB.ToDomainModel()
		blockPrograms = append(blockPrograms, &mappedBlockProgram)
	}
	return blockPrograms, nil
}

// BlockProgramDB is a struct representing the blockProgram database model.
type BlockProgramDB struct {
	UUID      uuid.UUID `db:"UUID"`
	BlockID   uuid.UUID `db:"blockUUID"`
	ProgramID uuid.UUID `db:"programUUID"`
	Position  int       `db:"position"`
}

// ToDomainModel converts a BlockProgramDB database model to a model.BlockProgram domain model.
func (db *BlockProgramDB) ToDomainModel() model.BlockProgram {
	return model.BlockProgram{
		ID:        db.UUID.String(),
		BlockID:   db.BlockID.String(),
		ProgramID: db.ProgramID.String(),
		Position:  db.Position,
	}
}

// FromDomainModel converts a model.BlockProgram domain model to a BlockProgramDB database model.
func (db *BlockProgramDB) FromDomainModel(domain model.BlockProgram) {
	db.UUID = uuid.MustParse(domain.ID)
	db.BlockID = uuid.MustParse(domain.BlockID)
	db.ProgramID = uuid.MustParse(domain.ProgramID)
	db.Position = domain.Position
}
