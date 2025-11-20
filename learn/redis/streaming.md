# STREAMING.md

## Known Limitation
Parser assumes entire message fits in 64KB buffer.

## Fails When
- Single value > 64KB
- Multiple commands in one TCP packet

## Fix Required
Rewrite parser to use bufio.Reader + io.ReadFull()
See: https://redis.io/docs/reference/protocol-spec/