package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// categoryAdapter is a struct that acts as an adapter for interacting with
// the category data in the MySQL database.
type categoryAdapter struct {
	client *client
}

// NewCategoryAdapter creates a new category adapter with the provided MySQL client.
func NewCategoryAdapter(client *client) port.CategoryPersister {
	return &categoryAdapter{
		client: client,
	}
}

// Create inserts a new category record into the database.
func (adapter *categoryAdapter) Create(ctx context.Context, category model.Category) error {
	// SQL query to insert a new category record
	const query = `
        INSERT INTO category (UUID, name, description, parentUUID)
        VALUES (UUID_TO_BIN(:UUID), :name, :description, UUID_TO_BIN(:parentUUID))
    `
	// Convert domain model to database model
	var categoryDB CategoryDB
	categoryDB.FromDomainModel(category)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, categoryDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a category record from the database based on its UUID.
func (adapter *categoryAdapter) Delete(ctx context.Context, categoryUUID string) error {
	// SQL query to delete a category record by UUID
	const query = `
        DELETE FROM category WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, categoryUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing category record in the database.
func (adapter *categoryAdapter) Update(ctx context.Context, categoryUUID string, updates model.Category) error {
	// SQL query to update a category record
	const query = `
        UPDATE category SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description),
                             parentUUID = UUID_TO_BIN(:parentUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = categoryUUID
	// Convert domain model to database model
	var categoryDB CategoryDB
	categoryDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, categoryDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all category records from the database.
func (adapter *categoryAdapter) FindAll(ctx context.Context) ([]*model.Category, error) {
	// SQL query to select all category records
	const query = `
        SELECT * FROM category
    `
	// Execute query and retrieve results
	var categoriesDB []*CategoryDB
	if err := adapter.client.db.SelectContext(ctx, &categoriesDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var categories []*model.Category
	for _, categoryDB := range categoriesDB {
		mappedCategory := categoryDB.ToDomainModel()
		categories = append(categories, &mappedCategory)
	}
	return categories, nil
}

// FindByUUID retrieves a category record from the database by its UUID.
func (adapter *categoryAdapter) Find(ctx context.Context, categoryUUID string) (*model.Category, error) {
	// SQL query to select a category record by UUID
	const query = `
        SELECT * FROM category WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var categoryDB CategoryDB
	if err := adapter.client.db.GetContext(ctx, &categoryDB, query, categoryUUID); err != nil {
		return nil, err
	}
	// Convert database model to domain model
	result := categoryDB.ToDomainModel()
	return &result, nil
}

// CategoryDB is a struct representing the category database model.
type CategoryDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	ParentID    sql.NullString `db:"parentUUID"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a CategoryDB database model to a model.Category domain model.
func (db *CategoryDB) ToDomainModel() model.Category {
	return model.Category{
		ID:          db.UUID.String(),
		Name:        db.Name,
		Description: db.Description,
		Parent: &model.Category{
			ID: db.ParentID.String,
		},
	}
}

// FromDomainModel converts a model.Category domain model to a CategoryDB database model.
func (db *CategoryDB) FromDomainModel(domain model.Category) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = domain.Name
	db.Description = domain.Description
	if domain.Parent.ID != "" {
		db.ParentID = sql.NullString{String: domain.Parent.ID, Valid: true}
	} else {
		db.ParentID = sql.NullString{Valid: false}
	}
}
