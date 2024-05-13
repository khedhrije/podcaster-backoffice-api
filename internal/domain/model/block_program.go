// Package model defines the data structures for the application domain.
package model

// BlockProgram represents the association between a block and a program.
// It includes the position of the program within the block.
type BlockProgram struct {
	ID        string // Unique identifier for the block-program association
	BlockID   string // Unique identifier for the associated block
	ProgramID string // Unique identifier for the associated program
	Position  int    // Position of the program within the block
}
