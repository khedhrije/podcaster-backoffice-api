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
func NewWallAdapter(client *client) port.WallPersister {
	return &wallAdapter{
		client: client,
	}
}

// Create inserts a new wall record into the database.
func (adapter wallAdapter) Create(ctx context.Context, wall model.Wall) error {
	// SQL query to insert a new wall record
	const query = `
        INSERT INTO wall (UUID, name, description)
        VALUES (UUID_TO_BIN(:UUID), :name, :description)
    `
	// Convert domain model to database model
	var wallDB WallDB
	wallDB.FromDomainModel(wall)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a wall record from the database based on its UUID.
func (adapter wallAdapter) Delete(ctx context.Context, wallUUID string) error {
	// SQL query to delete a wall record by UUID
	const query = `
        DELETE FROM wall WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, wallUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing wall record in the database.
func (adapter wallAdapter) Update(ctx context.Context, wallUUID string, updates model.Wall) error {
	// SQL query to update a wall record
	const query = `
        UPDATE wall SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description), 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = wallUUID
	// Convert domain model to database model
	var wallDB WallDB
	wallDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, wallDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all wall records from the database.
func (adapter wallAdapter) FindAll(ctx context.Context) ([]*model.Wall, error) {
	// SQL query to select all wall records
	const query = `
        SELECT * FROM wall
    `
	// Execute query and retrieve results
	var wallsDB []*WallDB
	if err := adapter.client.db.SelectContext(ctx, &wallsDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var walls []*model.Wall
	for _, wallDB := range wallsDB {
		mappedWall := wallDB.ToDomainModel()
		walls = append(walls, &mappedWall)
	}
	return walls, nil
}

// FindByUUID retrieves a wall record from the database by its UUID.
func (adapter wallAdapter) Find(ctx context.Context, wallUUID string) (*model.Wall, error) {
	// SQL query to select a wall record by UUID
	const query = `
        SELECT * FROM wall WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var wallsDB []*WallDB
	if err := adapter.client.db.SelectContext(ctx, &wallsDB, query, wallUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if len(wallsDB) == 0 {
		return nil, nil
	}
	// Convert database model to domain model
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
func (db *WallDB) ToDomainModel() model.Wall {
	return model.Wall{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Wall domain model to a WallDB database model.
func (db *WallDB) FromDomainModel(domain model.Wall) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
