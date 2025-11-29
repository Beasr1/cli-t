package store

import (
	"cli-t/internal/shared/logger"
	"context"
	"time"
)

// New creates a new in-memory store
func New() Store {
	// Later you might add:
	// - Load data from disk
	// - Pre-populate some keys
	// - Start background cleanup goroutines
	// - Initialize metrics
	// fayada of hacing interface pattern
	return &InMemoryStore{
		data: make(map[string]StoreValue),
	}
}

// Private helper - assumes lock is already held!
func (s *InMemoryStore) getIfValid(key string) (StoreValue, bool) {
	val, exists := s.data[key]
	if !exists {
		return StoreValue{}, false
	}

	if val.ExpiresAt != nil && time.Now().After(*val.ExpiresAt) {
		delete(s.data, key)
		return StoreValue{}, false
	}

	return val, true
}

// Get retrieves a value by key
// Returns (value, true) if found, ("", false) if not found
func (s *InMemoryStore) Get(key string) (StoreValue, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getIfValid(key)
}

// Set stores a key-value pair
func (s *InMemoryStore) Set(key string, value StoreValue) {
	s.mu.Lock() // Exclusive lock - blocks all readers and writers
	defer s.mu.Unlock()

	s.data[key] = value
}

// Add this method to InMemoryStore
func (s *InMemoryStore) CleanExpiredKeys() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	deleted := 0

	// Sample random keys (map iteration is random in Go!)
	for key, val := range s.data {
		// Check if expired
		if val.ExpiresAt != nil && now.After(*val.ExpiresAt) {
			delete(s.data, key)
			deleted++

			// Sample only ~25 keys per iteration (like Redis)
			if deleted >= 25 {
				break
			}
		}

		// Can Also limit total keys checked (in case map is huge)
	}

	return deleted
}

// StartExpiryWorker runs background expiry cleanup
// Call this from Server.Start()
func (s *InMemoryStore) StartExpiryWorker(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown
			logger.Debug("Expiry worker stopping")
			return

		case <-ticker.C:
			// Run cleanup
			deleted := s.CleanExpiredKeys()
			if deleted > 0 {
				logger.Verbose("Cleaned expired keys", "count", deleted)
			}
		}
	}
}

// GetTTL returns time-to-live in seconds
// Returns:
//
//	-2 if key doesn't exist
//	-1 if key exists but has no expiry
//	>0 seconds until expiry
func (s *InMemoryStore) GetTTL(key string) int64 {
	// TODO: Implement this
	// Hints:
	// - What lock do you need? (reading only)
	// - How to check if key exists?
	// - How to check if it has expiry?
	// - How to calculate time remaining? (time.Until)
	// - What if time.Until is negative? (already expired)
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.getIfValid(key)
	if !exists {
		return -2
	}

	if val.ExpiresAt == nil {
		return -1
	}

	ttl := time.Until(*val.ExpiresAt)
	if ttl < 0 {
		delete(s.data, key)
		return -2
	}

	return int64(ttl.Seconds())
}
