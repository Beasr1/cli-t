package store

import (
	"context"
	"sync"
	"time"
)

type StoreValue struct {
	Data      string
	ExpiresAt *time.Time // nil = no expiry
}

// Store defines the key-value storage interface
type Store interface {
	Get(key string) (StoreValue, bool)
	Set(key string, value StoreValue)
	CleanExpiredKeys() int
	StartExpiryWorker(ctx context.Context)
	GetTTL(key string) int64
	SetExpiry(key string, seconds int) bool
	Exists(keys ...string) int
	Delete(keys ...string) int
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	// We'll add more methods later (Delete, Exists, etc.)
}

// InMemoryStore is a thread-safe in-memory key-value store
type InMemoryStore struct {
	data map[string]StoreValue
	mu   sync.RWMutex // ← Use RWMutex, not Mutex!
}
