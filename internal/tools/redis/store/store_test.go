package store_test

import (
	"sync"
	"testing"
	"time"

	"cli-t/internal/tools/redis/store"

	"github.com/stretchr/testify/assert"
)

// Test concurrent Gets with expiry
func TestConcurrentGetWithExpiry(t *testing.T) {
	s := store.New()

	// Set a key with 100ms expiry
	expiry := time.Now().Add(100 * time.Millisecond)
	s.Set("key", store.StoreValue{
		Data:      "value",
		ExpiresAt: &expiry,
	})

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // ← Decrement when done
			// Call Get repeatedly over 200ms
			for j := 0; j < 20; j++ {
				s.Get("key")
				time.Sleep(10 * time.Millisecond)
			}
		}()
	}

	wg.Wait() // ← Wait for all to finish
}

/*
Goroutine 15 read from address 0x00c00009f170
Goroutine 10 wrote to same address
Both held RLock (not mutually exclusive)
Not synchronized → RACE!


Based on ThreadSanitizer (Google's Tool)
Go's race detector is built on ThreadSanitizer (TSan), originally from Google.

*/

// In store/store_test.go
func TestGetTTL(t *testing.T) {
	s := store.New()

	// Test 1: Non-existent key
	ttl := s.GetTTL("nonexistent")
	assert.Equal(t, ttl, int64(-2))

	// Test 2: Key with no expiry
	s.Set("permanent", store.StoreValue{Data: "val", ExpiresAt: nil})
	ttl = s.GetTTL("permanent")
	assert.Equal(t, ttl, int64(-1))

	// Test 3: Key with expiry
	expiry := time.Now().Add(5 * time.Second)
	s.Set("temp", store.StoreValue{Data: "val", ExpiresAt: &expiry})
	ttl = s.GetTTL("temp")
	assert.GreaterOrEqual(t, ttl, int64(4)) // At least 4 seconds
	assert.LessOrEqual(t, ttl, int64(5))    // At most 5 seconds

	// Test 4: Expired key (edge case!)
	pastExpiry := time.Now().Add(-1 * time.Second)
	s.Set("expired", store.StoreValue{Data: "val", ExpiresAt: &pastExpiry})
	ttl = s.GetTTL("expired")
	assert.Equal(t, ttl, int64(-2))
	// What should ttl be? Should you delete it?
}
