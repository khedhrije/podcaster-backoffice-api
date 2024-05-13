package pkg

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
}

type WallBlocksResponse struct {
	BlockResponse
	Position int
}

type BlockProgramsResponse struct {
	ProgramResponse
	Position int
}
