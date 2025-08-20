package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlftrace"
)

func main() {
	cmd := &cli.Command{
		Name:  "mysql-ftrace",
		Usage: "Instrument MySQL stored procedures and functions for call tracing",
		Commands: []*cli.Command{
			{
				Name:      "instrument",
				Usage:     "Instrument SQL files with function call tracing",
				ArgsUsage: "file1.sql [file2.sql ...]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "trace-statements",
						Usage: "Enable statement-level tracing within functions (adds significant overhead)",
					},
				},
				Action: instrumentAction,
			},
			{
				Name:  "init",
				Usage: "Initialize database with function tracing schema",
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
				Name:      "report",
				Usage:     "Generate function call trace report from FtraceEvent table",
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
					&cli.StringFlag{
						Name:  "format",
						Usage: "Output format: text, json, flamegraph",
						Value: "text",
					},
					&cli.IntFlag{
						Name:  "connection-id",
						Usage: "Filter reports by specific connection ID (default: show all connections)",
					},
				},
				Action: reportAction,
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

		// Default to {input}.ftraced naming convention
		outputFile := inputFilename + ".ftraced"

		// Open output file
		outFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", outputFile, err)
		}
		defer outFile.Close()
		output = outFile

		// Process the file
		if err := instrumentSQL(input, output, inputFilename, command.Bool("trace-statements")); err != nil {
			return fmt.Errorf("failed to instrument %s: %w", inputFilename, err)
		}

		fmt.Printf("Instrumented %s -> %s\n", inputFilename, outputFile)
	}

	return nil
}

func instrumentSQL(input io.Reader, output io.Writer, filename string, traceStatements bool) error {
	// Read all input content
	content, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Create AST-based instrumenter for function tracing
	instrumenter := sqlftrace.NewInstrumenter(filename)
	instrumenter.SetStatementTracing(traceStatements)

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

	// Drop and recreate the __FtraceEvent table
	dropTableSQL := `DROP TABLE IF EXISTS __FtraceEvent`
	if _, err := db.Exec(dropTableSQL); err != nil {
		return fmt.Errorf("failed to drop existing __FtraceEvent table: %w", err)
	}

	createTableSQL := `
		CREATE TABLE __FtraceEvent (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			connection_id INT NOT NULL,
			filename VARCHAR(255) NOT NULL,
			function_name VARCHAR(255) NOT NULL,
			object_type ENUM('function', 'procedure', 'statement', 'variable') NOT NULL,
			call_type ENUM('entry', 'exit', 'statement', 'set_variable') NOT NULL,
			arguments JSON,
			return_value JSON,
			call_depth INT NOT NULL DEFAULT 0,
			line_number INT,
			statement_type VARCHAR(50),
			variable_assignments JSON,
			timestamp TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6)
		) ENGINE = ARCHIVE`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create __FtraceEvent table: %w", err)
	}

	// Drop existing ftrace procedures if they exist
	dropProcedures := []string{
		`DROP PROCEDURE IF EXISTS __record_ftrace_entry`,
		`DROP PROCEDURE IF EXISTS __record_ftrace_exit`,
		`DROP PROCEDURE IF EXISTS __record_ftrace_statement`,
		`DROP PROCEDURE IF EXISTS __record_ftrace_set`,
		`DROP FUNCTION IF EXISTS __get_call_depth`,
		`DROP PROCEDURE IF EXISTS __increment_call_depth`,
		`DROP PROCEDURE IF EXISTS __decrement_call_depth`,
	}

	for _, dropSQL := range dropProcedures {
		if _, err := db.Exec(dropSQL); err != nil {
			return fmt.Errorf("failed to drop existing procedures: %w", err)
		}
	}

	// Create call depth management function
	createCallDepthSQL := `
		CREATE FUNCTION __get_call_depth() RETURNS INT READS SQL DATA DETERMINISTIC
		BEGIN
			DECLARE depth INT DEFAULT 0;
			SELECT COALESCE(@__ftrace_call_depth, 0) INTO depth;
			RETURN depth;
		END`

	if _, err := db.Exec(createCallDepthSQL); err != nil {
		return fmt.Errorf("failed to create __get_call_depth function: %w", err)
	}

	// Create increment call depth procedure
	createIncrementDepthSQL := `
		CREATE PROCEDURE __increment_call_depth()
		BEGIN
			SET @__ftrace_call_depth = COALESCE(@__ftrace_call_depth, 0) + 1;
		END`

	if _, err := db.Exec(createIncrementDepthSQL); err != nil {
		return fmt.Errorf("failed to create __increment_call_depth procedure: %w", err)
	}

	// Create decrement call depth procedure
	createDecrementDepthSQL := `
		CREATE PROCEDURE __decrement_call_depth()
		BEGIN
			SET @__ftrace_call_depth = GREATEST(COALESCE(@__ftrace_call_depth, 0) - 1, 0);
		END`

	if _, err := db.Exec(createDecrementDepthSQL); err != nil {
		return fmt.Errorf("failed to create __decrement_call_depth procedure: %w", err)
	}

	// Create the __record_ftrace_entry procedure
	createEntryProcedureSQL := `
		CREATE PROCEDURE __record_ftrace_entry(IN filename VARCHAR(255), IN function_name VARCHAR(255), IN object_type VARCHAR(10), IN arguments JSON)
		BEGIN
			CALL __increment_call_depth();
			INSERT INTO __FtraceEvent (connection_id, filename, function_name, object_type, call_type, arguments, call_depth)
			VALUES (CONNECTION_ID(), filename, function_name, object_type, 'entry', arguments, __get_call_depth());
		END`

	if _, err := db.Exec(createEntryProcedureSQL); err != nil {
		return fmt.Errorf("failed to create __record_ftrace_entry procedure: %w", err)
	}

	// Create the __record_ftrace_exit procedure
	createExitProcedureSQL := `
		CREATE PROCEDURE __record_ftrace_exit(IN filename VARCHAR(255), IN function_name VARCHAR(255), IN object_type VARCHAR(10), IN return_value JSON)
		BEGIN
			INSERT INTO __FtraceEvent (connection_id, filename, function_name, object_type, call_type, return_value, call_depth)
			VALUES (CONNECTION_ID(), filename, function_name, object_type, 'exit', return_value, __get_call_depth());
			CALL __decrement_call_depth();
		END`

	if _, err := db.Exec(createExitProcedureSQL); err != nil {
		return fmt.Errorf("failed to create __record_ftrace_exit procedure: %w", err)
	}

	// Create the __record_ftrace_statement procedure
	createStatementProcedureSQL := `
		CREATE PROCEDURE __record_ftrace_statement(IN filename VARCHAR(255), IN function_name VARCHAR(255), IN line_number INT, IN statement_type VARCHAR(50), IN statement_text TEXT)
		BEGIN
			INSERT INTO __FtraceEvent (connection_id, filename, function_name, object_type, call_type, call_depth, line_number, statement_type, return_value)
			VALUES (CONNECTION_ID(), filename, function_name, 'statement', 'statement', __get_call_depth(), line_number, statement_type, JSON_QUOTE(statement_text));
		END`

	if _, err := db.Exec(createStatementProcedureSQL); err != nil {
		return fmt.Errorf("failed to create __record_ftrace_statement procedure: %w", err)
	}

	// Create the __record_ftrace_set procedure
	createSetProcedureSQL := `
		CREATE PROCEDURE __record_ftrace_set(IN filename VARCHAR(255), IN function_name VARCHAR(255), IN line_number INT, IN variable_assignments JSON)
		BEGIN
			INSERT INTO __FtraceEvent (connection_id, filename, function_name, object_type, call_type, call_depth, line_number, statement_type, variable_assignments)
			VALUES (CONNECTION_ID(), filename, function_name, 'variable', 'set_variable', __get_call_depth(), line_number, 'SET', variable_assignments);
		END`

	if _, err := db.Exec(createSetProcedureSQL); err != nil {
		return fmt.Errorf("failed to create __record_ftrace_set procedure: %w", err)
	}

	fmt.Println("Successfully initialized database with function tracing schema")
	fmt.Println("- Created __FtraceEvent table")
	fmt.Println("- Created __record_ftrace_entry procedure")
	fmt.Println("- Created __record_ftrace_exit procedure")
	fmt.Println("- Created __record_ftrace_statement procedure")
	fmt.Println("- Created __record_ftrace_set procedure")
	fmt.Println("- Created call depth management functions")

	return nil
}

func reportAction(ctx context.Context, command *cli.Command) error {
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

	format := command.String("format")
	connectionID := command.Int("connection-id")

	switch format {
	case "text":
		return generateTextReport(db, output, connectionID)
	case "json":
		return generateJSONReport(db, output, connectionID)
	case "flamegraph":
		return generateFlamegraphReport(db, output, connectionID)
	default:
		return fmt.Errorf("unsupported format: %s (supported: text, json, flamegraph)", format)
	}
}

// FtraceEvent represents a function trace event
type FtraceEvent struct {
	ID                  int64
	ConnectionID        int
	Filename            string
	FunctionName        string
	ObjectType          string
	CallType            string
	Arguments           string
	ReturnValue         string
	CallDepth           int
	LineNumber          int
	StatementType       string
	VariableAssignments string
	Timestamp           string
}

func generateTextReport(db *sql.DB, output io.Writer, connectionID int) error {
	var query string
	var args []interface{}

	if connectionID > 0 {
		query = `
			SELECT id, connection_id, filename, function_name, object_type, call_type, arguments, return_value, call_depth, 
			       line_number, statement_type, variable_assignments, CAST(timestamp AS DATETIME(6))
			FROM __FtraceEvent
			WHERE connection_id = ?
			ORDER BY timestamp, id
		`
		args = []interface{}{connectionID}
	} else {
		query = `
			SELECT id, connection_id, filename, function_name, object_type, call_type, arguments, return_value, call_depth,
			       line_number, statement_type, variable_assignments, CAST(timestamp AS DATETIME(6))
			FROM __FtraceEvent
			ORDER BY connection_id, timestamp, id
		`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	writer.WriteString("MySQL Function Call Trace Report\n")
	writer.WriteString("================================\n\n")

	var currentConnectionID int = -1
	var pendingStatement *FtraceEvent

	// Helper function to format arguments as key=[JSON value] pairs
	formatArgs := func(jsonArgs string) string {
		if jsonArgs == "" {
			return ""
		}
		// Parse JSON and convert to key=[JSON value] format
		var argMap map[string]interface{}
		if err := json.Unmarshal([]byte(jsonArgs), &argMap); err == nil {
			var parts []string
			for key, value := range argMap {
				// Convert each value to proper JSON
				if jsonValue, err := json.Marshal(value); err == nil {
					parts = append(parts, fmt.Sprintf("%s=%s", key, string(jsonValue)))
				} else {
					parts = append(parts, fmt.Sprintf("%s=%v", key, value))
				}
			}
			// Sort for consistent output
			for i := 0; i < len(parts); i++ {
				for j := i + 1; j < len(parts); j++ {
					if parts[i] > parts[j] {
						parts[i], parts[j] = parts[j], parts[i]
					}
				}
			}
			return strings.Join(parts, ", ")
		}
		return jsonArgs
	}

	// Helper function to write a pending statement with its variable assignment
	flushPendingStatement := func() {
		if pendingStatement != nil {
			indentLevel := pendingStatement.CallDepth - 1
			if indentLevel < 0 {
				indentLevel = 0
			}
			indent := strings.Repeat("    ", indentLevel)
			timestampDisplay := pendingStatement.Timestamp
			if t, err := time.Parse("2006-01-02 15:04:05.000000", pendingStatement.Timestamp); err == nil {
				timestampDisplay = t.Format("15:04:05.000")
			}

			lineDisplay := ""
			if pendingStatement.LineNumber > 0 {
				lineDisplay = fmt.Sprintf("L%d: ", pendingStatement.LineNumber)
			}

			// Statements are indented more than ENTER/EXIT
			stmtIndent := indent + "    "

			writer.WriteString(fmt.Sprintf("[%s] %s%s%s\n",
				timestampDisplay, stmtIndent, lineDisplay, pendingStatement.ReturnValue))
			pendingStatement = nil
		}
	}

	for rows.Next() {
		var event FtraceEvent
		var args, retVal, stmtType, varName sql.NullString
		var lineNum sql.NullInt64

		if err := rows.Scan(&event.ID, &event.ConnectionID, &event.Filename, &event.FunctionName,
			&event.ObjectType, &event.CallType, &args, &retVal, &event.CallDepth,
			&lineNum, &stmtType, &varName, &event.Timestamp); err != nil {
			return err
		}

		event.Arguments = args.String
		event.ReturnValue = retVal.String
		event.LineNumber = int(lineNum.Int64)
		event.StatementType = stmtType.String
		event.VariableAssignments = varName.String

		// Show connection ID separator when it changes (only if showing multiple connections)
		if connectionID == 0 && event.ConnectionID != currentConnectionID {
			flushPendingStatement() // Flush any pending statement before connection change
			if currentConnectionID != -1 {
				writer.WriteString("\n")
			}
			writer.WriteString(fmt.Sprintf("=== Connection ID: %d ===\n", event.ConnectionID))
			currentConnectionID = event.ConnectionID
		}

		// Create indentation based on call depth (subtract 1 since call depths start at 1)
		indentLevel := event.CallDepth - 1
		if indentLevel < 0 {
			indentLevel = 0
		}
		indent := strings.Repeat("    ", indentLevel)

		// Parse timestamp for display
		timestampDisplay := event.Timestamp
		if t, err := time.Parse("2006-01-02 15:04:05.000000", event.Timestamp); err == nil {
			timestampDisplay = t.Format("15:04:05.000")
		}

		switch event.CallType {
		case "entry":
			flushPendingStatement() // Flush any pending statement before function entry
			writer.WriteString(fmt.Sprintf("[%s] %sENTER %s(%s)\n",
				timestampDisplay, indent, event.FunctionName, formatArgs(event.Arguments)))

		case "exit":
			flushPendingStatement() // Flush any pending statement before function exit
			writer.WriteString(fmt.Sprintf("[%s] %sRETURN %s\n",
				timestampDisplay, indent, event.ReturnValue))

		case "statement":
			// Check if this is a SET statement that will have a corresponding set_variable event
			if event.StatementType == "SET" {
				flushPendingStatement()   // Flush any previous pending statement
				pendingStatement = &event // Store this statement to be combined with variable assignment
			} else {
				flushPendingStatement() // Flush any previous pending statement

				lineDisplay := ""
				if event.LineNumber > 0 {
					lineDisplay = fmt.Sprintf("L%d: ", event.LineNumber)
				}

				// Statements are indented more than ENTER/EXIT
				stmtIndent := indent + "    "

				// For control flow statements, show the condition result
				if event.StatementType == "IF" || event.StatementType == "WHILE" || event.StatementType == "CASE" {
					writer.WriteString(fmt.Sprintf("[%s] %s%s%s\n",
						timestampDisplay, stmtIndent, lineDisplay, event.ReturnValue))
				} else {
					writer.WriteString(fmt.Sprintf("[%s] %s%s%s\n",
						timestampDisplay, stmtIndent, lineDisplay, event.ReturnValue))
				}
			}

		case "set_variable":
			// Combine with pending SET statement
			if pendingStatement != nil && pendingStatement.LineNumber == event.LineNumber {
				lineDisplay := ""
				if event.LineNumber > 0 {
					lineDisplay = fmt.Sprintf("L%d: ", event.LineNumber)
				}

				// Statements are indented more than ENTER/EXIT
				stmtIndent := indent + "    "

				writer.WriteString(fmt.Sprintf("[%s] %s%s%-30s → %s=%s\n",
					timestampDisplay, stmtIndent, lineDisplay, pendingStatement.ReturnValue,
					event.VariableAssignments, event.ReturnValue))
				pendingStatement = nil
			} else {
				// Standalone set_variable event (shouldn't normally happen)
				flushPendingStatement()
				lineDisplay := ""
				if event.LineNumber > 0 {
					lineDisplay = fmt.Sprintf("L%d: ", event.LineNumber)
				}
				stmtIndent := indent + "    "
				writer.WriteString(fmt.Sprintf("[%s] %s%sSET %s = %s\n",
					timestampDisplay, stmtIndent, lineDisplay, event.VariableAssignments, event.ReturnValue))
			}
		}
	}

	flushPendingStatement() // Flush any final pending statement

	return rows.Err()
}

func generateJSONReport(db *sql.DB, output io.Writer, connectionID int) error {
	var query string
	var args []interface{}

	if connectionID > 0 {
		query = `
			SELECT id, connection_id, filename, function_name, object_type, call_type, arguments, return_value, call_depth,
			       line_number, statement_type, variable_assignments, CAST(timestamp AS DATETIME(6))
			FROM __FtraceEvent
			WHERE connection_id = ?
			ORDER BY timestamp, id
		`
		args = []interface{}{connectionID}
	} else {
		query = `
			SELECT id, connection_id, filename, function_name, object_type, call_type, arguments, return_value, call_depth,
			       line_number, statement_type, variable_assignments, CAST(timestamp AS DATETIME(6))
			FROM __FtraceEvent
			ORDER BY connection_id, timestamp, id
		`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	writer.WriteString("[\n")
	first := true

	for rows.Next() {
		var event FtraceEvent
		var args, retVal, stmtType, varName sql.NullString
		var lineNum sql.NullInt64

		if err := rows.Scan(&event.ID, &event.ConnectionID, &event.Filename, &event.FunctionName,
			&event.ObjectType, &event.CallType, &args, &retVal, &event.CallDepth,
			&lineNum, &stmtType, &varName, &event.Timestamp); err != nil {
			return err
		}

		event.Arguments = args.String
		event.ReturnValue = retVal.String
		event.LineNumber = int(lineNum.Int64)
		event.StatementType = stmtType.String
		event.VariableAssignments = varName.String

		if !first {
			writer.WriteString(",\n")
		}
		first = false

		writer.WriteString(fmt.Sprintf(`  {
    "id": %d,
    "connection_id": %d,
    "filename": "%s",
    "function_name": "%s",
    "object_type": "%s",
    "call_type": "%s",
    "arguments": "%s",
    "return_value": "%s",
    "call_depth": %d,
    "line_number": %d,
    "statement_type": "%s",
    "variable_assignments": "%s",
    "timestamp": "%s"
  }`, event.ID, event.ConnectionID, event.Filename, event.FunctionName, event.ObjectType, event.CallType,
			strings.ReplaceAll(event.Arguments, `"`, `\"`),
			strings.ReplaceAll(event.ReturnValue, `"`, `\"`),
			event.CallDepth, event.LineNumber, event.StatementType, event.VariableAssignments, event.Timestamp))
	}

	writer.WriteString("\n]\n")
	return rows.Err()
}

func generateFlamegraphReport(db *sql.DB, output io.Writer, connectionID int) error {
	// TODO: Implement flamegraph generation
	// This would generate data suitable for brendangregg/FlameGraph tools
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	writer.WriteString("# Flamegraph generation not yet implemented\n")
	writer.WriteString("# This would generate stack traces suitable for flamegraph tools\n")
	if connectionID > 0 {
		writer.WriteString(fmt.Sprintf("# Would filter by connection_id = %d\n", connectionID))
	}

	return nil
}
