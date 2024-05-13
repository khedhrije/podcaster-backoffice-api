// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// wallAdapter is a struct that acts as an adapter for interacting with
// the wall data in the MySQL database.
type wallAdapter struct {
	client *client
}

// NewWallAdapter creates a new wall adapter with the provided MySQL client.
// It returns an implementation of the WallPersister interface.
func NewWallAdapter(client *client) port.WallPersister {
	return &wallAdapter{
		client: client,
	}
}

// Create inserts a new wall record into the database.
// It takes a context and a model.Wall, and returns an error if the operation fails.
func (adapter wallAdapter) Create(ctx context.Context, wall model.Wall) error {
	const query = `
        INSERT INTO wall (UUID, name, description)
        VALUES (UUID_TO_BIN(:UUID), :name, :description)
    `
	var wallDB WallDB
	wallDB.FromDomainModel(wall)
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallDB)
	return err
}

// Delete removes a wall record from the database based on its UUID.
// It takes a context and the wall's UUID, and returns an error if the operation fails.
func (adapter wallAdapter) Delete(ctx context.Context, wallUUID string) error {
	const query = `
        DELETE FROM wall WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, wallUUID)
	return err
}

// Update updates an existing wall record in the database.
// It takes a context, the wall's UUID, and the updated model.Wall, and returns an error if the operation fails.
func (adapter wallAdapter) Update(ctx context.Context, wallUUID string, updates model.Wall) error {
	const query = `
        UPDATE wall SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = wallUUID
	var wallDB WallDB
	wallDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallDB)
	return err
}

// FindAll retrieves all wall records from the database.
// It takes a context and returns a slice of model.Wall and an error if the operation fails.
func (adapter wallAdapter) FindAll(ctx context.Context) ([]*model.Wall, error) {
	const query = `
        SELECT * FROM wall
    `
	var wallsDB []*WallDB
	if err := adapter.client.db.SelectContext(ctx, &wallsDB, query); err != nil {
		return nil, err
	}
	var walls []*model.Wall
	for _, wallDB := range wallsDB {
		mappedWall := wallDB.ToDomainModel()
		walls = append(walls, &mappedWall)
	}
	return walls, nil
}

// Find retrieves a wall record from the database by its UUID.
// It takes a context and the wall's UUID, and returns a model.Wall and an error if the operation fails.
func (adapter wallAdapter) Find(ctx context.Context, wallUUID string) (*model.Wall, error) {
	const query = `
        SELECT * FROM wall WHERE UUID = UUID_TO_BIN(?)
    `
	var wallsDB []*WallDB
	if err := adapter.client.db.SelectContext(ctx, &wallsDB, query, wallUUID); err != nil {
		return nil, err
	}
	if len(wallsDB) == 0 {
		return nil, nil
	}
	result := wallsDB[0].ToDomainModel()
	return &result, nil
}

// WallDB is a struct representing the wall database model.
type WallDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a WallDB database model to a model.Wall domain model.
// It returns the corresponding model.Wall.
func (db *WallDB) ToDomainModel() model.Wall {
	return model.Wall{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Wall domain model to a WallDB database model.
// It sets the fields of the WallDB based on the given model.Wall.
func (db *WallDB) FromDomainModel(domain model.Wall) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
