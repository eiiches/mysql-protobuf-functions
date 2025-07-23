package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlinstrument"
)

func main() {
	cmd := &cli.Command{
		Name:  "mysql-coverage",
		Usage: "Instrument MySQL stored procedures and functions for coverage analysis",
		Commands: []*cli.Command{
			{
				Name:      "instrument",
				Usage:     "Instrument SQL files with coverage tracking calls",
				ArgsUsage: "[file1.sql file2.sql ...]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "output",
						Usage: "Output directory (only used with multiple files)",
					},
				},
				Action: instrumentAction,
			},
			{
				Name:  "init",
				Usage: "Initialize database with coverage tracking schema",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "database",
						Usage:    "Database connection string. Example: user:password@tcp(127.0.0.1:3306)/dbname",
						Required: true,
					},
				},
				Action: initAction,
			},
			{
				Name:  "lcov",
				Usage: "Generate LCOV format coverage report from MysqlCoverageEvents table",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "database",
						Usage:    "Database connection string. Example: user:password@tcp(127.0.0.1:3306)/dbname",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "output",
						Usage: "Output file (default: stdout)",
					},
					&cli.StringSliceFlag{
						Name:  "instrumented-file",
						Usage: "Path(s) to instrumented SQL file(s) (auto-detected if not specified)",
					},
				},
				Action: lcovAction,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func instrumentAction(ctx context.Context, command *cli.Command) error {
	args := command.Args().Slice()
	outputDir := command.String("output")

	// If no arguments, read from stdin
	if len(args) == 0 {
		return instrumentSQL(os.Stdin, os.Stdout, "stdin")
	}

	// Process multiple files
	for _, inputFilename := range args {
		var input io.Reader
		var output io.Writer

		// Open input file
		file, err := os.Open(inputFilename)
		if err != nil {
			return fmt.Errorf("failed to open %s: %v", inputFilename, err)
		}
		defer file.Close()
		input = file

		// Determine output file
		var outputFile string
		if outputDir != "" {
			// Output to specified directory
			baseFilename := filepath.Base(inputFilename)
			outputFile = filepath.Join(outputDir, baseFilename+".instrumented")
		} else {
			// Default to {input}.instrumented naming convention
			outputFile = inputFilename + ".instrumented"
		}

		// Open output file
		outFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create %s: %v", outputFile, err)
		}
		defer outFile.Close()
		output = outFile

		// Process the file
		baseFilename := filepath.Base(inputFilename)
		if err := instrumentSQL(input, output, baseFilename); err != nil {
			return fmt.Errorf("failed to instrument %s: %v", inputFilename, err)
		}

		fmt.Printf("Instrumented %s -> %s\n", inputFilename, outputFile)
	}

	return nil
}

func instrumentSQL(input io.Reader, output io.Writer, filename string) error {
	// Read all input content
	content, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Create AST-based instrumenter
	instrumenter := sqlinstrument.NewInstrumenter(filename)

	// Instrument the SQL content
	instrumentedSQL, err := instrumenter.InstrumentSQL(content)
	if err != nil {
		return fmt.Errorf("failed to instrument SQL: %w", err)
	}

	// Write the instrumented SQL to output
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	if _, err := writer.WriteString(instrumentedSQL); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func lcovAction(ctx context.Context, command *cli.Command) error {
	var output io.Writer = os.Stdout

	if outputFile := command.String("output"); outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		output = file
	}

	db, err := sql.Open("mysql", command.String("database"))
	if err != nil {
		return err
	}
	defer db.Close()

	// We need the instrumented files to know which lines were instrumented
	instrumentedFiles := command.StringSlice("instrumented-file")
	if len(instrumentedFiles) == 0 {
		// Auto-detect instrumented files using the new naming convention
		pattern := "*.sql.instrumented"
		matches, err := filepath.Glob(pattern)
		if err == nil {
			instrumentedFiles = append(instrumentedFiles, matches...)
		}
		// Also check for old naming convention for backward compatibility
		oldFiles := []string{"instrumented-protobuf.sql", "instrumented-protobuf-accessors.sql", "instrumented-protobuf-descriptor.sql", "instrumented-protobuf-json.sql", "instrumented.sql"}
		for _, file := range oldFiles {
			if _, err := os.Stat(file); err == nil {
				instrumentedFiles = append(instrumentedFiles, file)
			}
		}
	}

	return generateLCOVReport(db, output, instrumentedFiles)
}

type CoverageData struct {
	Filename   string
	Function   string
	LineNumber int
	HitCount   int
}

func generateLCOVReport(db *sql.DB, output io.Writer, instrumentedFiles []string) error {
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// Parse all instrumented files to get all instrumented lines
	instrumentedLines := make(map[string]map[int]string)
	for _, instrumentedFile := range instrumentedFiles {
		fileLines, err := parseInstrumentedLines(instrumentedFile)
		if err != nil {
			return fmt.Errorf("failed to parse instrumented file %s: %v", instrumentedFile, err)
		}
		// Merge into the main map
		for filename, lines := range fileLines {
			if instrumentedLines[filename] == nil {
				instrumentedLines[filename] = make(map[int]string)
			}
			for lineNumber, function := range lines {
				instrumentedLines[filename][lineNumber] = function
			}
		}
	}

	// Query to get coverage data
	query := `
		SELECT
			filename,
			function_name,
			line_number,
			COUNT(*) as hit_count
		FROM __CoverageEvent
		GROUP BY filename, function_name, line_number
		ORDER BY filename, function_name, line_number
	`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create map of hit counts by line
	hitCounts := make(map[string]map[int]int) // filename -> line_number -> hit_count
	var allFilenames []string

	for rows.Next() {
		var data CoverageData
		if err := rows.Scan(&data.Filename, &data.Function, &data.LineNumber, &data.HitCount); err != nil {
			return err
		}

		if hitCounts[data.Filename] == nil {
			hitCounts[data.Filename] = make(map[int]int)
			allFilenames = append(allFilenames, data.Filename)
		}
		hitCounts[data.Filename][data.LineNumber] = data.HitCount
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Merge instrumented lines with hit counts
	var coverageData []CoverageData
	for filename, lines := range instrumentedLines {
		for lineNumber, function := range lines {
			hitCount := 0
			if hitCounts[filename] != nil {
				if count, exists := hitCounts[filename][lineNumber]; exists {
					hitCount = count
				}
			}
			coverageData = append(coverageData, CoverageData{
				Filename:   filename,
				Function:   function,
				LineNumber: lineNumber,
				HitCount:   hitCount,
			})
		}
	}

	// Sort coverage data by filename, function, line number
	sort.Slice(coverageData, func(i, j int) bool {
		if coverageData[i].Filename != coverageData[j].Filename {
			return coverageData[i].Filename < coverageData[j].Filename
		}
		if coverageData[i].Function != coverageData[j].Function {
			return coverageData[i].Function < coverageData[j].Function
		}
		return coverageData[i].LineNumber < coverageData[j].LineNumber
	})

	// Group by filename
	fileGroups := make(map[string][]CoverageData)
	for _, data := range coverageData {
		fileGroups[data.Filename] = append(fileGroups[data.Filename], data)
	}

	// Generate LCOV format
	for filename, fileData := range fileGroups {
		// Test name (optional)
		if _, err := writer.WriteString("TN:\n"); err != nil {
			return err
		}

		// Source file
		if _, err := writer.WriteString(fmt.Sprintf("SF:%s\n", filename)); err != nil {
			return err
		}

		// Function data
		functionGroups := make(map[string][]CoverageData)
		for _, data := range fileData {
			functionGroups[data.Function] = append(functionGroups[data.Function], data)
		}

		// Function names and line numbers
		for funcName, funcData := range functionGroups {
			if len(funcData) > 0 {
				if _, err := writer.WriteString(fmt.Sprintf("FN:%d,%s\n", funcData[0].LineNumber, funcName)); err != nil {
					return err
				}
			}
		}

		// Function hit counts
		for funcName, funcData := range functionGroups {
			hitCount := 0
			for _, data := range funcData {
				if data.HitCount > 0 {
					hitCount++
					break
				}
			}
			if _, err := writer.WriteString(fmt.Sprintf("FNDA:%d,%s\n", hitCount, funcName)); err != nil {
				return err
			}
		}

		// Functions found and hit
		functionsFound := len(functionGroups)
		functionsHit := 0
		for _, funcData := range functionGroups {
			for _, data := range funcData {
				if data.HitCount > 0 {
					functionsHit++
					break
				}
			}
		}
		if _, err := writer.WriteString(fmt.Sprintf("FNF:%d\n", functionsFound)); err != nil {
			return err
		}
		if _, err := writer.WriteString(fmt.Sprintf("FNH:%d\n", functionsHit)); err != nil {
			return err
		}

		// Line coverage data
		for _, data := range fileData {
			if _, err := writer.WriteString(fmt.Sprintf("DA:%d,%d\n", data.LineNumber, data.HitCount)); err != nil {
				return err
			}
		}

		// Lines found and hit
		linesFound := len(fileData)
		linesHit := 0
		for _, data := range fileData {
			if data.HitCount > 0 {
				linesHit++
			}
		}
		if _, err := writer.WriteString(fmt.Sprintf("LF:%d\n", linesFound)); err != nil {
			return err
		}
		if _, err := writer.WriteString(fmt.Sprintf("LH:%d\n", linesHit)); err != nil {
			return err
		}

		// End of record
		if _, err := writer.WriteString("end_of_record\n"); err != nil {
			return err
		}
	}

	return nil
}

func initAction(ctx context.Context, command *cli.Command) error {
	db, err := sql.Open("mysql", command.String("database"))
	if err != nil {
		return err
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Drop and recreate the __CoverageEvent table
	dropTableSQL := `DROP TABLE IF EXISTS __CoverageEvent`
	if _, err := db.Exec(dropTableSQL); err != nil {
		return fmt.Errorf("failed to drop existing __CoverageEvent table: %v", err)
	}

	createTableSQL := `
		CREATE TABLE __CoverageEvent (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			filename VARCHAR(255) NOT NULL,
			function_name VARCHAR(255) NOT NULL,
			line_number INT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_coverage (filename, function_name, line_number)
		)`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create __CoverageEvent table: %v", err)
	}

	// Drop existing __record_coverage procedure if it exists
	dropProcedureSQL := `DROP PROCEDURE IF EXISTS __record_coverage`
	if _, err := db.Exec(dropProcedureSQL); err != nil {
		return fmt.Errorf("failed to drop existing __record_coverage procedure: %v", err)
	}

	// Create the __record_coverage procedure
	createProcedureSQL := `
		CREATE PROCEDURE __record_coverage(IN filename VARCHAR(255), IN function_name VARCHAR(255), IN line_number INT)
		BEGIN
			INSERT INTO __CoverageEvent (filename, function_name, line_number)
			VALUES (filename, function_name, line_number);
		END`

	if _, err := db.Exec(createProcedureSQL); err != nil {
		return fmt.Errorf("failed to create __record_coverage procedure: %v", err)
	}

	fmt.Println("Successfully initialized database with coverage tracking schema")
	fmt.Println("- Recreated __CoverageEvent table")
	fmt.Println("- Recreated __record_coverage procedure")

	return nil
}

func parseInstrumentedLines(instrumentedFile string) (map[string]map[int]string, error) {
	if instrumentedFile == "" {
		return make(map[string]map[int]string), nil
	}

	file, err := os.Open(instrumentedFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	instrumentedLines := make(map[string]map[int]string) // filename -> line_number -> function_name
	scanner := bufio.NewScanner(file)

	// Regex to match __record_coverage calls
	instrumentCallPattern := regexp.MustCompile(`CALL __record_coverage\('([^']+)', '([^']+)', (\d+)\)`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := instrumentCallPattern.FindStringSubmatch(line); matches != nil {
			filename := matches[1]
			functionName := matches[2]
			lineNumber, err := strconv.Atoi(matches[3])
			if err != nil {
				continue
			}

			if instrumentedLines[filename] == nil {
				instrumentedLines[filename] = make(map[int]string)
			}
			instrumentedLines[filename][lineNumber] = functionName
		}
	}

	return instrumentedLines, scanner.Err()
}
