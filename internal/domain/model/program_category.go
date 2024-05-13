// Package model defines the data structures for the application domain.
package model

// ProgramCategory represents the association between a program and a category.
// It links a program to a specific category.
type ProgramCategory struct {
	ID         string // Unique identifier for the program-category association
	ProgramID  string // Unique identifier for the associated program
	CategoryID string // Unique identifier for the associated category
}
