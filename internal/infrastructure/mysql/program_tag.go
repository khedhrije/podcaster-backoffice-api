package mysql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
)

// programTagAdapter is a struct that acts as an adapter for interacting with
// the program_tag data in the MySQL database.
type programTagAdapter struct {
	client *client
}

// NewProgramTagAdapter creates a new programTag adapter with the provided MySQL client.
func NewProgramTagAdapter(client *client) *programTagAdapter {
	return &programTagAdapter{
		client: client,
	}
}

// Create inserts a new programTag record into the database.
func (adapter *programTagAdapter) Create(ctx context.Context, programTag model.ProgramTag) error {
	// SQL query to insert a new programTag record
	const query = `
        INSERT INTO program_tag (UUID, programUUID, tagUUID)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:programUUID), UUID_TO_BIN(:tagUUID))
    `
	// Convert domain model to database model
	var programTagDB ProgramTagDB
	programTagDB.FromDomainModel(programTag)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, programTagDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a programTag record from the database based on its UUID.
func (adapter *programTagAdapter) Delete(ctx context.Context, programTagUUID string) error {
	// SQL query to delete a programTag record by UUID
	const query = `
        DELETE FROM program_tag WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, programTagUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing programTag record in the database.
func (adapter *programTagAdapter) Update(ctx context.Context, programTagUUID string, updates model.ProgramTag) error {
	// SQL query to update a programTag record
	const query = `
        UPDATE program_tag SET 
                             programUUID = UUID_TO_BIN(:programUUID), 
                             tagUUID = UUID_TO_BIN(:tagUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	// Set UUID for updates
	updates.ID = programTagUUID
	// Convert domain model to database model
	var programTagDB ProgramTagDB
	programTagDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, programTagDB)
	if err != nil {
		return err
	}
	return nil
}

// Find retrieves a programTag record from the database by its UUID.
func (adapter *programTagAdapter) Find(ctx context.Context, programTagUUID string) (*model.ProgramTag, error) {
	// SQL query to select a programTag record by UUID
	const query = `
        SELECT * FROM program_tag WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve result
	var programTagDB ProgramTagDB
	if err := adapter.client.db.GetContext(ctx, &programTagDB, query, programTagUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if programTagDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("programTag with ID %s not found", programTagUUID)
	}
	// Convert database model to domain model
	result := programTagDB.ToDomainModel()
	return &result, nil
}

// FindByProgramID retrieves all programTag records from the database for a given program ID.
func (adapter *programTagAdapter) FindByProgramID(ctx context.Context, programID string) ([]*model.ProgramTag, error) {
	// SQL query to select programTag records by program ID
	const query = `
        SELECT * FROM program_tag WHERE programUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programTagsDB []*ProgramTagDB
	if err := adapter.client.db.SelectContext(ctx, &programTagsDB, query, programID); err != nil {
		return nil, err
	}

	// Convert database models to domain models
	var programTags []*model.ProgramTag
	for _, programTagDB := range programTagsDB {
		mappedProgramTag := programTagDB.ToDomainModel()
		programTags = append(programTags, &mappedProgramTag)
	}

	return programTags, nil
}

// FindByTagIDAndProgramID retrieves all programTag records from the database
// for a given tag ID and program ID.
func (adapter *programTagAdapter) FindByTagIDAndProgramID(ctx context.Context, tagID, programID string) ([]*model.ProgramTag, error) {
	// SQL query to select programTag records by tag ID and program ID
	const query = `
        SELECT * FROM program_tag WHERE tagUUID = UUID_TO_BIN(?) AND programUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programTagsDB []*ProgramTagDB
	if err := adapter.client.db.SelectContext(ctx, &programTagsDB, query, tagID, programID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var programTags []*model.ProgramTag
	for _, programTagDB := range programTagsDB {
		mappedProgramTag := programTagDB.ToDomainModel()
		programTags = append(programTags, &mappedProgramTag)
	}
	return programTags, nil
}

// FindByTagID retrieves all programTag records from the database for a given tag ID.
func (adapter *programTagAdapter) FindByTagID(ctx context.Context, tagID string) ([]*model.ProgramTag, error) {
	// SQL query to select programTag records by tag ID
	const query = `
        SELECT * FROM program_tag WHERE tagUUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var programTagsDB []*ProgramTagDB
	if err := adapter.client.db.SelectContext(ctx, &programTagsDB, query, tagID); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var programTags []*model.ProgramTag
	for _, programTagDB := range programTagsDB {
		mappedProgramTag := programTagDB.ToDomainModel()
		programTags = append(programTags, &mappedProgramTag)
	}
	return programTags, nil
}

// ProgramTagDB is a struct representing the programTag database model.
type ProgramTagDB struct {
	UUID      uuid.UUID `db:"UUID"`
	ProgramID uuid.UUID `db:"programUUID"`
	TagID     uuid.UUID `db:"tagUUID"`
}

// ToDomainModel converts a ProgramTagDB database model to a model.ProgramTag domain model.
func (db *ProgramTagDB) ToDomainModel() model.ProgramTag {
	return model.ProgramTag{
		ID:        db.UUID.String(),
		ProgramID: db.ProgramID.String(),
		TagID:     db.TagID.String(),
	}
}

// FromDomainModel converts a model.ProgramTag domain model to a ProgramTagDB database model.
func (db *ProgramTagDB) FromDomainModel(domain model.ProgramTag) {
	db.UUID = uuid.MustParse(domain.ID)
	db.ProgramID = uuid.MustParse(domain.ProgramID)
	db.TagID = uuid.MustParse(domain.TagID)
}
