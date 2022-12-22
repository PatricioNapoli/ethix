package store

import (
	"encoding/gob"
	"log"
	"os"
)

type GobDatabase[T any] struct {
	File string

	tables map[string]*Table[T]
}

func NewGobDatabase[T any](file string) Database[T] {
	db := &GobDatabase[T]{
		File:   file,
		tables: map[string]*Table[T]{},
	}

	// Error discard, database could be new
	db.Retrieve()

	return db
}

func (db *GobDatabase[T]) Persist() error {
	file, err := os.Create(db.File)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(db.tables)
	if err != nil {
		return err
	}

	return nil
}

func (db *GobDatabase[T]) Retrieve() error {
	log.Printf("retrieving gob database from %s", db.File)

	gob.Register(map[string]*Table[T]{})

	file, err := os.Open(db.File)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&db.tables)

	if err != nil {
		return err
	}

	// Initialize non persisted mutex
	for _, t := range db.tables {
		t.Init()
	}

	return nil
}

func (db *GobDatabase[T]) Create(table string, data *Table[T]) {
	db.tables[table] = data
}

func (db *GobDatabase[T]) ReadOrCreate(table string) *Table[T] {
	if val, exists := db.tables[table]; exists {
		return val
	}

	data := NewTable[T]()
	db.Create(table, data)
	return db.tables[table]
}
