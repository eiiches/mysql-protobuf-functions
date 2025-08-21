package main

import (
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
)

// File naming functions

// flattenFileNameFunc converts "path/to/file.proto" to "path_to_file_methods.pb.sql"
func flattenFileNameFunc(protoPath string) string {
	filename := strings.ReplaceAll(protoPath, "/", "_")
	filename = strings.TrimSuffix(filename, ".proto") + ".pb.sql"
	return filename
}

// preserveFileNameFunc converts "path/to/file.proto" to "path/to/file_methods.pb.sql"
func preserveFileNameFunc(protoPath string) string {
	filename := strings.TrimSuffix(protoPath, ".proto") + ".pb.sql"
	return filename
}

// singleFileNameFunc returns "protobuf_methods.sql" for all files (single file mode)
func singleFileNameFunc(protoPath string) string {
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
		return flattenFileNameFunc
	default:
		panic("unknown file naming strategy: " + strategy)
	}
}
