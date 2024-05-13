// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an implementation of the BlockProgramPersister interface.
func NewBlockProgramAdapter(client *client) port.BlockProgramPersister {
	return &blockProgramAdapter{
		client: client,
	}
}

// Create inserts a new blockProgram record into the database.
// It takes a context and a model.BlockProgram, and returns an error if the operation fails.
func (adapter *blockProgramAdapter) Create(ctx context.Context, blockProgram model.BlockProgram) error {
	const query = `
        INSERT INTO block_program (UUID, blockUUID, programUUID, position)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:blockUUID), UUID_TO_BIN(:programUUID), :position)
    `
	var blockProgramDB BlockProgramDB
	blockProgramDB.FromDomainModel(blockProgram)
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockProgramDB)
	return err
}

// Delete removes a blockProgram record from the database based on its UUID.
// It takes a context and the blockProgram's UUID, and returns an error if the operation fails.
func (adapter *blockProgramAdapter) Delete(ctx context.Context, blockProgramUUID string) error {
	const query = `
        DELETE FROM block_program WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, blockProgramUUID)
	return err
}

// Update updates an existing blockProgram record in the database.
// It takes a context, the blockProgram's UUID, and the updated model.BlockProgram, and returns an error if the operation fails.
func (adapter *blockProgramAdapter) Update(ctx context.Context, blockProgramUUID string, updates model.BlockProgram) error {
	const query = `
        UPDATE block_program SET 
                             blockUUID = UUID_TO_BIN(:blockUUID), 
                             programUUID = UUID_TO_BIN(:programUUID), 
                             position = :position
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = blockProgramUUID
	var blockProgramDB BlockProgramDB
	blockProgramDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockProgramDB)
	return err
}

// Find retrieves a blockProgram record from the database by its UUID.
// It takes a context and the blockProgram's UUID, and returns a model.BlockProgram and an error if the operation fails.
func (adapter *blockProgramAdapter) Find(ctx context.Context, blockProgramUUID string) (*model.BlockProgram, error) {
	const query = `
        SELECT * FROM block_program WHERE UUID = UUID_TO_BIN(?)
    `
	var blockProgramDB BlockProgramDB
	if err := adapter.client.db.GetContext(ctx, &blockProgramDB, query, blockProgramUUID); err != nil {
		return nil, err
	}
	if blockProgramDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("blockProgram with ID %s not found", blockProgramUUID)
	}
	result := blockProgramDB.ToDomainModel()
	return &result, nil
}

// FindByBlockID retrieves all blockProgram records from the database for a given block ID.
// It takes a context and the block's ID, and returns a slice of model.BlockProgram and an error if the operation fails.
func (adapter *blockProgramAdapter) FindByBlockID(ctx context.Context, blockID string) ([]*model.BlockProgram, error) {
	const query = `
        SELECT * FROM block_program WHERE blockUUID = UUID_TO_BIN(?)
    `
	var blockProgramsDB []*BlockProgramDB
	if err := adapter.client.db.SelectContext(ctx, &blockProgramsDB, query, blockID); err != nil {
		return nil, err
	}
	var blockPrograms []*model.BlockProgram
	for _, blockProgramDB := range blockProgramsDB {
		mappedBlockProgram := blockProgramDB.ToDomainModel()
		blockPrograms = append(blockPrograms, &mappedBlockProgram)
	}
	return blockPrograms, nil
}

// FindByProgramID retrieves all blockProgram records from the database for a given program ID.
// It takes a context and the program's ID, and returns a slice of model.BlockProgram and an error if the operation fails.
func (adapter *blockProgramAdapter) FindByProgramID(ctx context.Context, programID string) ([]*model.BlockProgram, error) {
	const query = `
        SELECT * FROM block_program WHERE programUUID = UUID_TO_BIN(?)
    `
	var blockProgramsDB []*BlockProgramDB
	if err := adapter.client.db.SelectContext(ctx, &blockProgramsDB, query, programID); err != nil {
		return nil, err
	}
	var blockPrograms []*model.BlockProgram
	for _, blockProgramDB := range blockProgramsDB {
		mappedBlockProgram := blockProgramDB.ToDomainModel()
		blockPrograms = append(blockPrograms, &mappedBlockProgram)
	}
	return blockPrograms, nil
}

// FindByBlockIDAndProgramID retrieves all blockProgram records from the database for a given block ID and program ID.
// It takes a context, the block's ID, and the program's ID, and returns a slice of model.BlockProgram and an error if the operation fails.
func (adapter *blockProgramAdapter) FindByBlockIDAndProgramID(ctx context.Context, blockID, programID string) ([]*model.BlockProgram, error) {
	const query = `
        SELECT * FROM block_program WHERE blockUUID = UUID_TO_BIN(?) AND programUUID = UUID_TO_BIN(?)
    `
	var blockProgramsDB []*BlockProgramDB
	if err := adapter.client.db.SelectContext(ctx, &blockProgramsDB, query, blockID, programID); err != nil {
		return nil, err
	}
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
// It returns the corresponding model.BlockProgram.
func (db *BlockProgramDB) ToDomainModel() model.BlockProgram {
	return model.BlockProgram{
		ID:        db.UUID.String(),
		BlockID:   db.BlockID.String(),
		ProgramID: db.ProgramID.String(),
		Position:  db.Position,
	}
}

// FromDomainModel converts a model.BlockProgram domain model to a BlockProgramDB database model.
// It sets the fields of the BlockProgramDB based on the given model.BlockProgram.
func (db *BlockProgramDB) FromDomainModel(domain model.BlockProgram) {
	db.UUID = uuid.MustParse(domain.ID)
	db.BlockID = uuid.MustParse(domain.BlockID)
	db.ProgramID = uuid.MustParse(domain.ProgramID)
	db.Position = domain.Position
}
