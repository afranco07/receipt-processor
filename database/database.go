package database

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type InMemoryDatabase struct {
	data map[string]int
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		data: make(map[string]int),
	}
}

func (db *InMemoryDatabase) Insert(key string, value int) {
	db.data[key] = value
}

func (db *InMemoryDatabase) Get(key string) (int, error) {
	score, ok := db.data[key]
	if !ok {
		return 0, ErrNotFound
	}

	return score, nil
}
