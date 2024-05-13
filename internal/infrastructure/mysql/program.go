// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an implementation of the ProgramPersister interface.
func NewProgramAdapter(client *client) port.ProgramPersister {
	return &programAdapter{
		client: client,
	}
}

// Create inserts a new program record into the database.
// It takes a context and a model.Program, and returns an error if the operation fails.
func (adapter *programAdapter) Create(ctx context.Context, program model.Program) error {
	const query = `
        INSERT INTO program (UUID, name, description)
        VALUES (UUID_TO_BIN(:UUID), :name, :description)
    `
	var programDB ProgramDB
	programDB.FromDomainModel(program)
	_, err := adapter.client.db.NamedExecContext(ctx, query, programDB)
	return err
}

// Delete removes a program record from the database based on its UUID.
// It takes a context and the program's UUID, and returns an error if the operation fails.
func (adapter *programAdapter) Delete(ctx context.Context, programUUID string) error {
	const query = `
        DELETE FROM program WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, programUUID)
	return err
}

// Update updates an existing program record in the database.
// It takes a context, the program's UUID, and the updated model.Program, and returns an error if the operation fails.
func (adapter *programAdapter) Update(ctx context.Context, programUUID string, updates model.Program) error {
	const query = `
        UPDATE program SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description) 
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = programUUID
	var programDB ProgramDB
	programDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, programDB)
	return err
}

// FindAll retrieves all program records from the database.
// It takes a context and returns a slice of model.Program and an error if the operation fails.
func (adapter *programAdapter) FindAll(ctx context.Context) ([]*model.Program, error) {
	const query = `
        SELECT * FROM program
    `
	var programsDB []*ProgramDB
	if err := adapter.client.db.SelectContext(ctx, &programsDB, query); err != nil {
		return nil, err
	}
	var programs []*model.Program
	for _, programDB := range programsDB {
		mappedProgram := programDB.ToDomainModel()
		programs = append(programs, &mappedProgram)
	}
	return programs, nil
}

// Find retrieves a program record from the database by its UUID.
// It takes a context and the program's UUID, and returns a model.Program and an error if the operation fails.
func (adapter *programAdapter) Find(ctx context.Context, programUUID string) (*model.Program, error) {
	const query = `
        SELECT * FROM program WHERE UUID = UUID_TO_BIN(?)
    `
	var programsDB []*ProgramDB
	if err := adapter.client.db.SelectContext(ctx, &programsDB, query, programUUID); err != nil {
		return nil, err
	}
	if len(programsDB) == 0 {
		return nil, nil
	}
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
// It returns the corresponding model.Program.
func (db *ProgramDB) ToDomainModel() model.Program {
	return model.Program{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
	}
}

// FromDomainModel converts a model.Program domain model to a ProgramDB database model.
// It sets the fields of the ProgramDB based on the given model.Program.
func (db *ProgramDB) FromDomainModel(domain model.Program) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
}
