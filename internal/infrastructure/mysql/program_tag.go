// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an instance of programTagAdapter.
func NewProgramTagAdapter(client *client) *programTagAdapter {
	return &programTagAdapter{
		client: client,
	}
}

// Create inserts a new programTag record into the database.
// It takes a context and a model.ProgramTag, and returns an error if the operation fails.
func (adapter *programTagAdapter) Create(ctx context.Context, programTag model.ProgramTag) error {
	const query = `
        INSERT INTO program_tag (UUID, programUUID, tagUUID)
        VALUES (UUID_TO_BIN(:UUID), UUID_TO_BIN(:programUUID), UUID_TO_BIN(:tagUUID))
    `
	var programTagDB ProgramTagDB
	programTagDB.FromDomainModel(programTag)
	_, err := adapter.client.db.NamedExecContext(ctx, query, programTagDB)
	return err
}

// Delete removes a programTag record from the database based on its UUID.
// It takes a context and the programTag's UUID, and returns an error if the operation fails.
func (adapter *programTagAdapter) Delete(ctx context.Context, programTagUUID string) error {
	const query = `
        DELETE FROM program_tag WHERE UUID = UUID_TO_BIN(?)
    `
	_, err := adapter.client.db.ExecContext(ctx, query, programTagUUID)
	return err
}

// Update updates an existing programTag record in the database.
// It takes a context, the programTag's UUID, and the updated model.ProgramTag, and returns an error if the operation fails.
func (adapter *programTagAdapter) Update(ctx context.Context, programTagUUID string, updates model.ProgramTag) error {
	const query = `
        UPDATE program_tag SET 
                             programUUID = UUID_TO_BIN(:programUUID), 
                             tagUUID = UUID_TO_BIN(:tagUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = programTagUUID
	var programTagDB ProgramTagDB
	programTagDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, programTagDB)
	return err
}

// Find retrieves a programTag record from the database by its UUID.
// It takes a context and the programTag's UUID, and returns a model.ProgramTag and an error if the operation fails.
func (adapter *programTagAdapter) Find(ctx context.Context, programTagUUID string) (*model.ProgramTag, error) {
	const query = `
        SELECT * FROM program_tag WHERE UUID = UUID_TO_BIN(?)
    `
	var programTagDB ProgramTagDB
	if err := adapter.client.db.GetContext(ctx, &programTagDB, query, programTagUUID); err != nil {
		return nil, err
	}
	if programTagDB.UUID == uuid.Nil {
		return nil, fmt.Errorf("programTag with ID %s not found", programTagUUID)
	}
	result := programTagDB.ToDomainModel()
	return &result, nil
}

// FindByProgramID retrieves all programTag records from the database for a given program ID.
// It takes a context and the program's UUID, and returns a slice of model.ProgramTag and an error if the operation fails.
func (adapter *programTagAdapter) FindByProgramID(ctx context.Context, programID string) ([]*model.ProgramTag, error) {
	const query = `
        SELECT * FROM program_tag WHERE programUUID = UUID_TO_BIN(?)
    `
	var programTagsDB []*ProgramTagDB
	if err := adapter.client.db.SelectContext(ctx, &programTagsDB, query, programID); err != nil {
		return nil, err
	}
	var programTags []*model.ProgramTag
	for _, programTagDB := range programTagsDB {
		mappedProgramTag := programTagDB.ToDomainModel()
		programTags = append(programTags, &mappedProgramTag)
	}
	return programTags, nil
}

// FindByTagIDAndProgramID retrieves all programTag records from the database
// for a given tag ID and program ID.
// It takes a context, the tag's UUID, and the program's UUID, and returns a slice of model.ProgramTag and an error if the operation fails.
func (adapter *programTagAdapter) FindByTagIDAndProgramID(ctx context.Context, tagID, programID string) ([]*model.ProgramTag, error) {
	const query = `
        SELECT * FROM program_tag WHERE tagUUID = UUID_TO_BIN(?) AND programUUID = UUID_TO_BIN(?)
    `
	var programTagsDB []*ProgramTagDB
	if err := adapter.client.db.SelectContext(ctx, &programTagsDB, query, tagID, programID); err != nil {
		return nil, err
	}
	var programTags []*model.ProgramTag
	for _, programTagDB := range programTagsDB {
		mappedProgramTag := programTagDB.ToDomainModel()
		programTags = append(programTags, &mappedProgramTag)
	}
	return programTags, nil
}

// FindByTagID retrieves all programTag records from the database for a given tag ID.
// It takes a context and the tag's UUID, and returns a slice of model.ProgramTag and an error if the operation fails.
func (adapter *programTagAdapter) FindByTagID(ctx context.Context, tagID string) ([]*model.ProgramTag, error) {
	const query = `
        SELECT * FROM program_tag WHERE tagUUID = UUID_TO_BIN(?)
    `
	var programTagsDB []*ProgramTagDB
	if err := adapter.client.db.SelectContext(ctx, &programTagsDB, query, tagID); err != nil {
		return nil, err
	}
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
// It returns the corresponding model.ProgramTag.
func (db *ProgramTagDB) ToDomainModel() model.ProgramTag {
	return model.ProgramTag{
		ID:        db.UUID.String(),
		ProgramID: db.ProgramID.String(),
		TagID:     db.TagID.String(),
	}
}

// FromDomainModel converts a model.ProgramTag domain model to a ProgramTagDB database model.
// It sets the fields of the ProgramTagDB based on the given model.ProgramTag.
func (db *ProgramTagDB) FromDomainModel(domain model.ProgramTag) {
	db.UUID = uuid.MustParse(domain.ID)
	db.ProgramID = uuid.MustParse(domain.ProgramID)
	db.TagID = uuid.MustParse(domain.TagID)
}
