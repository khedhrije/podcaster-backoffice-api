// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an instance of tagAdapter.
func NewTagAdapter(client *client) port.TagPersister {
	return &tagAdapter{
		client: client,
	}
}

// Create inserts a new tag record into the database.
// It takes a context and a model.Tag, and returns an error if the operation fails.
func (adapter *tagAdapter) Create(ctx context.Context, tag model.Tag) error {
	const query = `
        INSERT INTO tag (UUID, name, description)
        VALUES (UUID_TO_BIN(:UUID), :name, :description)
    `
	var tagDB TagDB
	tagDB.FromDomainModel(tag)
	_, err := adapter.client.db.NamedExecContext(ctx, query, tagDB)
	return err
}

// Delete removes a tag record from the database based on its UUID.
// It takes a context and the tag's UUID, and returns an error if the operation fails.
func (adapter *tagAdapter) Delete(ctx context.Context, tagUUID string) error {
	const query = `
        DELETE FROM tag WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, tagUUID)
	return err
}

// Update updates an existing tag record in the database.
// It takes a context, the tag's UUID, and the updated model.Tag, and returns an error if the operation fails.
func (adapter *tagAdapter) Update(ctx context.Context, tagUUID string, updates model.Tag) error {
	const query = `
        UPDATE tag SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description) 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = tagUUID
	var tagDB TagDB
	tagDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, tagDB)
	return err
}

// FindAll retrieves all tag records from the database.
// It takes a context and returns a slice of model.Tag and an error if the operation fails.
func (adapter *tagAdapter) FindAll(ctx context.Context) ([]*model.Tag, error) {
	const query = `
        SELECT * FROM tag
    `
	var tagsDB []*TagDB
	if err := adapter.client.db.SelectContext(ctx, &tagsDB, query); err != nil {
		return nil, err
	}
	var tags []*model.Tag
	for _, tagDB := range tagsDB {
		mappedTag := tagDB.ToDomainModel()
		tags = append(tags, &mappedTag)
	}
	return tags, nil
}

// Find retrieves a tag record from the database by its UUID.
// It takes a context and the tag's UUID, and returns a model.Tag and an error if the operation fails.
func (adapter *tagAdapter) Find(ctx context.Context, tagUUID string) (*model.Tag, error) {
	const query = `
        SELECT * FROM tag WHERE UUID = UUID_TO_BIN(?)
    `
	var tagsDB []*TagDB
	if err := adapter.client.db.SelectContext(ctx, &tagsDB, query, tagUUID); err != nil {
		return nil, err
	}
	if len(tagsDB) == 0 {
		return nil, nil
	}
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
// It returns the corresponding model.Tag.
func (db *TagDB) ToDomainModel() model.Tag {
	return model.Tag{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Tag domain model to a TagDB database model.
// It sets the fields of the TagDB based on the given model.Tag.
func (db *TagDB) FromDomainModel(domain model.Tag) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
