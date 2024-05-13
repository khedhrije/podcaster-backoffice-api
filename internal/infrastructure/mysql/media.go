// Package mysql provides MySQL implementations of the persistence interfaces.
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
// It returns an implementation of the MediaPersister interface.
func NewMediaAdapter(client *client) port.MediaPersister {
	return &mediaAdapter{
		client: client,
	}
}

// Create inserts a new media record into the database.
// It takes a context and a model.Media, and returns an error if the operation fails.
func (adapter *mediaAdapter) Create(ctx context.Context, media model.Media) error {
	const query = `
        INSERT INTO media (UUID, direct_link, kind, episodeUUID)
        VALUES (UUID_TO_BIN(:UUID), :direct_link, :kind, UUID_TO_BIN(:episodeUUID))
    `
	var mediaDB MediaDB
	mediaDB.FromDomainModel(media)
	_, err := adapter.client.db.NamedExecContext(ctx, query, mediaDB)
	return err
}

// Delete removes a media record from the database based on its UUID.
// It takes a context and the media's UUID, and returns an error if the operation fails.
func (adapter *mediaAdapter) Delete(ctx context.Context, mediaUUID string) error {
	const query = `
        DELETE FROM media WHERE UUID = UUID_TO_BIN(?);
    `
	_, err := adapter.client.db.ExecContext(ctx, query, mediaUUID)
	return err
}

// Update updates an existing media record in the database.
// It takes a context, the media's UUID, and the updated model.Media, and returns an error if the operation fails.
func (adapter *mediaAdapter) Update(ctx context.Context, mediaUUID string, updates model.Media) error {
	const query = `
        UPDATE media SET 
                             direct_link = COALESCE(:direct_link, direct_link), 
                             kind = COALESCE(:kind, kind), 
                             episodeUUID = COALESCE(NULLIF(UUID_TO_BIN(:episodeUUID), UUID_TO_BIN('00000000-0000-0000-0000-000000000000')), episodeUUID)
                        WHERE UUID = UUID_TO_BIN(:UUID)
    `
	updates.ID = mediaUUID
	var mediaDB MediaDB
	mediaDB.FromDomainModel(updates)
	_, err := adapter.client.db.NamedExecContext(ctx, query, mediaDB)
	return err
}

// FindAll retrieves all media records from the database.
// It takes a context and returns a slice of model.Media and an error if the operation fails.
func (adapter *mediaAdapter) FindAll(ctx context.Context) ([]*model.Media, error) {
	const query = `
        SELECT * FROM media;
    `
	var mediaDB []*MediaDB
	if err := adapter.client.db.SelectContext(ctx, &mediaDB, query); err != nil {
		return nil, err
	}
	var media []*model.Media
	for _, mediaEntry := range mediaDB {
		mappedMedia := mediaEntry.ToDomainModel()
		media = append(media, &mappedMedia)
	}
	return media, nil
}

// Find retrieves a media record from the database by its UUID.
// It takes a context and the media's UUID, and returns a model.Media and an error if the operation fails.
func (adapter *mediaAdapter) Find(ctx context.Context, mediaUUID string) (*model.Media, error) {
	const query = `
        SELECT * FROM media WHERE UUID = UUID_TO_BIN(?);
    `
	var mediaDB MediaDB
	if err := adapter.client.db.GetContext(ctx, &mediaDB, query, mediaUUID); err != nil {
		return nil, err
	}
	result := mediaDB.ToDomainModel()
	return &result, nil
}

// MediaDB is a struct representing the media database model.
type MediaDB struct {
	UUID       uuid.UUID      `db:"UUID"`
	DirectLink sql.NullString `db:"direct_link"`
	Kind       sql.NullString `db:"kind"`
	EpisodeID  uuid.UUID      `db:"episodeUUID"`
}

// ToDomainModel converts a MediaDB database model to a model.Media domain model.
// It returns the corresponding model.Media.
func (db *MediaDB) ToDomainModel() model.Media {
	return model.Media{
		ID:         db.UUID.String(),
		DirectLink: db.DirectLink.String,
		Kind:       db.Kind.String,
		EpisodeID:  db.EpisodeID.String(),
	}
}

// FromDomainModel converts a model.Media domain model to a MediaDB database model.
// It sets the fields of the MediaDB based on the given model.Media.
func (db *MediaDB) FromDomainModel(domain model.Media) {
	db.UUID = uuid.MustParse(domain.ID)
	db.DirectLink = sql.NullString{String: domain.DirectLink, Valid: domain.DirectLink != ""}
	db.Kind = sql.NullString{String: domain.Kind, Valid: domain.Kind != ""}
	db.EpisodeID = uuid.Nil
	if domain.EpisodeID != "" {
		db.EpisodeID = uuid.MustParse(domain.EpisodeID)
	}
}
