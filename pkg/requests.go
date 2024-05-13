package pkg

// CreateWallRequestJSON represents a JSON request for creating walls.
type CreateWallRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
}

// Name returns the name of the wall.
func (req CreateWallRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the wall.
func (req CreateWallRequestJSON) Description() string {
	return req.DescriptionJSON
}

// UpdateWallRequestJSON represents a JSON request for updating walls.
type UpdateWallRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
}

// Name returns the name of the wall update request.
func (req UpdateWallRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the wall update request.
func (req UpdateWallRequestJSON) Description() string {
	return req.DescriptionJSON
}

// CreateTagRequestJSON represents a JSON request for creating tags.
type CreateTagRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
}

// Name returns the name of the tag.
func (req CreateTagRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the tag.
func (req CreateTagRequestJSON) Description() string {
	return req.DescriptionJSON
}

// UpdateTagRequestJSON represents a JSON request for updating tags.
type UpdateTagRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
}

// Name returns the name of the tag update request.
func (req UpdateTagRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the tag update request.
func (req UpdateTagRequestJSON) Description() string {
	return req.DescriptionJSON
}

// CreateProgramRequestJSON represents a JSON request for creating programs.
type CreateProgramRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
}

// Name returns the name of the program.
func (req CreateProgramRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the program.
func (req CreateProgramRequestJSON) Description() string {
	return req.DescriptionJSON
}

// UpdateProgramRequestJSON represents a JSON request for updating programs.
type UpdateProgramRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
}

// Name returns the name of the program update request.
func (req UpdateProgramRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the program update request.
func (req UpdateProgramRequestJSON) Description() string {
	return req.DescriptionJSON
}

// CreateMediaRequestJSON represents a JSON request for creating medias.
type CreateMediaRequestJSON struct {
	DirectLinkJSON string `json:"directLink"`
	KindJSON       string `json:"kind"`
}

// DirectLink returns the directLink of the media.
func (req CreateMediaRequestJSON) DirectLink() string {
	return req.DirectLinkJSON
}

// Kind returns the kind of the media.
func (req CreateMediaRequestJSON) Kind() string {
	return req.KindJSON
}

// UpdateMediaRequestJSON represents a JSON request for updating medias.
type UpdateMediaRequestJSON struct {
	DirectLinkJSON string `json:"directLink"`
	KindJSON       string `json:"kind"`
}

// DirectLink returns the directLink of the media update request.
func (req UpdateMediaRequestJSON) DirectLink() string {
	return req.DirectLinkJSON
}

// Kind returns the kind of the media update request.
func (req UpdateMediaRequestJSON) Kind() string {
	return req.KindJSON
}

// CreateCategoryRequestJSON represents a JSON request for creating categories.
type CreateCategoryRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
	ParentIDJSON    string `json:"parentID"`
}

// Name returns the name of the category.
func (req CreateCategoryRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the category.
func (req CreateCategoryRequestJSON) Description() string {
	return req.DescriptionJSON
}

// ParentID returns the parentID of the category.
func (req CreateCategoryRequestJSON) ParentID() string {
	return req.ParentIDJSON
}

// UpdateCategoryRequestJSON represents a JSON request for updating categories.
type UpdateCategoryRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
	ParentIDJSON    string `json:"parentID"`
}

// Name returns the name of the category update request.
func (req UpdateCategoryRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the category update request.
func (req UpdateCategoryRequestJSON) Description() string {
	return req.DescriptionJSON
}

// ParentID returns the parentID of the category.
func (req UpdateCategoryRequestJSON) ParentID() string {
	return req.ParentIDJSON
}

// CreateBlockRequestJSON represents a JSON request for creating blocks.
type CreateBlockRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
	KindJSON        string `json:"kind"`
}

// Name returns the name of the block.
func (req CreateBlockRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the block.
func (req CreateBlockRequestJSON) Description() string {
	return req.DescriptionJSON
}

// Kind returns the kind of the block.
func (req CreateBlockRequestJSON) Kind() string {
	return req.KindJSON
}

// UpdateBlockRequestJSON represents a JSON request for updating blocks.
type UpdateBlockRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
	KindJSON        string `json:"kind"`
}

// Name returns the name of the block update request.
func (req UpdateBlockRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the block update request.
func (req UpdateBlockRequestJSON) Description() string {
	return req.DescriptionJSON
}

// Kind returns the kind of the block.
func (req UpdateBlockRequestJSON) Kind() string {
	return req.KindJSON
}

// CreateEpisodeRequestJSON represents a JSON request for creating episodes.
type CreateEpisodeRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
	ProgramIDJSON   string `json:"programID"`
	PositionJSON    int    `json:"position"`
}

// Name returns the name of the episode.
func (req CreateEpisodeRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the episode.
func (req CreateEpisodeRequestJSON) Description() string {
	return req.DescriptionJSON
}

// ProgramID returns the programID of the episode.
func (req CreateEpisodeRequestJSON) ProgramID() string {
	return req.ProgramIDJSON
}

// Position returns the position of the episode.
func (req CreateEpisodeRequestJSON) Position() int {
	return req.PositionJSON
}

// UpdateEpisodeRequestJSON represents a JSON request for updating episodes.
type UpdateEpisodeRequestJSON struct {
	NameJSON        string `json:"name"`
	DescriptionJSON string `json:"description"`
	ProgramIDJSON   string `json:"programID"`
	PositionJSON    int    `json:"position"`
}

// Name returns the name of the episode update request.
func (req UpdateEpisodeRequestJSON) Name() string {
	return req.NameJSON
}

// Description returns the description of the episode update request.
func (req UpdateEpisodeRequestJSON) Description() string {
	return req.DescriptionJSON
}

// ProgramID returns the programID of the episode.
func (req UpdateEpisodeRequestJSON) ProgramID() string {
	return req.ProgramIDJSON
}

// Position returns the position of the episode.
func (req UpdateEpisodeRequestJSON) Position() int {
	return req.PositionJSON
}
