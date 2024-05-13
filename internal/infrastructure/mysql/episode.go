package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// episodeAdapter is a struct that acts as an adapter for interacting with
// the episode data in the MySQL database.
type episodeAdapter struct {
	client *client
}

// NewEpisodeAdapter creates a new episode adapter with the provided MySQL client.
func NewEpisodeAdapter(client *client) port.EpisodePersister {
	return &episodeAdapter{
		client: client,
	}
}

func (adapter *episodeAdapter) FindByProgramID(ctx context.Context, id string) ([]*model.Episode, error) {
	// SQL query to select all episode records
	const query = `
        SELECT * FROM episode WHERE programUUID = UUID_TO_BIN(?);
    `
	// Execute query and retrieve results
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query, id); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var episodes []*model.Episode
	for _, episodeDB := range episodesDB {
		mappedEpisode := episodeDB.ToDomainModel()
		episodes = append(episodes, &mappedEpisode)
	}
	return episodes, nil
}

// Create inserts a new episode record into the database.
func (adapter *episodeAdapter) Create(ctx context.Context, episode model.Episode) error {
	// SQL query to insert a new episode record
	const query = `
        INSERT INTO episode (UUID, name, description, position, programUUID)
        VALUES (UUID_TO_BIN(:UUID), :name, :description, :position, UUID_TO_BIN(:programUUID));
    `
	// Convert domain model to database model
	var episodeDB EpisodeDB
	episodeDB.FromDomainModel(episode)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, episodeDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes an episode record from the database based on its UUID.
func (adapter *episodeAdapter) Delete(ctx context.Context, episodeUUID string) error {
	// SQL query to delete an episode record by UUID
	const query = `
        DELETE FROM episode WHERE UUID = UUID_TO_BIN(?);
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, episodeUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing episode record in the database.
func (adapter *episodeAdapter) Update(ctx context.Context, episodeUUID string, updates model.Episode) error {
	// SQL query to update an episode record
	const query = `
        UPDATE episode SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description), 
                             position = COALESCE(NULLIF(:position, 0), position),
                             programUUID = COALESCE(NULLIF(UUID_TO_BIN(:programUUID), UUID_TO_BIN('00000000-0000-0000-0000-000000000000')), programUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID);
    `
	// Set UUID for updates
	updates.ID = episodeUUID
	// Convert domain model to database model
	var episodeDB EpisodeDB
	episodeDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, episodeDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all episode records from the database.
func (adapter *episodeAdapter) FindAll(ctx context.Context) ([]*model.Episode, error) {
	// SQL query to select all episode records
	const query = `
        SELECT * FROM episode;
    `
	// Execute query and retrieve results
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var episodes []*model.Episode
	for _, episodeDB := range episodesDB {
		mappedEpisode := episodeDB.ToDomainModel()
		episodes = append(episodes, &mappedEpisode)
	}
	return episodes, nil
}

// FindByUUID retrieves an episode record from the database by its UUID.
func (adapter *episodeAdapter) Find(ctx context.Context, episodeUUID string) (*model.Episode, error) {
	// SQL query to select an episode record by UUID
	const query = `
        SELECT * FROM episode WHERE UUID = UUID_TO_BIN(?)
    `
	// Execute query and retrieve results
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query, episodeUUID); err != nil {
		return nil, err
	}
	// Check if the record exists
	if len(episodesDB) == 0 {
		return nil, nil
	}
	// Convert database model to domain model
	result := episodesDB[0].ToDomainModel()
	return &result, nil
}

// EpisodeDB is a struct representing the episode database model.
type EpisodeDB struct {
	UUID        uuid.UUID      `db:"UUID"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Position    int            `db:"position"`
	ProgramID   uuid.UUID      `db:"programUUID"`
	CreatedAt   sql.NullTime   `db:"createdAt"`
	UpdatedAt   sql.NullTime   `db:"updatedAt"`
}

// ToDomainModel converts an EpisodeDB database model to a model.Episode domain model.
func (db *EpisodeDB) ToDomainModel() model.Episode {
	return model.Episode{
		ID:          db.UUID.String(),
		Name:        db.Name.String,
		Description: db.Description.String,
		Position:    db.Position,
		ProgramID:   db.ProgramID.String(),
	}
}

// FromDomainModel converts a model.Episode domain model to an EpisodeDB database model.
func (db *EpisodeDB) FromDomainModel(domain model.Episode) {
	db.UUID = uuid.MustParse(domain.ID)
	db.Name = sql.NullString{String: domain.Name, Valid: domain.Name != ""}
	db.Description = sql.NullString{String: domain.Description, Valid: domain.Description != ""}
	db.Position = domain.Position
	db.ProgramID = uuid.Nil
	if domain.ProgramID != "" {
		db.ProgramID = uuid.MustParse(domain.ProgramID)
	}

}
