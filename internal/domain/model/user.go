// Package model defines the data structures for the application domain.
package model

// User represents a user entity in the system.
// A User is an individual who interacts with the application, identified by their unique email.
type User struct {
	ID        string // Unique identifier for the user
	Firstname string // First name of the user
	Lastname  string // Last name of the user
	Email     string // Email address of the user
}
