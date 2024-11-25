package database

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/afranco07/receipt-processor/receipt"
	"github.com/google/uuid"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrReceiptAlreadyExists = errors.New("receipt already exists")
)

type InMemoryDatabase struct {
	data    map[string]int
	hashMap map[string]struct{}
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		data:    make(map[string]int),
		hashMap: make(map[string]struct{}),
	}
}

func (db *InMemoryDatabase) Insert(receipt receipt.Receipt, score int) (string, error) {
	exists, err := db.check(receipt)
	if err != nil {
		return "", err
	} else if exists {
		return "", ErrReceiptAlreadyExists
	}

	id := uuid.NewString()
	db.data[id] = score
	return id, nil
}

func (db *InMemoryDatabase) Get(key string) (int, error) {
	score, ok := db.data[key]
	if !ok {
		return 0, ErrNotFound
	}

	return score, nil
}

// check checks if the receipt has been submitted already by checking is
// sha256 hash
func (db *InMemoryDatabase) check(receipt receipt.Receipt) (bool, error) {
	b, err := json.Marshal(receipt)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(b)
	hashString := string(hash[:])

	if _, ok := db.hashMap[hashString]; ok {
		return true, nil
	}

	db.hashMap[hashString] = struct{}{}

	return false, nil
}
