# Protobuf Accessor Patterns Comparison Table

This table compares accessor method patterns across Java, Go (Opaque API), C++, C#, and the current MySQL implementation for different protobuf field types.

## Simple Fields

### Non-Optional Fields (Proto3 implicit presence)

| Operation | Java | Go | C++ | C# | MySQL | Notes |
|-----------|------|----|----|----|----|-------|
| **Getter** | `getFoo()` | `GetFoo()` | `foo()` | `Foo { get; }` | ✅ `test_get_foo(proto_data)` | Always returns default value when unset: 0 (numbers), false (bool), "" (string), empty bytes, first enum value |
| **Getter (Enum Name)** | ❌ | ❌ | ❌ | ❌ | ✅ `test_get_foo__as_name(proto_data)` | Returns enum value as string name (enum fields only) |
| **Setter** | `setFoo(value)` | `SetFoo(value)` | `set_foo(value)` | `Foo { set; }` | ✅ `test_set_foo(proto_data, value)` | Sets field value |
| **Setter (From Name)** | ❌ | ❌ | ❌ | ❌ | ✅ `test_set_foo__from_name(proto_data, "ENUM_NAME")` | Sets enum value from string name (enum fields only) |
| **Clear** | `clearFoo()` | `ClearFoo()` | `clear_foo()` | ❌ (no clear for non-optional) | ✅ `test_clear_foo(proto_data)` | Reset to default |
| **Has** | ❌ | ❌ | `has_foo()` | ❌ | ❌ | C++ has presence even for proto3 |

### Optional Fields (Proto2 or Proto3 explicit presence)

| Operation | Java | Go | C++ | C# | MySQL | Notes |
|-----------|------|----|----|----|----|-------|
| **Getter** | `getFoo()` | `GetFoo()` | `foo()` | `Foo { get; }` | ✅ `test_get_foo(proto_data)` | Returns default value when unset: 0 (numbers), false (bool), "" (string), empty bytes, first enum value |
| **Getter with Default** | ❌ | ❌ | ❌ | ❌ | ✅ `test_get_foo__or(proto_data, default_value)` | Returns field value if present, otherwise returns provided default value |
| **Getter (Enum Name)** | ❌ | ❌ | ❌ | ❌ | ✅ `test_get_foo__as_name(proto_data)` | Returns enum value as string name (enum fields only) |
| **Getter (Enum Name with Default)** | ❌ | ❌ | ❌ | ❌ | ✅ `test_get_foo__as_name_or(proto_data, default_name)` | Returns enum name if present, otherwise returns provided default name (optional enum fields only) |
| **Setter** | `setFoo(value)` | `SetFoo(value)` | `set_foo(value)` | `Foo { set; }` | ✅ `test_set_foo(proto_data, value)` | Sets field value |
| **Setter (From Name)** | ❌ | ❌ | ❌ | ❌ | ✅ `test_set_foo__from_name(proto_data, "ENUM_NAME")` | Sets enum value from string name (enum fields only) |
| **Has** | `hasFoo()` | `HasFoo()` | `has_foo()` | `HasFoo { get; }` | ✅ `test_has_foo(proto_data)` | Check if field is set |
| **Clear** | `clearFoo()` | `ClearFoo()` | `clear_foo()` | `ClearFoo()` | ✅ `test_clear_foo(proto_data)` | Clear presence |

## Repeated Fields

| Operation | Java | Go | C++ | C# | MySQL | Notes |
|-----------|------|----|-----|----|-------|-------|
| **Get All** | `getFooList()` | `GetFoo()` | `foo()` | `Foo { get; }` | ✅ `test_get_all_foo(proto_data)` | Returns entire collection |
| **Set All** | ❌ (builder only) | `SetFoo([]T)` | ❌ (direct access) | ❌ (readonly prop) | ✅ `test_set_all_foo(proto_data, array)` | Replace entire collection |
| **Get Count** | `getFooCount()` | ❌ (use `len()`) | `foo_size()` | `Foo.Count` | ✅ `test_count_foo(proto_data)` | Number of elements |
| **Get Index** | `getFoo(index)` | ❌ (use indexing) | `foo(index)` | `Foo[index]` | ✅ `test_get_foo(proto_data, index)` | Get element at index |
| **Set Index** | `setFoo(index, value)` | ❌ (use indexing) | `set_foo(index, val)` | `Foo[index] = val` | ✅ `test_set_foo(proto_data, index, value)` | Set element at index |
| **Insert at Index** | ❌ | ❌ | ❌ | `Foo.Insert(index, val)` | ✅ `test_insert_foo(proto_data, index, value)` | Insert element at index |
| **Remove at Index** | ❌ | ❌ | ❌ | `Foo.RemoveAt(index)` | ✅ `test_remove_foo(proto_data, index)` | Remove element at index |
| **Add/Append** | `addFoo(value)` | `AppendFoo(vals...)` | `add_foo(value)` | `Foo.Add(val)` | ✅ `test_add_foo(proto_data, value)` | Append element(s) |
| **Add All** | `addAllFoo(collection)` | `AppendFoo(vals...)` | ❌ (manual loop) | `Foo.AddRange(vals)` | ✅ `test_add_all_foo(proto_data, array)` | Append multiple |
| **Clear** | `clearFoo()` | `ClearFoo()` | `clear_foo()` | `Foo.Clear()` | ✅ `test_clear_foo(proto_data)` | Remove all elements |
| **Mutable** | ❌ | ❌ | `mutable_foo()` | N/A (direct access) | ❌ | Get mutable container |

## Map Fields

| Operation | Java | Go | C++ | C# | MySQL | Notes |
|-----------|------|----|-----|----|-------|-------|
| **Get All** | `getFooMap()` | `GetFoo()` | `foo()` | `Foo { get; }` | ✅ `test_get_all_foo(proto_data)` | Returns entire map as JSON object or `[]` if empty |
| **Set All** | ❌ (builder only) | `SetFoo(map[K]V)` | ❌ (direct access) | ❌ (readonly prop) | ✅ `test_set_all_foo(proto_data, map)` | Replace entire map |
| **Get Count** | ❌ (use `size()`) | ❌ (use `len()`) | `foo_size()` | `Foo.Count` | ✅ `test_count_foo(proto_data)` | Number of entries |
| **Contains Key** | `containsFoo(key)` | ❌ (use map check) | ❌ (use `.find()`) | `Foo.ContainsKey(key)` | ✅ `test_contains_foo(proto_data, key)` | Check if key exists |
| **Get Value** | `getFooOrDefault(key, def)` | ❌ (use map access) | `foo().at(key)` | `Foo[key]` | ✅ `test_get_foo(proto_data, key)` | Get value by key, returns type default if missing |
| **Get Value with Default** | ❌ | ❌ | ❌ | ❌ | ✅ `test_get_foo__or(proto_data, key, default_value)` | Returns value for key if exists, otherwise returns provided default |
| **Put Entry** | `putFoo(key, value)` | ❌ (set entire map) | `(*mutable_foo())[key] = val` | `Foo[key] = val` | ✅ `test_put_foo(proto_data, key, value)` | Add/update single entry |
| **Put All** | `putAllFoo(map)` | ❌ (set entire map) | ❌ (manual loop) | `Foo.Add(kvp)` | ✅ `test_put_all_foo(proto_data, map)` | Merge multiple entries |
| **Remove Key** | `removeFoo(key)` | ❌ (set entire map) | `mutable_foo()->erase(key)` | `Foo.Remove(key)` | ✅ `test_remove_foo(proto_data, key)` | Remove entry by key |
| **Clear** | `clearFoo()` | `ClearFoo()` | `clear_foo()` | `Foo.Clear()` | ✅ `test_clear_foo(proto_data)` | Remove all entries |
| **Mutable** | ❌ | ❌ | `mutable_foo()` | N/A (direct access) | ❌ | Get mutable map |

## Oneof Fields

Consider: `oneof choice { int32 foo_int = 4; string foo_string = 9; }`

| Operation | Java | Go | C++ | C# | MySQL | Notes |
|-----------|------|----|-----|----|-------|-------|
| **Which Case** | `getChoiceCase()` | `WhichChoice()` | `choice_case()` | `ChoiceCase { get; }` | ✅ `test_which_choice(proto_data)` | Returns current case (field number in MySQL) |
| **Clear Oneof** | `clearChoice()` | `ClearChoice()` | `clear_choice()` | `ClearChoice()` | ✅ `test_clear_choice(proto_data)` | Clear entire oneof |

**Individual Field Operations**: Individual fields within a oneof have the same accessor patterns as **Optional Fields** (see above table), since oneof fields always have presence semantics. The key difference is that setting any field in a oneof automatically clears all other fields in the same oneof group.

## Message Fields

### Message Fields (All have presence in proto3)

| Operation | Java | Go | C++ | C# | MySQL | Notes |
|-----------|------|----|-----|----|-------|-------|
| **Getter** | `getPerson()` | `GetPerson()` | `person()` | `Person { get; }` | ✅ `test_get_person(proto_data)` | Returns language-specific default: default instance (Java/C++/C#), nil (Go), empty JSON object (MySQL) |
| **Getter with Default** | ❌ | ❌ | ❌ | ❌ | ✅ `test_get_person__or(proto_data, default_message)` | Returns message field if present, otherwise returns provided default message |
| **Setter** | `setPerson(value)` | `SetPerson(value)` | `set_allocated_person()` | `Person { set; }` | ✅ `test_set_person(proto_data, message)` | Set message field |
| **Has** | `hasPerson()` | `HasPerson()` | `has_person()` | `HasPerson { get; }` | ✅ `test_has_person(proto_data)` | Check if field is set |
| **Clear** | `clearPerson()` | `ClearPerson()` | `clear_person()` | `Person = null` | ✅ `test_clear_person(proto_data)` | Clear presence |
| **Merge** | `mergePerson(value)` | ❌ | `mutable_person()->MergeFrom()` | ❌ | ❌ `test_merge_person(proto_data, message)` | Merge with existing |
| **Mutable** | ❌ | ❌ | `mutable_person()` | Direct access | ❌ | Get/create mutable ref |
| **Builder** | `getPersonBuilder()` | ❌ | ❌ | ❌ | ❌ | Get nested builder |
| **Release** | ❌ | ❌ | `release_person()` | ❌ | ❌ | Transfer ownership |

**Note**: In proto3, message fields always have field presence semantics, whether explicitly marked `optional` or not. This is different from scalar fields where presence is only available for explicit `optional` fields.

## Default Value Behavior for Unset Fields

Understanding getter behavior when fields are not present is crucial for protobuf usage:

### Key Principle: Getters Always Return a Value
**Protobuf getters never throw exceptions or return null/undefined for scalar fields**. They always return a type-appropriate default value when the field is unset.

### Default Values by Field Type

| Field Type | Default Value | Examples |
|------------|---------------|----------|
| **Numeric** (int32, int64, uint32, uint64, sint32, sint64, fixed32, fixed64, sfixed32, sfixed64, float, double) | `0` or `0.0` | `0`, `0.0` |
| **Boolean** | `false` | `false` |
| **String** | Empty string | `""` |
| **Bytes** | Empty bytes | `[]` or `""` (language-dependent) |
| **Enum** | First enum value | Usually `0` (e.g., `STATUS_UNSPECIFIED = 0`) |
| **Message** | Language-dependent | See table below |

### Message Field Default Behavior

Message fields behave differently across languages when unset:

| Language | Unset Message Field Returns | Notes |
|----------|----------------------------|-------|
| **Java** | Default instance (never `null`) | Always returns empty message instance; **CORRECTED** from previous incorrect documentation |
| **Go** | `nil` | Safe to check with `if msg.GetPerson() != nil` |
| **C++** | Default instance | Returns const reference to default instance |
| **Python** | Default instance (never `None`) | Always returns empty message instance; `WhichOneof()` returns `None` for unset oneof |
| **C#** | Default instance (never `null`) | Returns empty message with no fields set; cannot be null |
| **MySQL** | `{}` (empty JSON object) | Consistent with JSON representation |

### Presence vs Value

**Critical Distinction**: The value returned by a getter is separate from whether the field was explicitly set:

```java
// Java example
Person person = Person.newBuilder().build(); // Empty message

// These both return default values:
String name = person.getName();        // Returns ""
int age = person.getAge();            // Returns 0

// But presence checks show they weren't set:
boolean hasName = person.hasName();    // Returns false (if optional)
boolean hasAge = person.hasAge();      // Returns false (if optional)
```

### Proto2 vs Proto3 Differences

| Proto Version | Optional Fields | Default Behavior |
|---------------|-----------------|------------------|
| **Proto2** | All fields optional by default | Has presence tracking for all fields |
| **Proto3** | Fields without `optional` keyword | No presence tracking for scalars |
| **Proto3** | Fields with `optional` keyword | Has presence tracking like proto2 |

### Best Practices

1. **Always use `has*()` methods** to check if optional fields are set before using their values
2. **Don't rely on default values** to determine if a field was set - use presence methods
3. **For message fields**, only Go returns `nil` - all other languages return safe default instances
4. **Consider using wrapper types** (e.g., `google.protobuf.StringValue`) when you need to distinguish between unset and default values

### Python-Specific Behavior

Python protobuf getter behavior for unset message fields has evolved over time:

#### **Message Fields (Definitive Behavior)**
- **Never returns `None`**: Message field getters always return a default instance, never `None`
- **Same behavior across proto2 and proto3**: Both return empty message instances for unset fields
- **Read-through access**: Reading nested fields doesn't automatically set the parent field
- **Implicit setting**: Only setting (not reading) nested fields marks the parent message as present

#### **Example Behavior (Tested)**
```python
# Proto2 and Proto3 behavior is identical for message fields
person = Person()  # Any protobuf message
assert not person.HasField("address")  # False - field not set
address = person.address               # Returns empty Address() instance
assert address is not None             # Always True - never None
assert person.address.street == ""     # Default string value
assert not person.HasField("address")  # Still False - reading didn't set it

person.address.street = "Main St"      # Setting nested field sets parent
assert person.HasField("address")      # Now True - field was set
```

#### **Oneof Fields - The Only `None` Case**
```python
msg = TestOneof()
assert msg.WhichOneof("test_oneof") is None  # Only case where None is returned
msg.name = "test"
assert msg.WhichOneof("test_oneof") == "name"  # Now returns field name
```

#### **Key Takeaways**
- **Message fields**: Always return empty instances, never `None`
- **Consistent behavior**: Proto2 and proto3 behave identically for message fields
- **Only `None` case**: `WhichOneof()` method for unset oneof fields
- **Presence vs. access**: Reading doesn't set fields, only writing does

### Java-Specific Behavior

Java protobuf getter behavior for unset message fields was tested and documented:

#### **Message Fields (Tested Behavior)**
- **Never returns `null`**: Message field getters always return a default instance, never `null`
- **Same behavior across proto2 and proto3**: Both return empty message instances for unset fields
- **Includes explicit optional**: Even proto3 `optional` message fields return default instances
- **Safe access**: No null pointer exceptions when accessing nested fields of unset message fields

#### **Example Behavior (Tested)**
```java
// Proto2, Proto3, and Proto3 Optional behavior is identical for message fields
Test.Person person2 = Test.Person.newBuilder().build();
assert !person2.hasAddress();                    // false - field not set
Test.Address address = person2.getAddress();     // Returns default Address instance
assert address != null;                          // Always true - never null
assert address.getStreet().equals("");           // Default string value
// Safe to access nested fields without null checks

// Setting nested fields marks parent as present
person2.getAddress().toBuilder().setStreet("Main St").build();
// Note: In Java, you typically use builders for modifications
```

#### **Oneof Fields**
```java
Oneof.TestOneof msg = Oneof.TestOneof.newBuilder().build();
assert msg.getTestOneofCase() == Oneof.TestOneof.TestOneofCase.TESTONEOF_NOT_SET;
// Returns enum value, not null

Oneof.TestOneof msgWithName = Oneof.TestOneof.newBuilder().setName("test").build();
assert msgWithName.getTestOneofCase() == Oneof.TestOneof.TestOneofCase.NAME;
```

#### **Key Takeaways**
- **Message fields**: Always return empty instances, never `null` - contrary to common misconceptions
- **Null safety**: Java protobuf provides null-safe access to message fields
- **Consistent with other languages**: Java aligns with C++, C#, and Python behavior
- **Only Go is different**: Go is the sole major language that returns `nil` for unset message fields

### C#-Specific Behavior (Google.Protobuf)

C# protobuf behavior has some unique characteristics:

#### **Message Fields**
- **Never return `null`**: Like Java/C++/Python, C# message field getters never return `null`
- **Always return default instance**: Unset message fields return an empty message with no fields set
- **Can be set to `null`**: Setting a message field to `null` clears it (equivalent to calling `Clear()`)
- **Default instance checking**: You can check `if (msg.Person == PersonType.DefaultInstance)` to detect unset fields

#### **Scalar Fields**
- **String fields**: Return empty string `""` when unset; **throw `ArgumentNullException` if you try to set them to `null`**
- **Numeric fields**: Return `0` when unset
- **Bool fields**: Return `false` when unset
- **Bytes fields**: Return empty `ByteString` when unset; **throw exception if set to `null`**

#### **Wrapper Types for Nullability**
When you need nullable semantics, use wrapper types:
```protobuf
import "google/protobuf/wrappers.proto";
message Example {
    google.protobuf.Int32Value optional_int = 1;     // Maps to Nullable<int>
    google.protobuf.StringValue optional_string = 2; // Maps to string (can be null)
}
```

#### **Key Takeaways**
- **Message fields**: Safe from `null` - always get a valid object, but check for default instance
- **Scalar fields**: Never `null`, but don't try to set them to `null` (will throw exception)
- **Null safety**: C# protobuf is designed to avoid `null` reference exceptions for message fields
- **Presence detection**: Use `has*()` methods (proto2) or wrapper types (proto3) when you need to distinguish unset vs default values

### MySQL-Specific Behavior

The MySQL protobuf implementation follows these principles:
- **Scalar fields**: Return type-appropriate defaults (0, false, "", empty bytes, first enum value)
- **Message fields**: Return `{}` (empty JSON object) representing a default/empty message
- **Presence tracking**: Available through `has_*()` functions for optional fields
- **Consistency**: All getter functions are deterministic and always return valid JSON

## Language Behavior Summary

**Tested Languages**: We directly tested **Java and Python** protobuf behavior for unset message fields:

| Language | Tested Behavior | Result |
|----------|-----------------|--------|
| **Java** | ✅ Tested | Returns default instance (never `null`) |
| **Python** | ✅ Tested | Returns default instance (never `None`) |
| **Go** | 📚 Documented | Returns `nil` (from Go documentation) |
| **C++** | 📚 Documented | Returns default instance (from documentation) |
| **C#** | 📚 Documented | Returns default instance (from documentation) |
| **MySQL** | 🏠 Implementation | Returns `{}` (empty JSON object) |

**Key Finding from Testing**: Both Java and Python return default instances rather than null-like values, correcting previous misconceptions about Java behavior.

**MySQL Design Validation**: The tested languages (Java, Python) align with MySQL's `{}` approach, providing:
- **Type safety**: No null reference exceptions
- **Predictable behavior**: Always returns valid objects
- **Consistency**: Matches the pattern we verified in major languages

**Conclusion**: Based on our testing of Java and Python, MySQL's design choice to return `{}` for unset message fields aligns well with tested protobuf implementations and provides good type safety characteristics.

## References

- [Protocol Buffers Design Decision: No Nullable Setters/Getters Support](https://protobuf.dev/design-decisions/nullable-getters-setters/) - Official documentation explaining why protobuf doesn't support nullable getters and setters by design

## Method Naming Conventions

| Language | Field Name | Getter | Setter | Has | Clear | Notes |
|----------|------------|--------|--------|-----|-------|-------|
| **Java** | `birth_year` | `getBirthYear()` | `setBirthYear()` | `hasBirthYear()` | `clearBirthYear()` | CamelCase |
| **Go** | `birth_year` | `GetBirthYear()` | `SetBirthYear()` | `HasBirthYear()` | `ClearBirthYear()` | PascalCase |
| **MySQL** | `birth_year` | `test_get_birth_year()` | `test_set_birth_year()` | `test_has_birth_year()` | `test_clear_birth_year()` | snake_case |
| **MySQL (Current)** | `birth_year` | `test_get_birth_year()` | `test_set_birth_year()` | `test_has_birth_year()` | `test_clear_birth_year()` | snake_case with configurable prefix |

## Key Differences by Language

### Java
- **Builder Pattern**: All mutations through builder classes
- **Rich Collections**: Specialized methods for repeated/map operations
- **Type Safety**: Strong generic typing
- **Immutable Messages**: Thread-safe by design

### Go (Opaque API)
- **Direct Mutation**: Methods modify message in-place
- **Minimal Collections**: Basic get/set/clear for collections
- **Performance Focus**: Reduced allocator pressure
- **Memory Efficient**: Bit-field presence tracking

### MySQL (Recommended)
- **Functional Style**: Pure functions returning new JSON state
- **SQL Integration**: Designed for SQL function interface
- **JSON-Based**: Internal representation as JSON objects
- **Comprehensive**: Full method set for all operations

#### MySQL Implementation Status
- ✅ **Configurable Prefixing**: Functions use `${package}_${type}_` prefix pattern (configurable via prefix_map)
- ✅ **Complete Accessor API**: Full get/set/has/clear support for all field types
- ✅ **Proto3 Semantics**: Proper default value omission and presence handling
- ✅ **Advanced Features**: OneOf mutual exclusion, enum conversions, repeated field operations
- ✅ **Performance Optimized**: Number-JSON format, deterministic functions, efficient field number indexing
- ✅ **Production Ready**: Input validation, 64-character name limit handling with helpful errors
- ✅ **Extended Variants**: Supports `__modifier` pattern for specialized accessors

## Implementation Priority

For MySQL protobuf opaque API implementation, recommended priority:

1. ✅ **High Priority**: Get, Set, Clear operations for all field types
2. ✅ **Medium Priority**: Has operations for optional/oneof fields, Which operations for oneofs
3. ✅ **Low Priority**: Individual map operations (put/remove), merge operations
4. 🔄 **Future**: Index-based repeated field access (`get/set` with index parameters), ❌ builder patterns

**Current Status**: The MySQL protobuf implementation has achieved complete coverage of High, Medium, and Low priority features with sophisticated additional functionality including:
- ✅ **Complete Map API**: All CRUD operations (`get`, `put`, `remove`, `contains`, `count`, `put_all`) 
- ✅ **Advanced Accessors**: Enum handling with string name getters (`__as_name`, `__as_name_or`) and setters (`__from_name`)
- ✅ **Type Safety**: Proper MySQL type mapping with binary format support for floats/doubles
- ✅ **Performance Optimizations**: Efficient JSON path operations and deterministic functions

## MySQL Extended Variants

The MySQL implementation supports extended function variants using the `__modifier` pattern for specialized accessor operations:

### Enum Field Variants

| Field Type | Base Function | Extended Variant | Return Type | Description |
|------------|---------------|------------------|-------------|-------------|
| **Enum** | `test_get_status(proto_data)` → `1` | `test_get_status__as_name(proto_data)` → `"STATUS_ACTIVE"` | `LONGTEXT` | Returns enum value as string name ✅ **IMPLEMENTED** |
| **Optional Enum** | `test_get_status__or(proto_data, 999)` → `999` | `test_get_status__as_name_or(proto_data, "DEFAULT")` → `"DEFAULT"` | `LONGTEXT` | Returns enum name if present, else custom default ✅ **IMPLEMENTED** |
| **Enum** | `test_set_status(proto_data, 1)` | `test_set_status__from_name(proto_data, "STATUS_ACTIVE")` | `JSON` | Sets enum value from string name ✅ **IMPLEMENTED** |
| **Repeated Enum (All)** | `test_get_all_statuses(proto_data)` → `[1, 2]` | `test_get_all_statuses__as_names(proto_data)` → `["STATUS_ACTIVE", "STATUS_INACTIVE"]` | `JSON` | Returns array of enum names |
| **Repeated Enum (All)** | `test_set_all_statuses(proto_data, [1, 2])` | `test_set_all_statuses__from_names(proto_data, '["STATUS_ACTIVE", "STATUS_INACTIVE"]')` | `JSON` | Sets repeated enum from array of string names |
| **Repeated Enum (Index)** | `test_get_statuses(proto_data, index)` → `1` | `test_get_statuses__as_name(proto_data, index)` → `"STATUS_ACTIVE"` | `VARCHAR` | Returns enum at index as string name |
| **Repeated Enum (Index)** | `test_set_statuses(proto_data, index, 1)` | `test_set_statuses__from_name(proto_data, index, "STATUS_ACTIVE")` | `JSON` | Sets enum at index from string name |
| **Repeated Enum (Add)** | `test_add_statuses(proto_data, 1)` | `test_add_statuses__from_name(proto_data, "STATUS_ACTIVE")` | `JSON` | Adds enum element from string name |
| **Repeated Enum (Insert)** | `test_insert_statuses(proto_data, index, 1)` | `test_insert_statuses__from_name(proto_data, index, "STATUS_ACTIVE")` | `JSON` | Inserts enum at index from string name |

### Additional Index Operations

Beyond the standard get/set at index, MySQL supports additional index-based operations:

| Operation | Function Pattern | Parameters | Description |
|-----------|------------------|------------|-------------|
| **Get First** | `test_get_first_foo(proto_data)` | `proto_data JSON` | Get first element |
| **Get Last** | `test_get_last_foo(proto_data)` | `proto_data JSON` | Get last element |
| **Set First** | `test_set_first_foo(proto_data, value)` | `proto_data JSON, value T` | Set first element |
| **Set Last** | `test_set_last_foo(proto_data, value)` | `proto_data JSON, value T` | Set last element |

### Format Modifiers for Index Operations

The `__modifier` pattern can be applied to index operations for format conversion:

```sql
-- Get enum at specific index as string name
test_get_statuses__as_name(proto_data, index) → "STATUS_ACTIVE"

-- Set enum at specific index from string name
test_set_statuses__from_name(proto_data, index, "STATUS_ACTIVE")

-- Get first enum element as string name  
test_get_first_statuses__as_name(proto_data) → "STATUS_ACTIVE"

-- Get last enum element as string name
test_get_last_statuses__as_name(proto_data) → "STATUS_INACTIVE"

-- Insert enum at index from string name
test_insert_statuses__from_name(proto_data, index, "STATUS_PENDING")

-- Add enum from string name (already covered above)
test_add_statuses__from_name(proto_data, "STATUS_ACTIVE")
```

### Other Potential Variants

| Field Type | Base Function | Extended Variant | Description |
|------------|---------------|------------------|-------------|
| **String** | `test_get_name(proto_data)` | `test_get_name__as_base64(proto_data)` | Return string as base64 |
| **Bytes** | `test_get_data(proto_data)` | `test_get_data__as_hex(proto_data)` | Return bytes as hex string |
| **Timestamp** | `test_get_created_at(proto_data)` | `test_get_created_at__as_rfc3339(proto_data)` | Format timestamp as RFC 3339 |

### Benefits of `__modifier` Pattern

1. **Conflict-free**: Double underscore prevents field name conflicts
2. **Composable**: Multiple modifiers can be combined (`__at_index__as_name`)
3. **Self-documenting**: Clear indication of what the variant does
4. **Extensible**: Easy to add new modifiers without breaking existing functions
5. **Systematic**: Consistent pattern across all field types

This table provides a comprehensive reference for implementing consistent protobuf accessor patterns across different target languages.