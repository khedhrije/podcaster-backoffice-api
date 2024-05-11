package model

type Episode struct {
	ID          string
	Name        string
	Description string
	Position    int
	Media       Media
	ProgramID   string
}
