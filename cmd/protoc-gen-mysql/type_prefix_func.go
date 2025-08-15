package main

import (
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
)

// createTypePrefixFunc creates a function that maps package and type names to prefixes
func createTypePrefixFunc(packagePrefixMap map[string]string) protocgenmysql.TypePrefixFunc {
	return func(packageName string, fullTypeName string) string {
		// First check if there's a specific mapping for this package
		if prefix, exists := packagePrefixMap[packageName]; exists {
			// Extract just the type name (last part after last dot)
			typeParts := strings.Split(fullTypeName, ".")
			typeName := typeParts[len(typeParts)-1]
			return prefix + strings.ToLower(typeName)
		}

		// Default behavior: use package-based prefix
		if packageName == "" {
			// Extract just the type name for no package
			typeParts := strings.Split(fullTypeName, ".")
			typeName := typeParts[len(typeParts)-1]
			return "pb_" + strings.ToLower(typeName)
		}

		// Use package name as prefix with underscore, plus lowercase type name
		packagePrefix := strings.ReplaceAll(packageName, ".", "_") + "_"
		typeParts := strings.Split(fullTypeName, ".")
		typeName := typeParts[len(typeParts)-1]
		return packagePrefix + strings.ToLower(typeName)
	}
}
