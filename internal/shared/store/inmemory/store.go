package inmemory

import (
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"strconv"
	"time"
)

// New creates a new inmemory store
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

// Private helper - assumes lock is already held!
func (s *InMemoryStore) set(key string, value StoreValue) {
	s.data[key] = value
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

	s.set(key, value)
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

// SetExpiry sets an expiration time on an existing key
// Returns true if key exists and expiry was set
// Returns false if key doesn't exist
func (s *InMemoryStore) SetExpiry(key string, seconds int) bool {
	// TODO: Implement this
	// Hints:
	// - What lock do you need? (you're modifying)
	// - Check if key exists (use getIfValid helper?)
	// - Calculate expiry time (time.Now().Add(...))
	// - Update the ExpiresAt field
	// - Return true/false based on success

	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.getIfValid(key)
	if !exists {
		return false
	}

	expiresAt := time.Now().Add(time.Duration(seconds) * time.Second)
	val.ExpiresAt = &expiresAt
	s.set(key, val)
	return true
}

func (s *InMemoryStore) Exists(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, key := range keys {
		_, exists := s.getIfValid(key)
		if exists {
			count++
		}
	}
	return count
}

func (s *InMemoryStore) Delete(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, key := range keys {
		_, exists := s.getIfValid(key)
		if exists {
			delete(s.data, key)
			count++
		}
	}
	return count
}

func (s *InMemoryStore) Incr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get current value
	val, exists := s.getIfValid(key)

	var currentNum int64
	if exists {

		// Check if it's a string
		if val.Type != TypeString {
			return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}

		// Parse as integer
		// If parse fails, return error
		num, err := strconv.ParseInt(val.Data, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("ERR value is not an integer or out of range")
		}

		currentNum = num
	} else {
		// Doesn't exist, start at 0
		currentNum = 0
	}

	// Increment
	newValue := currentNum + 1

	// Store back as string
	// Return new value
	val.Data = strconv.FormatInt(newValue, 10)
	val.Type = TypeString
	s.set(key, val)
	return newValue, nil
}

func (s *InMemoryStore) Decr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get current value
	val, exists := s.getIfValid(key)

	var currentNum int64
	if exists {

		// Check if it's a string
		if val.Type != TypeString {
			return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}

		// Parse as integer
		// If parse fails, return error
		num, err := strconv.ParseInt(val.Data, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("ERR value is not an integer or out of range")
		}

		currentNum = num
	} else {
		// Doesn't exist, start at 0
		currentNum = 0
	}

	// Increment
	newValue := currentNum - 1

	// Store back as string
	// Return new value
	val.Data = strconv.FormatInt(newValue, 10)
	val.Type = TypeString
	s.set(key, val)
	return newValue, nil
}

func (s *InMemoryStore) LPush(key string, values ...string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.getIfValid(key)

	var list []string
	if exists {
		// Key exists - check type
		if val.Type != TypeList {
			return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		list = val.List
	} else {
		// New list
		list = []string{}
	}

	// Push to LEFT means PREPEND
	newList := append(values, list...)

	// Store back
	s.set(key, StoreValue{
		Type:      TypeList,
		List:      newList,
		ExpiresAt: val.ExpiresAt,
	})

	return int64(len(newList)), nil
}

func (s *InMemoryStore) RPush(key string, values ...string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.getIfValid(key)

	var list []string
	if exists {
		// Key exists - check type
		if val.Type != TypeList {
			return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		list = val.List
	} else {
		// New list
		list = []string{}
	}

	// Push to LEFT means PREPEND
	newList := append(list, values...)

	// Store back
	s.set(key, StoreValue{
		Type:      TypeList,
		List:      newList,
		ExpiresAt: val.ExpiresAt,
	})

	return int64(len(newList)), nil
}

func (s *InMemoryStore) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock() // Read-only!
	defer s.mu.RUnlock()

	val, exists := s.getIfValid(key)
	if !exists {
		return []string{}, nil // Empty array for non-existent
	}

	if val.Type != TypeList {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	list := val.List
	length := len(list)

	// Handle negative indices
	// -1 means last element, -2 means second-to-last, etc.
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	// Bounds checking
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return []string{}, nil
	}

	// Slice syntax: [start:stop+1]
	// (stop is inclusive in Redis!)
	return list[start : stop+1], nil
}
