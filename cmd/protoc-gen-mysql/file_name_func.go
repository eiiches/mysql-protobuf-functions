package main

import (
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
)

// File naming functions

// flattenFileNameFunc converts "path/to/file.proto" to "path_to_file_methods.sql"
func flattenFileNameFunc(protoPath string) string {
	if protoPath == "" {
		return "" // Individual file mode, not single file
	}
	filename := strings.ReplaceAll(protoPath, "/", "_")
	filename = strings.TrimSuffix(filename, ".proto") + ".sql"
	return filename
}

// preserveFileNameFunc converts "path/to/file.proto" to "path/to/file_methods.sql"
func preserveFileNameFunc(protoPath string) string {
	if protoPath == "" {
		return "" // Individual file mode, not single file
	}
	filename := strings.TrimSuffix(protoPath, ".proto") + ".sql"
	return filename
}

// singleFileNameFunc returns "protobuf_methods.sql" for all files (single file mode)
func singleFileNameFunc(protoPath string) string {
	if protoPath == "" {
		return "protobuf_methods.sql" // Single file mode
	}
	return "" // Skip individual files in single mode
}

// getFileNameFunc returns the appropriate file naming function
func getFileNameFunc(strategy string) protocgenmysql.FileNameFunc {
	switch strategy {
	case "preserve":
		return preserveFileNameFunc
	case "single":
		return singleFileNameFunc
	case "flatten":
		fallthrough
	default:
		return flattenFileNameFunc
	}
}
