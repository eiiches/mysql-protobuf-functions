# Troubleshooting

## Known Issues

### MySQL Stored Program Cache Bug

**Issue**: When many stored functions are used in a single connection and the `stored_program_cache` limit is reached, MySQL exhibits unpredictable behavior:

- **MySQL 9.3.0**: Functions silently return `NULL` instead of the expected result
- **MySQL 8.0.x**: Functions fail with `Function does not exist` error

**Root Cause**: [MySQL Bug #95825](https://bugs.mysql.com/bug.php?id=95825)

**Workaround**: Increase the stored program cache size:
```sql
SET GLOBAL stored_program_cache = 1024;  -- Default is 256
```

**Impact**:
- Most applications won't encounter this issue as the default cache size (256) is sufficient for typical usage
- This primarily affects comprehensive test suites or applications using many different protobuf functions in a single connection
- Related discussion: [Percona Forums](https://forums.percona.com/t/intermittent-stored-function-does-not-exist-problem/5143)