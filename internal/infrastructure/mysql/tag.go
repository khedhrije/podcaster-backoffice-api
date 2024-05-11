package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// tagAdapter is a struct that acts as an adapter for interacting with
// the tag data in the MySQL database.
type tagAdapter struct {
	client *client
}

// NewTagAdapter creates a new tag adapter with the provided MySQL client.
func NewTagAdapter(client *client) port.TagPersister {
	return &tagAdapter{
		client: client,
	}
}

// Create inserts a new tag record into the database.
func (adapter *tagAdapter) Create(ctx context.Context, tag model.Tag) error {
	// SQL query to insert a new tag record
	const query = `
        INSERT INTO tag (UUID, name, description)
        VALUES (UUID_TO_BIN(:UUID), :name, :description)
    `
	// Convert domain model to database model
	var tagDB TagDB
	tagDB.FromDomainModel(tag)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, tagDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a tag record from the database based on its UUID.
func (adapter *tagAdapter) Delete(ctx context.Context, tagUUID string) error {
	// SQL query to delete a tag record by UUID
	const query = `
        DELETE FROM tag WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, tagUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing tag record in the database.
func (adapter *tagAdapter) Update(ctx context.Context, tagUUID string, updates model.Tag) error {
	// SQL query to update a tag record
	const query = `
        UPDATE tag SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description) 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = tagUUID
	// Convert domain model to database model
	var tagDB TagDB
	tagDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, tagDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all tag records from the database.
func (adapter *tagAdapter) FindAll(ctx context.Context) ([]*model.Tag, error) {
	// SQL query to select all tag records
	const query = `
        SELECT * FROM tag
    `
	// Execute query and retrieve results
	var tagsDB []*TagDB
	if err := adapter.client.db.SelectContext(ctx, &tagsDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var tags []*model.Tag
	for _, tagDB := range tagsDB {
		mappedTag := tagDB.ToDomainModel()
		tags = append(tags, &mappedTag)
	}
	return tags, nil
}

// FindByUUID retrieves a tag record from the database by its UUID.
func (adapter *tagAdapter) Find(ctx context.Context, tagUUID string) (*model.Tag, error) {
	// SQL query to select a tag record by UUID
	const query = `
        SELECT * FROM tag WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var tagsDB []*TagDB
	if err := adapter.client.db.SelectContext(ctx, &tagsDB, query, tagUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if len(tagsDB) == 0 {
		return nil, nil
	}
	// Convert database model to domain model
	result := tagsDB[0].ToDomainModel()
	return &result, nil
}

// TagDB is a struct representing the tag database model.
type TagDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a TagDB database model to a model.Tag domain model.
func (db *TagDB) ToDomainModel() model.Tag {
	return model.Tag{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Tag domain model to a TagDB database model.
func (db *TagDB) FromDomainModel(domain model.Tag) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
