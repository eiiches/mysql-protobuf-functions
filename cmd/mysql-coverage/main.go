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
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-github/v66/github"
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
				ArgsUsage: "file1.sql [file2.sql ...]",
				Action:    instrumentAction,
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
				Name:      "lcov",
				Usage:     "Generate LCOV format coverage report from MysqlCoverageEvents table",
				ArgsUsage: "[file1.sql file2.sql ...]",
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
				},
				Action: lcovAction,
			},
			{
				Name:  "github-comment",
				Usage: "Post coverage summary as GitHub pull request comment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "lcov-file",
						Usage:    "Path to LCOV coverage file",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "github-token",
						Usage:   "GitHub personal access token (can also use GITHUB_TOKEN env var)",
						Sources: cli.EnvVars("GITHUB_TOKEN"),
					},
					&cli.StringFlag{
						Name:    "github-repository",
						Usage:   "GitHub repository in format owner/repo (can also use GITHUB_REPOSITORY env var)",
						Sources: cli.EnvVars("GITHUB_REPOSITORY"),
					},
					&cli.IntFlag{
						Name:  "github-pr-number",
						Usage: "GitHub pull request number",
					},
					&cli.StringFlag{
						Name:  "title",
						Usage: "Title for the comment",
						Value: "📊 MySQL Coverage Report",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Output the comment to stdout instead of posting to GitHub",
					},
					&cli.StringSliceFlag{
						Name:  "exclude-file",
						Usage: "Exclude files from coverage rate calculation and output (can be repeated)",
					},
				},
				Action: githubCommentAction,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func instrumentAction(ctx context.Context, command *cli.Command) error {
	args := command.Args().Slice()

	// Require at least one file argument
	if len(args) == 0 {
		return fmt.Errorf("at least one SQL file must be specified")
	}

	// Process multiple files
	for _, inputFilename := range args {
		var input io.Reader
		var output io.Writer

		// Open input file
		file, err := os.Open(inputFilename)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", inputFilename, err)
		}
		defer file.Close()
		input = file

		// Default to {input}.instrumented naming convention
		outputFile := inputFilename + ".instrumented"

		// Open output file
		outFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", outputFile, err)
		}
		defer outFile.Close()
		output = outFile

		// Process the file
		if err := instrumentSQL(input, output, inputFilename); err != nil {
			return fmt.Errorf("failed to instrument %s: %w", inputFilename, err)
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

	// Get source files from positional arguments or auto-detect from database
	sourceFiles := command.Args().Slice()
	if len(sourceFiles) == 0 {
		// Auto-detect source files from database coverage events
		var err error
		sourceFiles, err = getSourceFilesFromDatabase(db)
		if err != nil {
			return fmt.Errorf("failed to auto-detect source files from database: %w", err)
		}
	}

	return generateLCOVReport(db, output, sourceFiles)
}

func getSourceFilesFromDatabase(db *sql.DB) ([]string, error) {
	query := "SELECT DISTINCT filename FROM __CoverageEvent ORDER BY filename"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sourceFiles []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		sourceFiles = append(sourceFiles, filename)
	}

	return sourceFiles, rows.Err()
}

type CoverageData struct {
	Filename   string
	Function   string
	LineNumber int
	HitCount   int
}

func generateLCOVReport(db *sql.DB, output io.Writer, sourceFiles []string) error {
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// Parse all instrumented files to get all instrumented lines
	instrumentedLines := make(map[string]map[int]string)
	for _, sourceFile := range sourceFiles {
		// Try to find corresponding instrumented file
		instrumentedFile := sourceFile + ".instrumented"
		fileLines, err := parseInstrumentedLines(instrumentedFile)
		if err != nil {
			// If instrumented file not found, skip this source file
			// This allows partial coverage reports when some files haven't been instrumented
			continue
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

	for rows.Next() {
		var data CoverageData
		if err := rows.Scan(&data.Filename, &data.Function, &data.LineNumber, &data.HitCount); err != nil {
			return err
		}

		if hitCounts[data.Filename] == nil {
			hitCounts[data.Filename] = make(map[int]int)
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
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Drop and recreate the __CoverageEvent table
	dropTableSQL := `DROP TABLE IF EXISTS __CoverageEvent`
	if _, err := db.Exec(dropTableSQL); err != nil {
		return fmt.Errorf("failed to drop existing __CoverageEvent table: %w", err)
	}

	createTableSQL := `
		CREATE TABLE __CoverageEvent (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			filename VARCHAR(255) NOT NULL,
			function_name VARCHAR(255) NOT NULL,
			line_number INT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE = ARCHIVE`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create __CoverageEvent table: %w", err)
	}

	// Drop existing __record_coverage procedure if it exists
	dropProcedureSQL := `DROP PROCEDURE IF EXISTS __record_coverage`
	if _, err := db.Exec(dropProcedureSQL); err != nil {
		return fmt.Errorf("failed to drop existing __record_coverage procedure: %w", err)
	}

	// Create the __record_coverage procedure
	createProcedureSQL := `
		CREATE PROCEDURE __record_coverage(IN filename VARCHAR(255), IN function_name VARCHAR(255), IN line_number INT)
		BEGIN
			INSERT INTO __CoverageEvent (filename, function_name, line_number)
			VALUES (filename, function_name, line_number);
		END`

	if _, err := db.Exec(createProcedureSQL); err != nil {
		return fmt.Errorf("failed to create __record_coverage procedure: %w", err)
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

// LCOV parsing structures
type LCOVFile struct {
	SourceFile     string
	Functions      []LCOVFunction
	Lines          []LCOVLine
	FunctionsFound int
	FunctionsHit   int
	LinesFound     int
	LinesHit       int
}

type LCOVFunction struct {
	Name       string
	LineNumber int
	HitCount   int
}

type LCOVLine struct {
	LineNumber int
	HitCount   int
}

type CoverageSummary struct {
	Files             []LCOVFile
	TotalFunctions    int
	TotalFunctionsHit int
	TotalLines        int
	TotalLinesHit     int
	FunctionCoverage  float64
	LineCoverage      float64
}

func githubCommentAction(ctx context.Context, command *cli.Command) error {
	lcovFile := command.String("lcov-file")
	githubToken := command.String("github-token")
	repo := command.String("github-repository")
	prNumber := command.Int("github-pr-number")
	title := command.String("title")
	dryRun := command.Bool("dry-run")

	// Validate required flags for GitHub posting
	if !dryRun {
		if githubToken == "" {
			return fmt.Errorf("--github-token is required when not using --dry-run")
		}
		if repo == "" {
			return fmt.Errorf("--github-repository is required when not using --dry-run")
		}
		if prNumber == 0 {
			return fmt.Errorf("--github-pr-number is required when not using --dry-run")
		}
	}

	// Parse LCOV file
	excludeFiles := command.StringSlice("exclude-file")
	summary, err := parseLCOVFile(lcovFile, excludeFiles)
	if err != nil {
		return fmt.Errorf("failed to parse LCOV file: %w", err)
	}

	// Generate coverage report
	report := generateCoverageReport(summary, title)

	// Output to stdout or post to GitHub
	if dryRun {
		fmt.Print(report)
		return nil
	}

	// Post to GitHub
	return postGitHubComment(githubToken, repo, prNumber, report)
}

func parseLCOVFile(filePath string, excludeFiles []string) (*CoverageSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Helper function to check if a file should be excluded
	shouldExclude := func(filename string) bool {
		for _, excludePattern := range excludeFiles {
			// Support both exact matches and simple wildcard patterns
			if matched, _ := filepath.Match(excludePattern, filename); matched {
				return true
			}
			// Also check basename matching
			if matched, _ := filepath.Match(excludePattern, filepath.Base(filename)); matched {
				return true
			}
		}
		return false
	}

	scanner := bufio.NewScanner(file)
	summary := &CoverageSummary{}
	var currentFile *LCOVFile

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Extract prefix (everything before first colon)
		prefix := line
		if colonIndex := strings.Index(line, ":"); colonIndex >= 0 {
			prefix = line[:colonIndex]
		}

		switch prefix {
		case "TN":
			// Test name - start of new record
			continue
		case "SF":
			// Source file
			if currentFile != nil {
				summary.Files = append(summary.Files, *currentFile)
			}
			sourceFile := strings.TrimPrefix(line, "SF:")
			// Skip excluded files
			if shouldExclude(sourceFile) {
				currentFile = nil
				continue
			}
			currentFile = &LCOVFile{
				SourceFile: sourceFile,
			}
		case "FN":
			// Function definition: FN:line_number,function_name
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			parts := strings.SplitN(strings.TrimPrefix(line, "FN:"), ",", 2)
			if len(parts) == 2 {
				lineNum, _ := strconv.Atoi(parts[0])
				currentFile.Functions = append(currentFile.Functions, LCOVFunction{
					Name:       parts[1],
					LineNumber: lineNum,
				})
			}
		case "FNDA":
			// Function hit data: FNDA:hit_count,function_name
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			parts := strings.SplitN(strings.TrimPrefix(line, "FNDA:"), ",", 2)
			if len(parts) == 2 {
				hitCount, _ := strconv.Atoi(parts[0])
				functionName := parts[1]
				// Find the function and update hit count
				for i := range currentFile.Functions {
					if currentFile.Functions[i].Name == functionName {
						currentFile.Functions[i].HitCount = hitCount
						break
					}
				}
			}
		case "FNF":
			// Functions found
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			currentFile.FunctionsFound, _ = strconv.Atoi(strings.TrimPrefix(line, "FNF:"))
		case "FNH":
			// Functions hit
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			currentFile.FunctionsHit, _ = strconv.Atoi(strings.TrimPrefix(line, "FNH:"))
		case "DA":
			// Line hit data: DA:line_number,hit_count
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			parts := strings.SplitN(strings.TrimPrefix(line, "DA:"), ",", 2)
			if len(parts) == 2 {
				lineNum, _ := strconv.Atoi(parts[0])
				hitCount, _ := strconv.Atoi(parts[1])
				currentFile.Lines = append(currentFile.Lines, LCOVLine{
					LineNumber: lineNum,
					HitCount:   hitCount,
				})
			}
		case "LF":
			// Lines found
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			currentFile.LinesFound, _ = strconv.Atoi(strings.TrimPrefix(line, "LF:"))
		case "LH":
			// Lines hit
			if currentFile == nil {
				continue // Skip if file is excluded
			}
			currentFile.LinesHit, _ = strconv.Atoi(strings.TrimPrefix(line, "LH:"))
		default:
			if line == "end_of_record" {
				// End of current file record
				if currentFile != nil {
					summary.Files = append(summary.Files, *currentFile)
					currentFile = nil
				}
			}
		}
	}

	// Add final file if exists
	if currentFile != nil {
		summary.Files = append(summary.Files, *currentFile)
	}

	// Calculate totals
	for _, file := range summary.Files {
		summary.TotalFunctions += file.FunctionsFound
		summary.TotalFunctionsHit += file.FunctionsHit
		summary.TotalLines += file.LinesFound
		summary.TotalLinesHit += file.LinesHit
	}

	// Calculate percentages
	if summary.TotalFunctions > 0 {
		summary.FunctionCoverage = float64(summary.TotalFunctionsHit) / float64(summary.TotalFunctions) * 100
	}
	if summary.TotalLines > 0 {
		summary.LineCoverage = float64(summary.TotalLinesHit) / float64(summary.TotalLines) * 100
	}

	return summary, scanner.Err()
}

func generateCoverageReport(summary *CoverageSummary, title string) string {
	var report strings.Builder

	report.WriteString(fmt.Sprintf("## %s\n\n", title))

	// Overall summary
	report.WriteString("### 📈 Overall Coverage\n\n")
	report.WriteString("| Metric | Coverage | Hit | Total |\n")
	report.WriteString("|--------|----------|-----|-------|\n")
	report.WriteString(fmt.Sprintf("| **Functions** | **%.1f%%** | %d | %d |\n",
		summary.FunctionCoverage, summary.TotalFunctionsHit, summary.TotalFunctions))
	report.WriteString(fmt.Sprintf("| **Lines** | **%.1f%%** | %d | %d |\n",
		summary.LineCoverage, summary.TotalLinesHit, summary.TotalLines))
	report.WriteString("\n")

	// Per-file coverage
	if len(summary.Files) > 1 {
		report.WriteString("### 📁 Coverage by File\n\n")
		report.WriteString("| File | Function Coverage | Line Coverage | Functions | Lines |\n")
		report.WriteString("|------|-------------------|---------------|-----------|-------|\n")

		for _, file := range summary.Files {
			functionPct := 0.0
			if file.FunctionsFound > 0 {
				functionPct = float64(file.FunctionsHit) / float64(file.FunctionsFound) * 100
			}
			linePct := 0.0
			if file.LinesFound > 0 {
				linePct = float64(file.LinesHit) / float64(file.LinesFound) * 100
			}

			// Get coverage emoji
			funcEmoji := getCoverageEmoji(functionPct)
			lineEmoji := getCoverageEmoji(linePct)

			report.WriteString(fmt.Sprintf("| `%s` | %s %.1f%% | %s %.1f%% | %d/%d | %d/%d |\n",
				filepath.Base(file.SourceFile),
				funcEmoji, functionPct,
				lineEmoji, linePct,
				file.FunctionsHit, file.FunctionsFound,
				file.LinesHit, file.LinesFound))
		}
		report.WriteString("\n")
	}

	// Coverage quality assessment
	report.WriteString("### 🎯 Coverage Quality\n\n")
	overallEmoji := getCoverageEmoji(summary.LineCoverage)
	qualityMessage := getCoverageQualityMessage(summary.LineCoverage)
	report.WriteString(fmt.Sprintf("%s **%s** (%.1f%% line coverage)\n\n", overallEmoji, qualityMessage, summary.LineCoverage))

	// Footer
	report.WriteString("---\n")
	report.WriteString("*Generated by mysql-coverage tool*")

	return report.String()
}

func getCoverageEmoji(percentage float64) string {
	switch {
	case percentage >= 90:
		return "🟢"
	case percentage >= 70:
		return "🟡"
	case percentage >= 50:
		return "🟠"
	default:
		return "🔴"
	}
}

func getCoverageQualityMessage(percentage float64) string {
	switch {
	case percentage >= 90:
		return "Excellent Coverage"
	case percentage >= 70:
		return "Good Coverage"
	case percentage >= 50:
		return "Moderate Coverage"
	default:
		return "Low Coverage"
	}
}

func postGitHubComment(token, repo string, prNumber int, comment string) error {
	// Parse owner/repo
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format, expected 'owner/repo', got '%s'", repo)
	}
	owner, repoName := parts[0], parts[1]

	// Create GitHub client
	client := github.NewClient(nil).WithAuthToken(token)

	// Create the comment
	issueComment := &github.IssueComment{
		Body: &comment,
	}

	// Post the comment
	ctx := context.Background()
	_, _, err := client.Issues.CreateComment(ctx, owner, repoName, prNumber, issueComment)
	if err != nil {
		return fmt.Errorf("failed to post GitHub comment: %w", err)
	}

	fmt.Printf("Successfully posted coverage comment to PR #%d in %s\n", prNumber, repo)
	return nil
}
