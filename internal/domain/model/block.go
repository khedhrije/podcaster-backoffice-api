package model

type Block struct {
	ID          string
	Name        string
	Description string
	Kind        string
	Programs    []Program
}
