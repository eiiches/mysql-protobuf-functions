# MySQL Protocol Buffers Functions

[![MySQL Version](https://img.shields.io/badge/MySQL-8.0.17%2B-blue)](https://dev.mysql.com/downloads/mysql/)
[![Aurora MySQL](https://img.shields.io/badge/Aurora%20MySQL-3.04.0%2B-orange)](https://aws.amazon.com/rds/aurora/)

A comprehensive library of MySQL stored functions and procedures for working with Protocol Buffers (protobuf) encoded data directly within MySQL databases. This project enables you to parse, query, and manipulate protobuf messages without requiring external applications or services.

## Features

- 🔍 **Field Access**: Extract specific fields from protobuf messages using field numbers
- ✏️ **Message Manipulation**: Create, modify, and update protobuf messages directly in MySQL - set fields, add/remove repeated elements, and clear fields
- 🔄 **JSON Conversion**: Convert protobuf messages to JSON format for easier debugging
- 🛠️ **Pure MySQL Implementation**: Written entirely in MySQL stored functions and procedures - no native libraries or external dependencies required

## Quick Start

1. **Install core functions**:

   Download [protobuf.sql](https://raw.githubusercontent.com/eiiches/mysql-protobuf-functions/refs/heads/main/build/protobuf.sql) and load it into your MySQL database:
   ```bash
   curl -fLO https://raw.githubusercontent.com/eiiches/mysql-protobuf-functions/refs/heads/main/build/protobuf.sql
   mysql -u your_username -p your_database < protobuf.sql
   ```

2. **Try it out**:
   ```sql
   -- Create new protobuf message
   SELECT pb_message_set_string_field(pb_message_new(), 1, 'Hello World');
   -- Result: _binary X'0A0B48656C6C6F20576F726C64'

   -- Extract field from protobuf message
   SELECT pb_message_get_string_field(_binary X'0A0B48656C6C6F20576F726C64', 1, '');
   -- Result: "Hello World"

   -- Convert to JSON
   -- This requires protobuf-json.sql and schema. See docs/tutorial.md for ways to load schema into MySQL.
   SELECT pb_message_to_json(greeting_schema(), '.Greeting', _binary X'0A0B48656C6C6F20576F726C64');
   -- Result: {
   --   "message": "Hello World",
   -- }
   ```

## Documentation

> **Work in Progress:** The documentation was written with the help of AI and is currently under review.

- **[Installation Guide](docs/installation.md)** - Detailed installation instructions and requirements
- **[Tutorial](docs/tutorial.md)** - Complete tutorial with examples using a sample schema
- **[Advanced Usage](docs/advanced-usage.md)** - Indexing, triggers, and performance optimization
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions
- **[API Reference](docs/function-reference.md)** - Complete function reference
- **[Roadmap](docs/roadmap.md)** - Planned features and current limitations

## Requirements

- **MySQL**: 8.0.17 or later
- **Aurora MySQL**: 3.04.0 or later

See the [Installation Guide](docs/installation.md) for detailed requirements and setup instructions.

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
