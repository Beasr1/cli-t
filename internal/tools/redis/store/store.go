package store

import "sync"

// Store defines the key-value storage interface
type Store interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	// We'll add more methods later (Delete, Exists, etc.)
}

// InMemoryStore is a thread-safe in-memory key-value store
type InMemoryStore struct {
	data map[string]string
	mu   sync.RWMutex // ← Use RWMutex, not Mutex!
}

// New creates a new in-memory store
func New() Store {
	// Later you might add:
	// - Load data from disk
	// - Pre-populate some keys
	// - Start background cleanup goroutines
	// - Initialize metrics
	// fayada of hacing interface pattern
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

// Get retrieves a value by key
// Returns (value, true) if found, ("", false) if not found
func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.RLock() // Multiple Gets can happen simultaneously!
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Set stores a key-value pair
func (s *InMemoryStore) Set(key string, value string) {
	s.mu.Lock() // Exclusive lock - blocks all readers and writers
	defer s.mu.Unlock()

	s.data[key] = value
}
