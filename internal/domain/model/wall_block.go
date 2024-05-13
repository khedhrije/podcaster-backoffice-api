// Package model defines the data structures for the application domain.
package model

// WallBlock represents the association between a wall and a block.
// It includes the position of the block within the wall.
type WallBlock struct {
	ID       string // Unique identifier for the wall-block association
	WallID   string // Unique identifier for the associated wall
	BlockID  string // Unique identifier for the associated block
	Position int    // Position of the block within the wall
}
