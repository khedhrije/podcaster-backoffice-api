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
func NewProgramCategoryAdapter(client *client) port.ProgramCategoryPersister {
	return &programCategoryAdapter{
		client: client,
	}
}

// Create inserts a new programCategory record into the database.
func (adapter *programCategoryAdapter) Create(ctx context.Context, programCategory model.ProgramCategory) error {
	// SQL query to insert a new programCategory record
	const query = `
        INSERT INTO program_category (UUID, programUUID, categoryUUID)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:programUUID), UUID_TO_BIN(:categoryUUID))
    `
	// Convert domain model to database model
	var programCategoryDB ProgramCategoryDB
	programCategoryDB.FromDomainModel(programCategory)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, programCategoryDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a programCategory record from the database based on its UUID.
func (adapter *programCategoryAdapter) Delete(ctx context.Context, programCategoryUUID string) error {
	// SQL query to delete a programCategory record by UUID
	const query = `
        DELETE FROM program_category WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, programCategoryUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing programCategory record in the database.
func (adapter *programCategoryAdapter) Update(ctx context.Context, programCategoryUUID string, updates model.ProgramCategory) error {
	// SQL query to update a programCategory record
	const query = `
        UPDATE program_category SET 
                             programUUID = UUID_TO_BIN(:programUUID), 
                             categoryUUID = UUID_TO_BIN(:categoryUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = programCategoryUUID
	// Convert domain model to database model
	var programCategoryDB ProgramCategoryDB
	programCategoryDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, programCategoryDB)
	if err != nil {
		return err
	}
	return nil
}

// Find retrieves a programCategory record from the database by its UUID.
func (adapter *programCategoryAdapter) Find(ctx context.Context, programCategoryUUID string) (*model.ProgramCategory, error) {
	// SQL query to select a programCategory record by UUID
	const query = `
        SELECT * FROM program_category WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve result
	var programCategoryDB ProgramCategoryDB
	if err := adapter.client.db.GetContext(ctx, &programCategoryDB, query, programCategoryUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if programCategoryDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("programCategory with ID %s not found", programCategoryUUID)
	}
	// Convert database model to domain model
	result := programCategoryDB.ToDomainModel()
	return &result, nil
}

// FindByProgramID retrieves all programCategory records from the database for a given program ID.
func (adapter *programCategoryAdapter) FindByProgramID(ctx context.Context, programID string) ([]*model.ProgramCategory, error) {
	// SQL query to select programCategory records by program ID
	const query = `
        SELECT * FROM program_category WHERE programUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programCategoriesDB []*ProgramCategoryDB
	if err := adapter.client.db.SelectContext(ctx, &programCategoriesDB, query, programID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var programCategories []*model.ProgramCategory
	for _, programCategoryDB := range programCategoriesDB {
		mappedProgramCategory := programCategoryDB.ToDomainModel()
		programCategories = append(programCategories, &mappedProgramCategory)
	}
	return programCategories, nil
}

// FindByCategoryIDAndProgramID retrieves all programCategory records from the database
// for a given category ID and program ID.
func (adapter *programCategoryAdapter) FindByCategoryIDAndProgramID(ctx context.Context, categoryID, programID string) ([]*model.ProgramCategory, error) {
	// SQL query to select programCategory records by category ID and program ID
	const query = `
        SELECT * FROM program_category WHERE categoryUUID = UUID_TO_BIN(?) AND programUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programCategoriesDB []*ProgramCategoryDB
	if err := adapter.client.db.SelectContext(ctx, &programCategoriesDB, query, categoryID, programID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var programCategories []*model.ProgramCategory
	for _, programCategoryDB := range programCategoriesDB {
		mappedProgramCategory := programCategoryDB.ToDomainModel()
		programCategories = append(programCategories, &mappedProgramCategory)
	}
	return programCategories, nil
}

// FindByCategoryID retrieves all programCategory records from the database for a given category ID.
func (adapter *programCategoryAdapter) FindByCategoryID(ctx context.Context, categoryID string) ([]*model.ProgramCategory, error) {
	// SQL query to select programCategory records by category ID
	const query = `
        SELECT * FROM program_category WHERE categoryUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programCategoriesDB []*ProgramCategoryDB
	if err := adapter.client.db.SelectContext(ctx, &programCategoriesDB, query, categoryID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
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
func (db *ProgramCategoryDB) ToDomainModel() model.ProgramCategory {
	return model.ProgramCategory{
		ID:         db.UUID.String(),
		ProgramID:  db.ProgramID.String(),
		CategoryID: db.CategoryID.String(),
	}
}

// FromDomainModel converts a model.ProgramCategory domain model to a ProgramCategoryDB database model.
func (db *ProgramCategoryDB) FromDomainModel(domain model.ProgramCategory) {
	db.UUID = uuid.MustParse(domain.ID)
	db.ProgramID = uuid.MustParse(domain.ProgramID)
	db.CategoryID = uuid.MustParse(domain.CategoryID)
}
