package mysql

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/model"
	"github.com/khedhrije/podcaster-backoffice-api/internal/domain/port"
)

// mediaAdapter is a struct that acts as an adapter for interacting with
// the media data in the MySQL database.
type mediaAdapter struct {
	client *client
}

// NewMediaAdapter creates a new media adapter with the provided MySQL client.
func NewMediaAdapter(client *client) port.MediaPersister {
	return &mediaAdapter{
		client: client,
	}
}

// Create inserts a new media record into the database.
func (adapter *mediaAdapter) Create(ctx context.Context, media model.Media) error {
	// SQL query to insert a new media record
	const query = `
        INSERT INTO media (UUID, direct_link, kind, episodeUUID)
        VALUES (UUID_TO_BIN(:UUID), :direct_link, :kind, UUID_TO_BIN(:episodeUUID))
    `
	// Convert domain model to database model
	var mediaDB MediaDB
	mediaDB.FromDomainModel(media)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, mediaDB)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a media record from the database based on its UUID.
func (adapter *mediaAdapter) Delete(ctx context.Context, mediaUUID string) error {
	// SQL query to delete a media record by UUID
	const query = `
        DELETE FROM media WHERE UUID = UUID_TO_BIN(?);
    `
	// Execute named query
	_, err := adapter.client.db.ExecContext(ctx, query, mediaUUID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an existing media record in the database.
func (adapter *mediaAdapter) Update(ctx context.Context, mediaUUID string, updates model.Media) error {
	// SQL query to update a media record
	const query = `
        UPDATE media SET 
                             direct_link = COALESCE(:direct_link, direct_link), 
                             kind = COALESCE(:kind, kind), 
                             episodeUUID = COALESCE(NULLIF(UUID_TO_BIN(:episodeUUID), UUID_TO_BIN('00000000-0000-0000-0000-000000000000')), episodeUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `

	// Set UUID for updates
	updates.ID = mediaUUID
	// Convert domain model to database model
	var mediaDB MediaDB
	mediaDB.FromDomainModel(updates)
	// Execute named query
	_, err := adapter.client.db.NamedExecContext(ctx, query, mediaDB)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves all media records from the database.
func (adapter *mediaAdapter) FindAll(ctx context.Context) ([]*model.Media, error) {
	// SQL query to select all media records
	const query = `
        SELECT * FROM media;
    `
	// Execute query and retrieve results
	var mediaDB []*MediaDB
	if err := adapter.client.db.SelectContext(ctx, &mediaDB, query); err != nil {
		return nil, err
	}
	// Convert database models to domain models
	var media []*model.Media
	for _, mediaEntry := range mediaDB {
		mappedMedia := mediaEntry.ToDomainModel()
		media = append(media, &mappedMedia)
	}

	return media, nil
}

// Find retrieves a media record from the database by its UUID.
func (adapter *mediaAdapter) Find(ctx context.Context, mediaUUID string) (*model.Media, error) {
	// SQL query to select a media record by UUID
	const query = `
        SELECT * FROM media WHERE UUID = UUID_TO_BIN(?);
    `
	// Execute query and retrieve results
	var mediaDB MediaDB
	if err := adapter.client.db.GetContext(ctx, &mediaDB, query, mediaUUID); err != nil {
		return nil, err
	}
	// Convert database model to domain model
	result := mediaDB.ToDomainModel()
	return &result, nil
}

// MediaDB is a struct representing the media dactabase model.
type MediaDB struct {
	UUID       uuid.UUID      `db:"UUID"`
	DirectLink sql.NullString `db:"direct_link"`
	Kind       sql.NullString `db:"kind"`
	EpisodeID  uuid.UUID      `db:"episodeUUID"`
}

// ToDomainModel converts a MediaDB database model to a model.Media domain model.
func (db *MediaDB) ToDomainModel() model.Media {
	return model.Media{
		ID:         db.UUID.String(),
		DirectLink: db.DirectLink.String,
		Kind:       db.Kind.String,
		EpisodeID:  db.EpisodeID.String(),
	}
}

// FromDomainModel converts a model.Media domain model to a MediaDB database model.
func (db *MediaDB) FromDomainModel(domain model.Media) {
	db.UUID = uuid.MustParse(domain.ID)
	db.DirectLink = sql.NullString{String: domain.DirectLink, Valid: domain.DirectLink != ""}
	db.Kind = sql.NullString{String: domain.Kind, Valid: domain.Kind != ""}
	db.EpisodeID = uuid.Nil
	if domain.EpisodeID != "" {
		db.EpisodeID = uuid.MustParse(domain.EpisodeID)
	}

}
