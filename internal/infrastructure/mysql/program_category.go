// Package mysql provides MySQL implementations of the persistence interfaces.
package mysql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// programCategoryAdapter is a struct that acts as an adapter for interacting with
// the program_category data in the MySQL database.
type programCategoryAdapter struct {
	client *client
}

// NewProgramCategoryAdapter creates a new programCategory adapter with the provided MySQL client.
// It returns an implementation of the ProgramCategoryPersister interface.
func NewProgramCategoryAdapter(client *client) port.ProgramCategoryPersister {
	return &programCategoryAdapter{
		client: client,
	}
}

// Create inserts a new programCategory record into the database.
// It takes a context and a model.ProgramCategory, and returns an error if the operation fails.
func (adapter *programCategoryAdapter) Create(ctx context.Context, programCategory model.ProgramCategory) error {
	const query = `
        INSERT INTO program_category (UUID, programUUID, categoryUUID)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:programUUID), UUID_TO_BIN(:categoryUUID))
    `
	var programCategoryDB ProgramCategoryDB
	programCategoryDB.FromDomainModel(programCategory)
	_, err := adapter.client.db.NamedExecContext(ctx, query, programCategoryDB)
	return err
}

// Delete removes a programCategory record from the database based on its UUID.
// It takes a context and the programCategory's UUID, and returns an error if the operation fails.
func (adapter *programCategoryAdapter) Delete(ctx context.Context, programCategoryUUID string) error {
	const query = `
        DELETE FROM program_category WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, programCategoryUUID)
	return err
}

// Update updates an existing programCategory record in the database.
// It takes a context, the programCategory's UUID, and the updated model.ProgramCategory, and returns an error if the operation fails.
func (adapter *programCategoryAdapter) Update(ctx context.Context, programCategoryUUID string, updates model.ProgramCategory) error {
	const query = `
        UPDATE program_category SET 
                             programUUID = UUID_TO_BIN(:programUUID), 
                             categoryUUID = UUID_TO_BIN(:categoryUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = programCategoryUUID
	var programCategoryDB ProgramCategoryDB
	programCategoryDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, programCategoryDB)
	return err
}

// Find retrieves a programCategory record from the database by its UUID.
// It takes a context and the programCategory's UUID, and returns a model.ProgramCategory and an error if the operation fails.
func (adapter *programCategoryAdapter) Find(ctx context.Context, programCategoryUUID string) (*model.ProgramCategory, error) {
	const query = `
        SELECT * FROM program_category WHERE UUID = UUID_TO_BIN(?)
    `
	var programCategoryDB ProgramCategoryDB
	if err := adapter.client.db.GetContext(ctx, &programCategoryDB, query, programCategoryUUID); err != nil {
		return nil, err
	}
	if programCategoryDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("programCategory with ID %s not found", programCategoryUUID)
	}
	result := programCategoryDB.ToDomainModel()
	return &result, nil
}

// FindByProgramID retrieves all programCategory records from the database for a given program ID.
// It takes a context and the program's UUID, and returns a slice of model.ProgramCategory and an error if the operation fails.
func (adapter *programCategoryAdapter) FindByProgramID(ctx context.Context, programID string) ([]*model.ProgramCategory, error) {
	const query = `
        SELECT * FROM program_category WHERE programUUID = UUID_TO_BIN(?)
    `
	var programCategoriesDB []*ProgramCategoryDB
	if err := adapter.client.db.SelectContext(ctx, &programCategoriesDB, query, programID); err != nil {
		return nil, err
	}
	var programCategories []*model.ProgramCategory
	for _, programCategoryDB := range programCategoriesDB {
		mappedProgramCategory := programCategoryDB.ToDomainModel()
		programCategories = append(programCategories, &mappedProgramCategory)
	}
	return programCategories, nil
}

// FindByCategoryIDAndProgramID retrieves all programCategory records from the database
// for a given category ID and program ID.
// It takes a context, the category's UUID, and the program's UUID, and returns a slice of model.ProgramCategory and an error if the operation fails.
func (adapter *programCategoryAdapter) FindByCategoryIDAndProgramID(ctx context.Context, categoryID, programID string) ([]*model.ProgramCategory, error) {
	const query = `
        SELECT * FROM program_category WHERE categoryUUID = UUID_TO_BIN(?) AND programUUID = UUID_TO_BIN(?)
    `
	var programCategoriesDB []*ProgramCategoryDB
	if err := adapter.client.db.SelectContext(ctx, &programCategoriesDB, query, categoryID, programID); err != nil {
		return nil, err
	}
	var programCategories []*model.ProgramCategory
	for _, programCategoryDB := range programCategoriesDB {
		mappedProgramCategory := programCategoryDB.ToDomainModel()
		programCategories = append(programCategories, &mappedProgramCategory)
	}
	return programCategories, nil
}

// FindByCategoryID retrieves all programCategory records from the database for a given category ID.
// It takes a context and the category's UUID, and returns a slice of model.ProgramCategory and an error if the operation fails.
func (adapter *programCategoryAdapter) FindByCategoryID(ctx context.Context, categoryID string) ([]*model.ProgramCategory, error) {
	const query = `
        SELECT * FROM program_category WHERE categoryUUID = UUID_TO_BIN(?)
    `
	var programCategoriesDB []*ProgramCategoryDB
	if err := adapter.client.db.SelectContext(ctx, &programCategoriesDB, query, categoryID); err != nil {
		return nil, err
	}
	var programCategories []*model.ProgramCategory
	for _, programCategoryDB := range programCategoriesDB {
		mappedProgramCategory := programCategoryDB.ToDomainModel()
		programCategories = append(programCategories, &mappedProgramCategory)
	}
	return programCategories, nil
}

// ProgramCategoryDB is a struct representing the programCategory database model.
type ProgramCategoryDB struct {
	UUID       uuid.UUID `db:"UUID"`
	ProgramID  uuid.UUID `db:"programUUID"`
	CategoryID uuid.UUID `db:"categoryUUID"`
}

// ToDomainModel converts a ProgramCategoryDB database model to a model.ProgramCategory domain model.
// It returns the corresponding model.ProgramCategory.
func (db *ProgramCategoryDB) ToDomainModel() model.ProgramCategory {
	return model.ProgramCategory{
		ID:         db.UUID.String(),
		ProgramID:  db.ProgramID.String(),
		CategoryID: db.CategoryID.String(),
	}
}

// FromDomainModel converts a model.ProgramCategory domain model to a ProgramCategoryDB database model.
// It sets the fields of the ProgramCategoryDB based on the given model.ProgramCategory.
func (db *ProgramCategoryDB) FromDomainModel(domain model.ProgramCategory) {
	db.UUID = uuid.MustParse(domain.ID)
	db.ProgramID = uuid.MustParse(domain.ProgramID)
	db.CategoryID = uuid.MustParse(domain.CategoryID)
}
