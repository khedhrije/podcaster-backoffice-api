package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// blockAdapter is a struct that acts as an adapter for interacting with
// the block data in the MySQL database.
type blockAdapter struct {
	client *client
}

// NewBlockAdapter creates a new block adapter with the provided MySQL client.
func NewBlockAdapter(client *client) port.BlockPersister {
	return &blockAdapter{
		client: client,
	}
}

// Create inserts a new block record into the database.
func (adapter *blockAdapter) Create(ctx context.Context, block model.Block) error {
	// SQL query to insert a new block record
	const query = `
        INSERT INTO block (UUID, name, description, kind)
        VALUES (UUID_TO_BIN(:UUID), :name, :description, :kind)
    `
	// Convert domain model to database model
	var blockDB BlockDB
	blockDB.FromDomainModel(block)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a block record from the database based on its UUID.
func (adapter *blockAdapter) Delete(ctx context.Context, blockUUID string) error {
	// SQL query to delete a block record by UUID
	const query = `
        DELETE FROM block WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, blockUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing block record in the database.
func (adapter *blockAdapter) Update(ctx context.Context, blockUUID string, updates model.Block) error {
	// SQL query to update a block record
	const query = `
        UPDATE block SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description), 
                             kind = COALESCE(:kind, kind) 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = blockUUID
	// Convert domain model to database model
	var blockDB BlockDB
	blockDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all block records from the database.
func (adapter *blockAdapter) FindAll(ctx context.Context) ([]*model.Block, error) {
	// SQL query to select all block records
	const query = `
        SELECT * FROM block
    `
	// Execute query and retrieve results
	var blocksDB []*BlockDB
	if err := adapter.client.db.SelectContext(ctx, &blocksDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var blocks []*model.Block
	for _, blockDB := range blocksDB {
		mappedBlock := blockDB.ToDomainModel()
		blocks = append(blocks, &mappedBlock)
	}
	return blocks, nil
}

// Find retrieves a block record from the database by its UUID.
func (adapter *blockAdapter) Find(ctx context.Context, blockUUID string) (*model.Block, error) {
	// SQL query to select a block record by UUID
	const query = `
        SELECT * FROM block WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var blocksDB []*BlockDB
	if err := adapter.client.db.SelectContext(ctx, &blocksDB, query, blockUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if len(blocksDB) == 0 {
		return nil, nil
	}
	// Convert database model to domain model
	result := blocksDB[0].ToDomainModel()
	return &result, nil
}

// BlockDB is a struct representing the block database model.
type BlockDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Kind        sql.NullString `db:"kind"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a BlockDB database model to a model.Block domain model.
func (db *BlockDB) ToDomainModel() model.Block {
	return model.Block{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
		Kind:        db.Kind.String,
	}
}

// FromDomainModel converts a model.Block domain model to a BlockDB database model.
func (db *BlockDB) FromDomainModel(domain model.Block) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
	db.Kind = sql.NullString{String: domain.Kind, Valid: domain.Kind != ""}
}
