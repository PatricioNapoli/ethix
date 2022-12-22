package store

type Database[T any] interface {
	Persist() error
	Retrieve() error

	Create(table string, data *Table[T])
	ReadOrCreate(table string) *Table[T]
}
