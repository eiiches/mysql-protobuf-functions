package main

import (
	"fmt"
	"log"

	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	"google.golang.org/protobuf/types/descriptorpb"
)

func main() {
	// Build a FileDescriptorSet with dependencies
	fileDescriptorSet := protoreflectutils.BuildFileDescriptorSetWithDependencies(
		descriptorpb.File_google_protobuf_descriptor_proto,
	)

	// Convert to JSON using the new package
	jsonStr, err := descriptorsetjson.ToJson(fileDescriptorSet)
	if err != nil {
		log.Fatalf("Failed to convert descriptor.proto to JSON: %v", err)
	}

	// Output CREATE FUNCTION statement
	fmt.Printf(`DROP FUNCTION IF EXISTS _pb_get_descriptor_proto_set $$
CREATE FUNCTION _pb_get_descriptor_proto_set() RETURNS JSON DETERMINISTIC
BEGIN
	RETURN CAST('%s' AS JSON);
END $$
`, escapeSQLString(jsonStr))
}

func escapeSQLString(s string) string {
	// Escape single quotes and backslashes for SQL string literals
	result := ""
	for _, char := range s {
		switch char {
		case '\'':
			result += "''"
		case '\\':
			result += "\\\\"
		default:
			result += string(char)
		}
	}
	return result
}
