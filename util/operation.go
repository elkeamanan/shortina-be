package util

type SQLOperation int

const (
	SelectOperation SQLOperation = iota
	JoinOperation
	InsertOperation
	UpdateOperation
	DeleteOperation
)
