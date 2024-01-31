package db

import (
	"sync"
	"time"
)

// InMemoryDb represents an in-memory database
type InMemoryDb struct {
	m   map[string]string
	mux sync.RWMutex
}

func NewDB() *InMemoryDb {
	db := &InMemoryDb{
		m: make(map[string]string),
	}
	return db
}

// Set sets the value for the specified key in the inMemoryDb.
// It locks the mutex, updates the map with the key-value pair,
// and then unlocks the mutex.
// Additionally, it sets a timer to automatically delete the key
// from the map after 3 minutes.
func (db *InMemoryDb) Set(key string, value string) {
	db.mux.Lock()
	db.m[key] = value
	db.mux.Unlock()

	// Set the key to be automatically deleted after 3 minutes
	time.AfterFunc(3*time.Minute, func() {
		db.Delete(key)
	})
}

// Get retrieves a value for a key
func (db *InMemoryDb) Get(key string) (string, bool) {
	db.mux.RLock()
	value, exists := db.m[key]
	db.mux.RUnlock()
	return value, exists
}

// Delete removes a value for a key
func (db *InMemoryDb) Delete(key string) {
	db.mux.Lock()
	delete(db.m, key)
	db.mux.Unlock()
}
