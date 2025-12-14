package inmemory

import (
	"context"
	"sync"
	"time"
)

type ValueType string

const (
	TypeString ValueType = "string"
	TypeList   ValueType = "list"
)

type StoreValue struct {
	Type      ValueType
	Data      string     // For strings
	List      []string   // For lists
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
	LPush(key string, values ...string) (int64, error)
	RPush(key string, values ...string) (int64, error)
	LRange(key string, start, stop int) ([]string, error)
	// We'll add more methods later (Delete, Exists, etc.)
}

// InMemoryStore is a thread-safe inmemory key-value store
type InMemoryStore struct {
	data map[string]StoreValue
	mu   sync.RWMutex // ← Use RWMutex, not Mutex!
}
