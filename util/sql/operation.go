package sql

type Operation int

const (
	SelectOperation Operation = iota
	JoinOperation
	InsertOperation
	UpdateOperation
	DeleteOperation
)
