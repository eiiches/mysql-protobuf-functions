# Installation Guide

## Requirements

- **MySQL**: 8.0.17 or later
  - JSON_TABLE() was added in 8.0.4 but requires 8.0.17 for [this critical bugfix](https://bugs.mysql.com/bug.php?id=92976)
- **Aurora MySQL**: 3.04.0 (oldest available 3.x version as of June 2025) or later

## Quick Start

1. **Clone the repository** to get the SQL files:
   ```bash
   git clone https://github.com/eiiches/mysql-protobuf-functions.git
   cd mysql-protobuf-functions
   ```

2. **Install core functions**:
   ```bash
   mysql -u your_username -p your_database < build/protobuf.sql
   ```

3. **Install optional components**
   ```bash
   # For JSON conversion support.
   mysql -u your_username -p your_database < build/protobuf-json.sql  # depends on protobuf.sql

   # For dynamic schema loading from FileDescriptorSet binary.
   mysql -u your_username -p your_database < build/protobuf-descriptor.sql  # depends on protobuf.sql and protobuf-json.sql
   ```

## Component Dependencies

The installation components have the following dependency chain:

- `protobuf.sql` - Core wire format parsing (required)
- `protobuf-json.sql` - JSON conversion (depends on protobuf.sql)
- `protobuf-descriptor.sql` - Schema loading (depends on protobuf.sql and protobuf-json.sql)

## Important Installation Notes

- ⚠️ All functions and procedures use `_pb_` or `pb_` prefixes to avoid naming conflicts
- 📝 Verify existing routines before installation to prevent overwrites
- 🔍 Check that your MySQL version meets the minimum requirements before installation
- 🛠️ The functions are implemented entirely in MySQL stored procedures - no native libraries required

## Verification

After installation, verify the functions are available:

```sql
-- Check core functions
SELECT pb_message_new() IS NOT NULL AS core_installed;

-- Check JSON functions (if installed)
SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES 
WHERE ROUTINE_SCHEMA = DATABASE() 
AND ROUTINE_NAME LIKE 'pb_%json%';

-- Check descriptor functions (if installed)  
SELECT ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES 
WHERE ROUTINE_SCHEMA = DATABASE() 
AND ROUTINE_NAME LIKE 'pb_%descriptor%';
```