// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an implementation of the EpisodePersister interface.
func NewEpisodeAdapter(client *client) port.EpisodePersister {
	return &episodeAdapter{
		client: client,
	}
}

// FindByProgramID retrieves all episode records from the database for a given program ID.
// It takes a context and the program's ID, and returns a slice of model.Episode and an error if the operation fails.
func (adapter *episodeAdapter) FindByProgramID(ctx context.Context, id string) ([]*model.Episode, error) {
	const query = `
        SELECT * FROM episode WHERE programUUID = UUID_TO_BIN(?);
    `
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query, id); err != nil {
		return nil, err
	}
	var episodes []*model.Episode
	for _, episodeDB := range episodesDB {
		mappedEpisode := episodeDB.ToDomainModel()
		episodes = append(episodes, &mappedEpisode)
	}
	return episodes, nil
}

// Create inserts a new episode record into the database.
// It takes a context and a model.Episode, and returns an error if the operation fails.
func (adapter *episodeAdapter) Create(ctx context.Context, episode model.Episode) error {
	const query = `
        INSERT INTO episode (UUID, name, description, position, programUUID)
        VALUES (UUID_TO_BIN(:UUID), :name, :description, :position, UUID_TO_BIN(:programUUID));
    `
	var episodeDB EpisodeDB
	episodeDB.FromDomainModel(episode)
	_, err := adapter.client.db.NamedExecContext(ctx, query, episodeDB)
	return err
}

// Delete removes an episode record from the database based on its UUID.
// It takes a context and the episode's UUID, and returns an error if the operation fails.
func (adapter *episodeAdapter) Delete(ctx context.Context, episodeUUID string) error {
	const query = `
        DELETE FROM episode WHERE UUID = UUID_TO_BIN(?);
    `
	_, err := adapter.client.db.ExecContext(ctx, query, episodeUUID)
	return err
}

// Update updates an existing episode record in the database.
// It takes a context, the episode's UUID, and the updated model.Episode, and returns an error if the operation fails.
func (adapter *episodeAdapter) Update(ctx context.Context, episodeUUID string, updates model.Episode) error {
	const query = `
        UPDATE episode SET 
                             name = COALESCE(:name, name), 
                             description = COALESCE(:description, description), 
                             position = COALESCE(NULLIF(:position, 0), position),
                             programUUID = COALESCE(NULLIF(UUID_TO_BIN(:programUUID), UUID_TO_BIN('00000000-0000-0000-0000-000000000000')), programUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID);
    `
	updates.ID = episodeUUID
	var episodeDB EpisodeDB
	episodeDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, episodeDB)
	return err
}

// FindAll retrieves all episode records from the database.
// It takes a context and returns a slice of model.Episode and an error if the operation fails.
func (adapter *episodeAdapter) FindAll(ctx context.Context) ([]*model.Episode, error) {
	const query = `
        SELECT * FROM episode;
    `
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query); err != nil {
		return nil, err
	}
	var episodes []*model.Episode
	for _, episodeDB := range episodesDB {
		mappedEpisode := episodeDB.ToDomainModel()
		episodes = append(episodes, &mappedEpisode)
	}
	return episodes, nil
}

// Find retrieves an episode record from the database by its UUID.
// It takes a context and the episode's UUID, and returns a model.Episode and an error if the operation fails.
func (adapter *episodeAdapter) Find(ctx context.Context, episodeUUID string) (*model.Episode, error) {
	const query = `
        SELECT * FROM episode WHERE UUID = UUID_TO_BIN(?)
    `
	var episodesDB []*EpisodeDB
	if err := adapter.client.db.SelectContext(ctx, &episodesDB, query, episodeUUID); err != nil {
		return nil, err
	}
	if len(episodesDB) == 0 {
		return nil, nil
	}
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
// It returns the corresponding model.Episode.
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
// It sets the fields of the EpisodeDB based on the given model.Episode.
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
