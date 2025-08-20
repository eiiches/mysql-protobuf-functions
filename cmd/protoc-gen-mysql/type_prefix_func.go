package main

import (
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// createTypePrefixFunc creates a function that maps package and type names to prefixes
func createTypePrefixFunc(packagePrefixMap map[protoreflect.FullName]string) protocgenmysql.TypePrefixFunc {
	return func(packageName protoreflect.FullName, fullTypeName protoreflect.FullName) string {
		packageNameStr := string(packageName)
		fullTypeNameStr := string(fullTypeName)

		// Extract the type name part after the package name
		var typeName string
		if packageNameStr != "" && strings.HasPrefix(fullTypeNameStr, packageNameStr+".") {
			// Remove package name and leading dot
			typeName = fullTypeNameStr[len(packageNameStr)+1:]
		} else {
			// Use full type name if no package or package doesn't match
			typeName = fullTypeNameStr
		}

		// Look up package prefix recursively
		matchedPrefix, remainingPackage := findPackagePrefix(packageName, packagePrefixMap)

		// Build the final name
		var finalName string
		if remainingPackage != "" {
			// Include remaining package parts: ${packagePrefix}${remainingPackage}_${type_part}
			finalName = matchedPrefix + strings.ReplaceAll(strings.ToLower(string(remainingPackage)), ".", "_") + "_" + toSnakeTypeName(typeName)
		} else {
			// Direct match: ${packagePrefix}${type_part}
			finalName = matchedPrefix + toSnakeTypeName(typeName)
		}

		return finalName
	}
}

// findPackagePrefix looks up package prefix recursively
// Returns the matched prefix and any remaining package parts
func findPackagePrefix(packageName protoreflect.FullName, packagePrefixMap map[protoreflect.FullName]string) (string, protoreflect.FullName) {
	currentPackage := packageName
	for {
		if prefix, exists := packagePrefixMap[currentPackage]; exists {
			// Found a match, calculate remaining package parts
			if string(currentPackage) == "" {
				return prefix, packageName
			} else if packageName == currentPackage {
				return prefix, "" // Exact match, no remaining parts
			} else {
				return prefix, protoreflect.FullName(strings.TrimPrefix(string(packageName), string(currentPackage)+"."))
			}
		}
		if currentPackage == "" {
			break
		}
		currentPackage = currentPackage.Parent()
	}

	// No match found, use default behavior
	return "", packageName
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
