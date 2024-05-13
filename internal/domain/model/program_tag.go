// Package model defines the data structures for the application domain.
package model

// ProgramTag represents the association between a program and a tag.
// It links a program to a specific tag.
type ProgramTag struct {
	ID        string // Unique identifier for the program-tag association
	ProgramID string // Unique identifier for the associated program
	TagID     string // Unique identifier for the associated tag
}
