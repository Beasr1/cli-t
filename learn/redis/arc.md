# Redis Server Implementation

## Overview
redis alt : server that communicates using resp protocol and is simple data store map
fast af and efficient data retrival since it stores data within memory (local storage) --> that can be persisted as well
can handle concurrency, one to one data store, can help reduce the load on db for frequent data point access
tries to do all things redis does/did

## RESP
Simple String: "+OK\r\n"
Error: "-ERR unknown command\r\n"
Integer: ":42\r\n"
Bulk String: "$5\r\nhello\r\n" or "$-1\r\n" (null)
Array: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" or "*-1\r\n" (null)
resp is used since we can send this as stream of bytes
WHY does Redis use this protocol instead of JSON? resp is easier to parse as it tells has metadata of the data its parsing. json is not convient as it does not handle well with steam and will not be able to parse without loading in whole message first
What's significant about \r\n terminators? they help with metadata : size of bulk string
Why have 5 different types?  String is fast, Bulk String handles binary data, Integer is efficient, Error is distinguishable, Array groups things. optimizes things


## Components
server (listners that listen in the requests of the clients)
per client spawns one goroutines. can handle n number of clients concurrently
when client disconnect goroutiones stopes gracefully (no memory leak)
shutdown is graceful. first listners stops -> giving time for active requests that were being processed

parser that converts []bytes as resp format --> resp values
What's hard about parsing RESP? buffer limt, partition over multiple packets. one buffer can have multiple commands
How do you handle partial reads? Currently: Assumes full message in buffer (4KB limit). Doesn't handle partial reads - this is a known limitation. TODO: fix that
What's the buffer limitation? accepting steam of request in buffer. buffer size limit limatation for now
Buffer is 4KB. Messages larger than 4KB fail. Partial messages across multiple TCP packets aren't handled. Needs streaming parser to fix

map is used since its fast af o(1) lookup
locks are used (read locks are used when there is read>>>> writes)

combination of multiple resp values makes up a command
redis-cli communitcates in RESP protocol : bulk strings



v0 originally supported string
now it also supports list
Why not create separate maps for strings vs lists? accessing 2 data stores will add to redudnancy. also will have to take care of key collision and cleanup overhead. single source truth
What are the tradeoffs of having both Data and List fields? will need empty 24 bytes minimum for pointer ref even if unused
How much memory is "wasted" when a string doesn't use List field? memory occupied by pointer
If you have 1 million string keys, how much memory is "wasted"? (24 MB)
 ~24 bytes overhead per string, ~16 bytes overhead per list, acceptable tradeoff for simpler code, memory vs complexity.
Could you use interface{} instead? Why didn't you? interface would have been too generic. ype assertions needed everywhere, less explicit, harder to debug, lose compile-time safety, harder to see structure.


Why use type ValueType string instead of just string? could have but would be better for verbosity. Type safety, autocomplete in IDE, prevents typos, documents intent, compile-time checking.
What happens if Type is empty string ""? no type set to storevalue. throws wrong type error
Could you use iota here? Why or why not? idk. string is more readable in logs/debugging, easier to JSON marshal, self-documenting, small enum so no performance benefit from int.


type StoreValue struct {
    Type      ValueType    // string = 16 bytes (pointer + len)
    Data      string       // string = 16 bytes  
    List      []string     // slice = 24 bytes (pointer + len + cap)
    ExpiresAt *time.Time   // pointer = 8 bytes
}

## TCP
redis uses tcp 
forms a connection one time and listen in the requests from same client
we can;t risk losing packets and order matters as well for correct parsing of commands
persistent connection matters since connecting client again using handshake will take longer



## Data store
map with usage of locks
pointer of time is used since data store since it uses 8 butes tather than memory of default time value


## expiry
hybrid expiry (passive and active)
passive expiry check whether data point is expired when data is accessed
active also keeps on checking
both is better since active expiry cron will not need to be tight and its necessary to delete the stale data if not accessed to prevent memory bloating



## Bugs Fixed
1. **RLock → Lock:** Used RLock in Get() but deleted from map. Fixed by using Lock().
2. **Deadlock in GetTTL():** Called Get() while holding lock. Fixed with getIfValid() helper.

## Learned
- RLock doesn't protect writes
- Can't call locked method from locked method
- Race detector working overview



## TODO - Future Redis Features

### High Priority
- **SAVE/LOAD:** Persist database to disk (RDB format)
  - Implement SAVE command to snapshot data
  - Load snapshot on server startup
  - Handle corruption/partial writes

### List Operations (Nice to Have)
- **LPOP/RPOP:** Remove and return elements from list ends
- **LLEN:** Return list length
- **LINDEX:** Get element at index
- **LSET:** Set element at index

### New Data Types
- **Hash:** HSET, HGET, HDEL, HGETALL
- **Set:** SADD, SMEMBERS, SISMEMBER, SREM
- **Sorted Set:** ZADD, ZRANGE, ZREM

### Protocol Improvements
- **Streaming Parser:** Handle messages > 4KB
- **Partial Read Support:** Handle TCP fragmentation
- **Pipelining:** Batch multiple commands

### Advanced Features
- **Pub/Sub:** PUBLISH, SUBSCRIBE channels
- **Transactions:** MULTI, EXEC, DISCARD
- **Replication:** Master-slave setup
- **Cluster Mode:** Distributed sharding

### Performance
- **Benchmarking:** Compare with real Redis
- **Memory Optimization:** Reduce per-key overhead
- **Lock Optimization:** Per-key locks instead of global