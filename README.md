[Experimental &amp; Incomplete] Pure-MySQL Protocol Buffer Functions
====

Disclaimer: This is just a toy project - very incomplete, likely abandoned, and not intended for production use.

Usage
----

1. Load functions to MySQL.

   ```console
   $ mysql < protobuf.sql
   ```

2. Use functions in statements.

   ```sql
   -- Get the value of a string field.
   mysql> SELECT pb_message_get_string_field(_binary X'100a2a03616263', 5 /* field_number */, NULL /* repeated_index */);
   +-----------------------------------------------------------------+
   | pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL) |
   +-----------------------------------------------------------------+
   | abc                                                             |
   +-----------------------------------------------------------------+
   1 row in set (0.00 sec)

   -- Get the first element of a repeated int32 (packed) field.
   mysql> SELECT pb_message_get_int32_field(_binary X'3a03010203', 7 /* field_number */, 0 /* repeated_index */);
   +---------------------------------------------------------+
   | pb_message_get_int32_field(_binary X'3a03010203', 7, 0) |
   +---------------------------------------------------------+
   |                                                       1 |
   +---------------------------------------------------------+
   1 row in set (0.01 sec)

   -- Get int32 field from a nested message field.
   mysql> SELECT pb_message_get_int32_field(pb_message_get_message_field(_binary X'4202080a', 8, NULL), 1, NULL);
   +-------------------------------------------------------------------------------------------------+
   | pb_message_get_int32_field(pb_message_get_message_field(_binary X'4202080a', 8, NULL), 1, NULL) |
   +-------------------------------------------------------------------------------------------------+
   |                                                                                              10 |
   +-------------------------------------------------------------------------------------------------+
   ```

   Protobuf schema used in this example:

   ```protobuf
   syntax = "proto3";

   message Test {
       string string_field = 5;
       repeated int32 repeated_int32_field = 7;
       TestMessage message_field = 8;
   }

   message TestMessage {
       int32 int32_field = 1;
   }
   ```

Project Status
------

* Implementation is very incomplete. Implemented field types are: (repeated) int32, int64, string, bytes, Message. Large (>=2**31) or negative integers are not tested.
  * No map support.
* Currently, only getters are implemented. Setters are technically possible, I just didn't have time.

Limitation
----------

* I never benchmarked the functions. I have no idea how slow these functions are.
  * MySQL doesn't have an ARRAY type. With the current API, each element of a repeated field must be retrieved one by one. Yes, O(n^2) to retrieve all.
* Currently, MySQL doesn't allow using stored functions in functional indexes or generated columns. That very much limits the usefulness of the functions.
