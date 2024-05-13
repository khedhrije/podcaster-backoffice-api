// Package pkg provides the response structs for handling JSON responses.
package pkg

// ErrorJSON represents the structure for error messages in JSON responses.
type ErrorJSON struct {
	Error string `json:"error" description:"error message"`
}

// WallResponse represents the response structure for walls.
type WallResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TagResponse represents the response structure for tags.
type TagResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProgramResponse represents the response structure for programs.
type ProgramResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// MediaResponse represents the response structure for medias.
type MediaResponse struct {
	ID         string `json:"ID"`
	DirectLink string `json:"directLink"`
	Kind       string `json:"kind"`
	EpisodeID  string `json:"episodeID"`
}

// CategoryResponse represents the response structure for categories.
type CategoryResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parentID"`
}

// BlockResponse represents the response structure for blocks.
type BlockResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

// EpisodeResponse represents the response structure for episodes.
type EpisodeResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProgramID   string `json:"programID"`
	Position    int    `json:"position"`
}

// WallBlocksResponse represents the response structure for blocks within a wall,
// including positional information.
type WallBlocksResponse struct {
	BlockResponse
	Position int `json:"position"`
}

// BlockProgramsResponse represents the response structure for programs within a block,
// including positional information.
type BlockProgramsResponse struct {
	ProgramResponse
	Position int `json:"position"`
}
