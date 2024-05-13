// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an implementation of the BlockPersister interface.
func NewBlockAdapter(client *client) port.BlockPersister {
	return &blockAdapter{
		client: client,
	}
}

// Create inserts a new block record into the database.
// It takes a context and a model.Block, and returns an error if the operation fails.
func (adapter *blockAdapter) Create(ctx context.Context, block model.Block) error {
	const query = `
        INSERT INTO block (UUID, name, description, kind)
        VALUES (UUID_TO_BIN(:UUID), :name, :description, :kind)
    `
	var blockDB BlockDB
	blockDB.FromDomainModel(block)
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockDB)
	return err
}

// Delete removes a block record from the database based on its UUID.
// It takes a context and the block's UUID, and returns an error if the operation fails.
func (adapter *blockAdapter) Delete(ctx context.Context, blockUUID string) error {
	const query = `
        DELETE FROM block WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, blockUUID)
	return err
}

// Update updates an existing block record in the database.
// It takes a context, the block's UUID, and the updated model.Block, and returns an error if the operation fails.
func (adapter *blockAdapter) Update(ctx context.Context, blockUUID string, updates model.Block) error {
	const query = `
        UPDATE block SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description), 
                             kind = COALESCE(:kind, kind) 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = blockUUID
	var blockDB BlockDB
	blockDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, blockDB)
	return err
}

// FindAll retrieves all block records from the database.
// It takes a context and returns a slice of model.Block and an error if the operation fails.
func (adapter *blockAdapter) FindAll(ctx context.Context) ([]*model.Block, error) {
	const query = `
        SELECT * FROM block;
    `
	var blocksDB []*BlockDB
	if err := adapter.client.db.SelectContext(ctx, &blocksDB, query); err != nil {
		return nil, err
	}
	var blocks []*model.Block
	for _, blockDB := range blocksDB {
		mappedBlock := blockDB.ToDomainModel()
		blocks = append(blocks, &mappedBlock)
	}
	return blocks, nil
}

// Find retrieves a block record from the database by its UUID.
// It takes a context and the block's UUID, and returns a model.Block and an error if the operation fails.
func (adapter *blockAdapter) Find(ctx context.Context, blockUUID string) (*model.Block, error) {
	const query = `
        SELECT * FROM block WHERE UUID = UUID_TO_BIN(?);
    `
	var blocksDB []*BlockDB
	if err := adapter.client.db.SelectContext(ctx, &blocksDB, query, blockUUID); err != nil {
		return nil, err
	}
	if len(blocksDB) == 0 {
		return nil, nil
	}
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
// It returns the corresponding model.Block.
func (db *BlockDB) ToDomainModel() model.Block {
	return model.Block{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
		Kind:        db.Kind.String,
	}
}

// FromDomainModel converts a model.Block domain model to a BlockDB database model.
// It sets the fields of the BlockDB based on the given model.Block.
func (db *BlockDB) FromDomainModel(domain model.Block) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
	db.Kind = sql.NullString{String: domain.Kind, Valid: domain.Kind != ""}
}
