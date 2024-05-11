package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// programAdapter is a struct that acts as an adapter for interacting with
// the program data in the MySQL database.
type programAdapter struct {
	client *client
}

// NewProgramAdapter creates a new program adapter with the provided MySQL client.
func NewProgramAdapter(client *client) port.ProgramPersister {
	return &programAdapter{
		client: client,
	}
}

// Create inserts a new program record into the database.
func (adapter *programAdapter) Create(ctx context.Context, program model.Program) error {
	// SQL query to insert a new program record
	const query = `
        INSERT INTO program (UUID, name, description)
        VALUES (UUID_TO_BIN(:UUID), :name, :description)
    `
	// Convert domain model to database model
	var programDB ProgramDB
	programDB.FromDomainModel(program)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, programDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a program record from the database based on its UUID.
func (adapter *programAdapter) Delete(ctx context.Context, programUUID string) error {
	// SQL query to delete a program record by UUID
	const query = `
        DELETE FROM program WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, programUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing program record in the database.
func (adapter *programAdapter) Update(ctx context.Context, programUUID string, updates model.Program) error {
	// SQL query to update a program record
	const query = `
        UPDATE program SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description) 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = programUUID
	// Convert domain model to database model
	var programDB ProgramDB
	programDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, programDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all program records from the database.
func (adapter *programAdapter) FindAll(ctx context.Context) ([]*model.Program, error) {
	// SQL query to select all program records
	const query = `
        SELECT * FROM program
    `
	// Execute query and retrieve results
	var programsDB []*ProgramDB
	if err := adapter.client.db.SelectContext(ctx, &programsDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var programs []*model.Program
	for _, programDB := range programsDB {
		mappedProgram := programDB.ToDomainModel()
		programs = append(programs, &mappedProgram)
	}
	return programs, nil
}

// FindByUUID retrieves a program record from the database by its UUID.
func (adapter *programAdapter) Find(ctx context.Context, programUUID string) (*model.Program, error) {
	// SQL query to select a program record by UUID
	const query = `
        SELECT * FROM program WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programsDB []*ProgramDB
	if err := adapter.client.db.SelectContext(ctx, &programsDB, query, programUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if len(programsDB) == 0 {
		return nil, nil
	}
	// Convert database model to domain model
	result := programsDB[0].ToDomainModel()
	return &result, nil
}

// ProgramDB is a struct representing the program database model.
type ProgramDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts a ProgramDB database model to a model.Program domain model.
func (db *ProgramDB) ToDomainModel() model.Program {
	return model.Program{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Program domain model to a ProgramDB database model.
func (db *ProgramDB) FromDomainModel(domain model.Program) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
