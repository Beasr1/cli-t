package inmemory_test

import (
	inmemory "cli-t/internal/shared/store/inmemory"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test concurrent Gets with expiry
func TestConcurrentGetWithExpiry(t *testing.T) {
	s := inmemory.New()

	// Set a key with 100ms expiry
	expiry := time.Now().Add(100 * time.Millisecond)
	s.Set("key", inmemory.StoreValue{
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

// In inmemory/store_test.go
func TestGetTTL(t *testing.T) {
	s := inmemory.New()

	// Test 1: Non-existent key
	ttl := s.GetTTL("nonexistent")
	assert.Equal(t, ttl, int64(-2))

	// Test 2: Key with no expiry
	s.Set("permanent", inmemory.StoreValue{Data: "val", ExpiresAt: nil})
	ttl = s.GetTTL("permanent")
	assert.Equal(t, ttl, int64(-1))

	// Test 3: Key with expiry
	expiry := time.Now().Add(5 * time.Second)
	s.Set("temp", inmemory.StoreValue{Data: "val", ExpiresAt: &expiry})
	ttl = s.GetTTL("temp")
	assert.GreaterOrEqual(t, ttl, int64(4)) // At least 4 seconds
	assert.LessOrEqual(t, ttl, int64(5))    // At most 5 seconds

	// Test 4: Expired key (edge case!)
	pastExpiry := time.Now().Add(-1 * time.Second)
	s.Set("expired", inmemory.StoreValue{Data: "val", ExpiresAt: &pastExpiry})
	ttl = s.GetTTL("expired")
	assert.Equal(t, ttl, int64(-2))
}

func TestSetExpiry(t *testing.T) {
	s := inmemory.New()

	// Test 1: Set expiry on existing key
	s.Set("mykey", inmemory.StoreValue{Data: "value", ExpiresAt: nil})
	success := s.SetExpiry("mykey", 5)
	assert.True(t, success, "Should succeed on existing key")

	// Verify TTL
	ttl := s.GetTTL("mykey")
	assert.GreaterOrEqual(t, ttl, int64(4))
	assert.LessOrEqual(t, ttl, int64(5))

	// Test 2: Set expiry on non-existent key
	success = s.SetExpiry("nonexistent", 10)
	assert.False(t, success, "Should fail on non-existent key")

	// Test 3: Override existing expiry
	expiry := time.Now().Add(100 * time.Second)
	s.Set("override", inmemory.StoreValue{Data: "val", ExpiresAt: &expiry})
	ttl = s.GetTTL("override")
	assert.Greater(t, ttl, int64(90)) // Should be ~100

	success = s.SetExpiry("override", 5)
	assert.True(t, success)
	ttl = s.GetTTL("override")
	assert.LessOrEqual(t, ttl, int64(5)) // Now should be ~5

	// Test 4: Key expires after time
	s.Set("temp", inmemory.StoreValue{Data: "val", ExpiresAt: nil})
	s.SetExpiry("temp", 1) // 1 second

	time.Sleep(1100 * time.Millisecond) // Wait for expiry

	_, exists := s.Get("temp")
	assert.False(t, exists, "Key should be expired")

	// Test 5: EXPIRE 0 (immediate expiry)
	s.Set("immediate", inmemory.StoreValue{Data: "val", ExpiresAt: nil})
	success = s.SetExpiry("immediate", 0)
	assert.True(t, success)

	// Key should expire immediately
	ttl = s.GetTTL("immediate")
	assert.Equal(t, int64(-2), ttl, "Key with 0 second expiry should be gone")
}
