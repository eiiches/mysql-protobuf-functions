# Roadmap and Limitations

## Planned Features

- [ ] Add `pb_{message,wire_json}_which_oneof` function
- [ ] Add `pb_{message,wire_json}_get_map_entry_by_{type}_key(message, field_number, key, default_value)` function (finds the last one)
- [ ] Add `pb_{message,wire_json}_search_repeated_message_field_by_{type}_key(message, field_number, key_field_number, key, default_value)` function (finds the first one)
- [ ] JSON to Protobuf Conversion
- [x] Protobuf to JSON Conversion
  - [ ] **[Editions](https://protobuf.dev/editions/overview/) Support**

## Current Limitations

- [Groups](https://protobuf.dev/programming-guides/encoding/#groups) are not supported.