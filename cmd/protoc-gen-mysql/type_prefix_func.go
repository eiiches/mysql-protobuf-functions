package main

import (
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// createTypePrefixFunc creates a function that maps package and type names to prefixes
func createTypePrefixFunc(prefixMap map[protoreflect.FullName]string) protocgenmysql.TypePrefixFunc {
	return func(packageName protoreflect.FullName, fullTypeName protoreflect.FullName) string {
		// Use the same recursive lookup logic on the full type name
		// This allows matching both package names and specific type names
		matchedPrefix, remainingName := findPrefix(fullTypeName, prefixMap)

		// Build the final name
		if remainingName != "" {
			// Include remaining parts: ${prefix}${remaining_parts}
			return matchedPrefix + toSnakeTypeName(string(remainingName))
		} else {
			// Exact match, just use the prefix
			return matchedPrefix
		}
	}
}

// findPrefix looks up prefix recursively for both packages and types
// Returns the matched prefix and any remaining name parts
func findPrefix(name protoreflect.FullName, prefixMap map[protoreflect.FullName]string) (string, protoreflect.FullName) {
	currentName := name
	for {
		if prefix, exists := prefixMap[currentName]; exists {
			// Found a match, calculate remaining name parts
			if string(currentName) == "" {
				return prefix, name
			} else if name == currentName {
				return prefix, "" // Exact match, no remaining parts
			} else {
				return prefix, protoreflect.FullName(strings.TrimPrefix(string(name), string(currentName)+"."))
			}
		}
		if currentName == "" {
			break
		}
		currentName = currentName.Parent()
	}

	// No match found, use default behavior
	return "", name
}

func toSnakeTypeName(s string) string {
	// Split by dots to handle nested types like "FieldDescriptorProto.Label"
	parts := strings.Split(s, ".")
	var snakeParts []string
	for _, part := range parts {
		snakeParts = append(snakeParts, toSnakeCase(part))
	}
	return strings.Join(snakeParts, "_")
}

// toSnakeCase converts a type name like "FieldDescriptorProto" to "field_descriptor_proto"
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
