package database

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

func (db *InMemoryDatabase) Get(key string) int {
	return db.data[key]
}
