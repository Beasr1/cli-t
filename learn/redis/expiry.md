# Expiry Implementation

## Design Decisions

**Approach:** Hybrid (Passive + Active) like Redis

**Data Structure:**
- `map[string]StoreValue` with embedded `*time.Time`
- No separate time-ordered index (simplicity > optimization)

**Passive Deletion:**
- Check expiry on GET/SET/EXISTS
- Delete immediately if expired

**Active Deletion:**
- Background goroutine every 100ms
- Samples first 25 keys from map iteration
- Deletes expired ones

## Known Limitations

1. Random sampling might miss expired keys for a while
2. No support for key eviction policies (maxmemory-policy)
3. No persistence of expiry times

## Commands Supported

- `SET key value EX seconds` ✅
- `SET key value PX milliseconds` ✅
- `EXPIRE key seconds` ❌ (not implemented)
- `TTL key` ❌ (not implemented)

## Future Enhancements

- Sorted expiry index (min-heap or skip list)
- EXPIRE/PERSIST commands
- TTL command