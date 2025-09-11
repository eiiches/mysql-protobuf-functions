package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
	"github.com/urfave/cli/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Check if running in standalone mode (any command line arguments)
	if len(os.Args) > 1 {
		runStandalone()
		return
	}

	// Original protoc plugin mode
	runAsProtocPlugin()
}

func runStandalone() {
	app := &cli.Command{
		Name:  "protoc-gen-mysql",
		Usage: "Generate MySQL stored functions and method helpers for protobuf schemas",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "descriptor_set_in",
				Usage:    "Path to binary FileDescriptorSet file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the generated SQL function",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "include_source_info",
				Usage: "Include source code info in output (increases output size significantly)",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "mysql_out",
				Usage: "Output directory for generated SQL files",
				Value: ".",
			},
			&cli.StringFlag{
				Name:  "file_naming_strategy",
				Usage: "Naming strategy for generated method files: flatten, preserve, or single",
				Value: "single",
			},
			&cli.StringFlag{
				Name:  "prefix_map",
				Usage: "Map proto packages or types to function prefixes (e.g., 'google.protobuf=pb_,com.example.MyMessage=msg_')",
				Value: "",
			},
			&cli.BoolFlag{
				Name:  "generate_methods",
				Usage: "Generate type method functions (constructors, getters, setters)",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "include_wkt",
				Usage: "Include well-known types (google.protobuf.*) in descriptor set and method generation",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "validate",
				Usage: "Validate the file descriptor set",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "skip_functions",
				Usage: "Comma-separated list of function names to skip (generate as commented out)",
				Value: "",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			descriptorSetIn := cmd.String("descriptor_set_in")
			name := cmd.String("name")
			includeSourceInfo := cmd.Bool("include_source_info")
			mysqlOut := cmd.String("mysql_out")
			namingStrategy := cmd.String("file_naming_strategy")
			prefixMapStr := cmd.String("prefix_map")
			generateMethods := cmd.Bool("generate_methods")
			includeWkt := cmd.Bool("include_wkt")
			validate := cmd.Bool("validate")
			skipFunctionsStr := cmd.String("skip_functions")

			// Parse skip functions list
			var skippedFunctions []string
			if skipFunctionsStr != "" {
				skippedFunctions = strings.Split(skipFunctionsStr, ",")
				for i := range skippedFunctions {
					skippedFunctions[i] = strings.TrimSpace(skippedFunctions[i])
				}
			}

			// Read binary FileDescriptorSet from file
			data, err := os.ReadFile(descriptorSetIn)
			if err != nil {
				return fmt.Errorf("failed to read descriptor set file: %w", err)
			}

			var fileDescriptorSet descriptorpb.FileDescriptorSet
			if unmarshalErr := proto.Unmarshal(data, &fileDescriptorSet); unmarshalErr != nil {
				return fmt.Errorf("failed to unmarshal FileDescriptorSet: %w", unmarshalErr)
			}

			// Create generation config
			config := protocgenmysql.GenerateConfig{
				DescriptorSetName: name,
				IncludeSourceInfo: includeSourceInfo,
				GenerateMethods:   generateMethods,
				IncludeWkt:        includeWkt,
				FileNameFunc:      getFileNameFunc(namingStrategy),
				TypePrefixFunc:    createTypePrefixFunc(parsePrefixMap(prefixMapStr)),
				FieldFilterFunc:   createFieldFilterFunc(skippedFunctions),
			}

			if validate {
				_, err := protodesc.NewFiles(&fileDescriptorSet)
				if err != nil {
					return fmt.Errorf("validation failed: %w", err)
				}
			}

			// Process and generate files
			response, err := protocgenmysql.Generate(&fileDescriptorSet, config)
			if err != nil {
				return err
			}

			// Write all files
			for _, file := range response.File {
				if file.Name == nil || file.Content == nil {
					continue
				}
				outputFile := filepath.Join(mysqlOut, *file.Name)
				//nolint:gosec // 0o644 permissions are intentional for generated SQL files
				if err := os.WriteFile(outputFile, []byte(*file.Content), 0o644); err != nil {
					return fmt.Errorf("failed to write file: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Generated %s\n", outputFile)
			}

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runAsProtocPlugin() {
	// Read CodeGeneratorRequest from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	var req pluginpb.CodeGeneratorRequest
	if unmarshalErr := proto.Unmarshal(input, &req); unmarshalErr != nil {
		log.Fatalf("Failed to unmarshal CodeGeneratorRequest: %v", unmarshalErr)
	}

	// Parse plugin options
	var descriptorSetName string
	includeSourceInfo := false
	namingStrategy := "flatten"
	generateMethods := true
	includeWkt := false
	prefixMap := make(map[protoreflect.FullName]string)
	var skippedFunctions []string
	if req.Parameter != nil && *req.Parameter != "" {
		params := parseParameters(*req.Parameter)
		if name, ok := params["name"]; ok {
			descriptorSetName = name
		}
		if include, ok := params["include_source_info"]; ok {
			includeSourceInfo = include == "true"
		}
		if strategy, ok := params["file_naming_strategy"]; ok {
			namingStrategy = strategy
		}
		if methods, ok := params["generate_methods"]; ok {
			generateMethods = methods == "true"
		}
		if wkt, ok := params["include_wkt"]; ok {
			includeWkt = wkt == "true"
		}
		if prefix, ok := params["prefix_map"]; ok {
			prefixMap = parsePrefixMap(prefix)
		}
		if skipFuncs, ok := params["skip_functions"]; ok {
			skippedFunctions = strings.Split(skipFuncs, ",")
			for i := range skippedFunctions {
				skippedFunctions[i] = strings.TrimSpace(skippedFunctions[i])
			}
		}
	}

	if descriptorSetName == "" {
		sendError("name parameter is required. Use --mysql_opt=name=your_function_name")
		return
	}

	// Build FileDescriptorSet from the request
	fileDescriptorSet := &descriptorpb.FileDescriptorSet{
		File: req.ProtoFile,
	}

	// Create generation config
	config := protocgenmysql.GenerateConfig{
		DescriptorSetName: descriptorSetName,
		IncludeSourceInfo: includeSourceInfo,
		GenerateMethods:   generateMethods,
		IncludeWkt:        includeWkt,
		FileNameFunc:      getFileNameFunc(namingStrategy),
		TypePrefixFunc:    createTypePrefixFunc(prefixMap),
		FieldFilterFunc:   createFieldFilterFunc(skippedFunctions),
	}

	// Process and generate files
	response, err := protocgenmysql.Generate(fileDescriptorSet, config)
	if err != nil {
		sendError(err.Error())
		return
	}

	// Marshal and write response
	output, err := proto.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}
}

func sendError(message string) {
	response := &pluginpb.CodeGeneratorResponse{
		Error: proto.String(message),
	}

	output, err := proto.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal error response: %v", err)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		log.Fatalf("Failed to write error output: %v", err)
	}
}

func parseParameters(paramStr string) map[string]string {
	params := make(map[string]string)
	pairs := strings.Split(paramStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			params[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return params
}

func parsePrefixMap(mapStr string) map[protoreflect.FullName]string {
	result := make(map[protoreflect.FullName]string)
	if mapStr == "" {
		return result
	}

	pairs := strings.Split(mapStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			pkg := strings.TrimSpace(kv[0])
			prefix := strings.TrimSpace(kv[1])
			result[protoreflect.FullName(pkg)] = prefix
		}
	}
	return result
}

// createFieldFilterFunc creates a FieldFilterFunc from a list of skipped function names
func createFieldFilterFunc(skippedFunctions []string) protocgenmysql.FieldFilterFunc {
	if len(skippedFunctions) == 0 {
		return nil
	}

	skippedSet := make(map[string]bool)
	for _, funcName := range skippedFunctions {
		skippedSet[funcName] = true
	}

	return func(field protoreflect.FieldDescriptor, functionName string) protocgenmysql.FilterDecision {
		if skippedSet[functionName] {
			return protocgenmysql.DecisionCommentOut
		}
		return protocgenmysql.DecisionInclude
	}
}
