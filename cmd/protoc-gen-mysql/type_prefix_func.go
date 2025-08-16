package main

import (
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// createTypePrefixFunc creates a function that maps package and type names to prefixes
func createTypePrefixFunc(packagePrefixMap map[string]string) protocgenmysql.TypePrefixFunc {
	return func(packageName protoreflect.FullName, fullTypeName protoreflect.FullName) string {
		packageNameStr := string(packageName)
		fullTypeNameStr := string(fullTypeName)

		// First check if there's a specific mapping for this package
		if prefix, exists := packagePrefixMap[packageNameStr]; exists {
			// Extract just the type name (last part after last dot)
			typeParts := strings.Split(fullTypeNameStr, ".")
			typeName := typeParts[len(typeParts)-1]
			return prefix + strings.ToLower(typeName)
		}

		// Default behavior: use package-based prefix
		if packageNameStr == "" {
			// Extract just the type name for no package
			typeParts := strings.Split(fullTypeNameStr, ".")
			typeName := typeParts[len(typeParts)-1]
			return "pb_" + strings.ToLower(typeName)
		}

		// Use package name as prefix with underscore, plus lowercase type name
		packagePrefix := strings.ReplaceAll(packageNameStr, ".", "_") + "_"
		typeParts := strings.Split(fullTypeNameStr, ".")
		typeName := typeParts[len(typeParts)-1]
		return packagePrefix + strings.ToLower(typeName)
	}
}
