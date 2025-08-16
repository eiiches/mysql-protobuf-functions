package sqlftrace

import (
	"fmt"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlsplitter"
)

// isBinaryType checks if a MySQL data type is a binary type that needs base64 encoding
func isBinaryType(dataType string) bool {
	if dataType == "" {
		return false
	}

	// Convert to uppercase for case-insensitive comparison
	upperType := strings.ToUpper(strings.TrimSpace(dataType))

	// Check for binary types
	return strings.Contains(upperType, "BLOB") ||
		strings.Contains(upperType, "BINARY") ||
		strings.HasPrefix(upperType, "VARBINARY")
}

// ASTInstrumenter adds function tracing instrumentation using AST reconstruction
type ASTInstrumenter struct {
	filename string
	codegen  *CodeGenerator
}

// NewASTInstrumenter creates a new AST-based instrumenter for function tracing
func NewASTInstrumenter(filename string) *ASTInstrumenter {
	return &ASTInstrumenter{
		filename: filename,
		codegen:  NewCodeGenerator(),
	}
}

// InstrumentSQL adds function tracing instrumentation to SQL content using AST reconstruction
func (i *ASTInstrumenter) InstrumentSQL(content []byte) (string, error) {
	// First, split the SQL into statements
	splitter := sqlsplitter.NewParser(content)
	statements, err := splitter.Parse()
	if err != nil {
		return "", fmt.Errorf("failed to split SQL: %w", err)
	}

	var result strings.Builder
	currentDelimiter := ";"

	for _, stmt := range statements {
		switch stmt.Type {
		case "DELIMITER":
			// Update current delimiter and pass through unchanged
			parts := strings.Fields(stmt.Text)
			if len(parts) >= 2 {
				currentDelimiter = parts[1]
			}
			result.WriteString(stmt.Text)
			result.WriteString("\n")
		case "COMMENT":
			// Pass through comment statements unchanged
			result.WriteString(stmt.Text)
			result.WriteString("\n")
		case "SQL":
			// Instrument SQL statements using AST reconstruction
			instrumented, err := i.instrumentStatementWithAST(stmt)
			if err != nil {
				return "", fmt.Errorf("failed to instrument statement at line %d: %w", stmt.LineNo, err)
			}
			result.WriteString(instrumented)
			result.WriteString(" ")
			result.WriteString(currentDelimiter)
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// instrumentStatementWithAST instruments a statement using AST reconstruction for function tracing
func (i *ASTInstrumenter) instrumentStatementWithAST(stmt sqlsplitter.Statement) (string, error) {
	// Parse the statement with sqlflowparser to get the AST
	ast, err := sqlflowparser.Parse("", []byte(stmt.Text))
	if err != nil {
		// Extract statement type from the beginning of the text for better error context
		stmtType := "unknown"
		textLines := strings.Split(stmt.Text, "\n")
		if len(textLines) > 0 {
			firstLine := strings.TrimSpace(textLines[0])
			if len(firstLine) > 50 {
				firstLine = firstLine[:50] + "..."
			}
			stmtType = firstLine
		}

		// If parsing fails, report the error with file context
		return "", fmt.Errorf("failed to parse statement starting at file line %d (%s): %w", stmt.LineNo, stmtType, err)
	}

	switch node := ast.(type) {
	case *sqlflowparser.CreateFunctionStmt:
		instrumentedAST := i.instrumentCreateFunctionAST(stmt, node)
		return i.codegen.GenerateSQL(instrumentedAST), nil
	case *sqlflowparser.CreateProcedureStmt:
		instrumentedAST := i.instrumentCreateProcedureAST(stmt, node)
		return i.codegen.GenerateSQL(instrumentedAST), nil
	default:
		// For other statements (DROP, etc.), return unchanged
		return stmt.Text, nil
	}
}

// instrumentCreateFunctionAST instruments a CREATE FUNCTION statement's AST for tracing
func (i *ASTInstrumenter) instrumentCreateFunctionAST(stmt sqlsplitter.Statement, node *sqlflowparser.CreateFunctionStmt) *sqlflowparser.CreateFunctionStmt {
	instrumentedBody := i.instrumentFunctionBody(stmt, node.Name, node.Parameters, node.Body, true, node.ReturnType)
	return &sqlflowparser.CreateFunctionStmt{
		BaseStatement: node.BaseStatement,
		Name:          node.Name,
		Parameters:    node.Parameters,
		ReturnType:    node.ReturnType,
		Body:          instrumentedBody,
	}
}

// instrumentCreateProcedureAST instruments a CREATE PROCEDURE statement's AST for tracing
func (i *ASTInstrumenter) instrumentCreateProcedureAST(stmt sqlsplitter.Statement, node *sqlflowparser.CreateProcedureStmt) *sqlflowparser.CreateProcedureStmt {
	instrumentedBody := i.instrumentFunctionBody(stmt, node.Name, node.Parameters, node.Body, false, "")
	return &sqlflowparser.CreateProcedureStmt{
		BaseStatement: node.BaseStatement,
		Name:          node.Name,
		Parameters:    node.Parameters,
		Body:          instrumentedBody,
	}
}

// instrumentFunctionBody instruments the body of a function or procedure for tracing
func (i *ASTInstrumenter) instrumentFunctionBody(stmt sqlsplitter.Statement, functionName string, params []sqlflowparser.Parameter, statements []sqlflowparser.StatementAST, isFunction bool, returnType string) []sqlflowparser.StatementAST {
	// If the body contains a single BEGIN statement, instrument inside it
	if len(statements) == 1 {
		if beginStmt, ok := statements[0].(*sqlflowparser.BeginStmt); ok {
			// Create instrumented BEGIN block with entry/exit tracing
			return []sqlflowparser.StatementAST{i.instrumentBeginBlock(beginStmt, stmt, functionName, params, isFunction, returnType)}
		}
	}

	// For other cases, add tracing around the statements
	var result []sqlflowparser.StatementAST

	// Add function entry tracing at the beginning
	entryCall := i.createEntryTracingCall(functionName, params, isFunction)
	result = append(result, entryCall)

	// Process each statement in the body
	for _, astStmt := range statements {
		// Add instrumented statement
		instrumentedStmt := i.instrumentSingleStatement(stmt, functionName, astStmt, isFunction, returnType)
		result = append(result, instrumentedStmt)
	}

	// For functions, add exit tracing before implicit return (at end)
	// For procedures, add exit tracing at the end
	if isFunction {
		// Check if the last statement is already a return statement
		hasExplicitReturn := false
		if len(statements) > 0 {
			if _, ok := statements[len(statements)-1].(*sqlflowparser.ReturnStmt); ok {
				hasExplicitReturn = true
			}
		}

		// If no explicit return, add exit tracing for implicit return
		if !hasExplicitReturn {
			exitCall := i.createExitTracingCall(functionName, "NULL", true, returnType)
			result = append(result, exitCall)
		}
	} else {
		// For procedures, always add exit tracing at the end with OUT parameters
		// Note: params may not be available in this context, so use basic exit call
		exitCall := i.createExitTracingCall(functionName, "NULL", false, "")
		result = append(result, exitCall)
	}

	return result
}

// instrumentBeginBlock instruments a BEGIN block for function/procedure tracing
func (i *ASTInstrumenter) instrumentBeginBlock(beginStmt *sqlflowparser.BeginStmt, stmt sqlsplitter.Statement, functionName string, params []sqlflowparser.Parameter, isFunction bool, returnType string) *sqlflowparser.BeginStmt {
	var instrumentedBody []sqlflowparser.StatementAST

	// First, add all DECLARE statements
	var declareStatements []sqlflowparser.StatementAST
	var executableStatements []sqlflowparser.StatementAST

	for _, astStmt := range beginStmt.Body {
		if _, ok := astStmt.(*sqlflowparser.DeclareStmt); ok {
			declareStatements = append(declareStatements, astStmt)
		} else {
			executableStatements = append(executableStatements, astStmt)
		}
	}

	// Add DECLARE statements first
	instrumentedBody = append(instrumentedBody, declareStatements...)

	// Add function entry tracing after DECLARE statements
	entryCall := i.createEntryTracingCall(functionName, params, isFunction)
	instrumentedBody = append(instrumentedBody, entryCall)

	// Process executable statements
	processedStatements := i.instrumentStatementList(stmt, functionName, executableStatements, isFunction, returnType)
	instrumentedBody = append(instrumentedBody, processedStatements...)

	// For functions, add exit tracing before implicit return (at end)
	// For procedures, add exit tracing at the end
	if isFunction {
		// Check if the last executable statement is already a return statement
		hasExplicitReturn := false
		if len(executableStatements) > 0 {
			if _, ok := executableStatements[len(executableStatements)-1].(*sqlflowparser.ReturnStmt); ok {
				hasExplicitReturn = true
			}
		}

		// If no explicit return, add exit tracing for implicit return
		if !hasExplicitReturn {
			exitCall := i.createExitTracingCall(functionName, "NULL", true, returnType)
			instrumentedBody = append(instrumentedBody, exitCall)
		}
	} else {
		// For procedures, always add exit tracing at the end with OUT parameters
		exitCall := i.createProcedureExitTracingCall(functionName, params)
		instrumentedBody = append(instrumentedBody, exitCall)
	}

	return &sqlflowparser.BeginStmt{
		BaseStatement: beginStmt.BaseStatement,
		Body:          instrumentedBody,
	}
}

// instrumentSingleStatement instruments a single statement for function tracing
func (i *ASTInstrumenter) instrumentSingleStatement(stmt sqlsplitter.Statement, functionName string, astStmt sqlflowparser.StatementAST, isFunction bool, returnType string) sqlflowparser.StatementAST {
	switch s := astStmt.(type) {
	case *sqlflowparser.BeginStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body, isFunction, returnType)
		return &sqlflowparser.BeginStmt{
			BaseStatement: s.BaseStatement,
			Body:          instrumentedBody,
		}
	case *sqlflowparser.IfStmt:
		instrumentedThen := i.instrumentStatementList(stmt, functionName, s.Then, isFunction, returnType)
		instrumentedElse := i.instrumentStatementList(stmt, functionName, s.Else, isFunction, returnType)

		var instrumentedElseIfs []sqlflowparser.ElseIfClause
		for _, elseif := range s.ElseIfs {
			instrumentedElseIfThen := i.instrumentStatementList(stmt, functionName, elseif.Then, isFunction, returnType)
			instrumentedElseIfs = append(instrumentedElseIfs, sqlflowparser.ElseIfClause{
				BaseStatement: elseif.BaseStatement,
				Condition:     elseif.Condition,
				Then:          instrumentedElseIfThen,
			})
		}

		return &sqlflowparser.IfStmt{
			BaseStatement: s.BaseStatement,
			Condition:     s.Condition,
			Then:          instrumentedThen,
			ElseIfs:       instrumentedElseIfs,
			Else:          instrumentedElse,
		}
	case *sqlflowparser.WhileStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body, isFunction, returnType)
		return &sqlflowparser.WhileStmt{
			BaseStatement: s.BaseStatement,
			Condition:     s.Condition,
			Body:          instrumentedBody,
		}
	case *sqlflowparser.LoopStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body, isFunction, returnType)
		return &sqlflowparser.LoopStmt{
			BaseStatement: s.BaseStatement,
			Body:          instrumentedBody,
		}
	case *sqlflowparser.RepeatStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body, isFunction, returnType)
		return &sqlflowparser.RepeatStmt{
			BaseStatement: s.BaseStatement,
			Body:          instrumentedBody,
			Condition:     s.Condition,
		}
	case *sqlflowparser.CaseStmt:
		var instrumentedWhenClauses []sqlflowparser.WhenClause
		for _, when := range s.WhenClauses {
			instrumentedWhenThen := i.instrumentStatementList(stmt, functionName, when.Then, isFunction, returnType)
			instrumentedWhenClauses = append(instrumentedWhenClauses, sqlflowparser.WhenClause{
				BaseStatement: when.BaseStatement,
				Condition:     when.Condition,
				Then:          instrumentedWhenThen,
			})
		}

		instrumentedElse := i.instrumentStatementList(stmt, functionName, s.Else, isFunction, returnType)

		return &sqlflowparser.CaseStmt{
			BaseStatement: s.BaseStatement,
			Expression:    s.Expression,
			WhenClauses:   instrumentedWhenClauses,
			Else:          instrumentedElse,
		}
	case *sqlflowparser.ReturnStmt:
		// Return statements are handled in the statement list processing
		// to avoid unnecessary BEGIN/END blocks
		return s
	default:
		// For statements without nested structure, return as-is
		return astStmt
	}
}

// instrumentStatementList instruments a list of statements for function tracing
func (i *ASTInstrumenter) instrumentStatementList(stmt sqlsplitter.Statement, functionName string, statements []sqlflowparser.StatementAST, isFunction bool, returnType string) []sqlflowparser.StatementAST {
	var result []sqlflowparser.StatementAST

	for _, astStmt := range statements {
		// Special handling for RETURN statements in functions - add exit tracing before the return
		if returnStmt, ok := astStmt.(*sqlflowparser.ReturnStmt); ok && isFunction {
			returnValue := i.extractReturnValue(returnStmt.Text)
			exitCall := i.createExitTracingCall(functionName, returnValue, true, returnType)
			result = append(result, exitCall)
			result = append(result, returnStmt)
		} else {
			// Add the instrumented statement
			instrumentedStmt := i.instrumentSingleStatement(stmt, functionName, astStmt, isFunction, returnType)
			result = append(result, instrumentedStmt)
		}
	}

	return result
}

// createEntryTracingCall creates a tracing call for function entry
func (i *ASTInstrumenter) createEntryTracingCall(functionName string, params []sqlflowparser.Parameter, isFunction bool) sqlflowparser.StatementAST {
	// Build arguments JSON using MySQL JSON_OBJECT function
	var argsExpr string
	if len(params) > 0 {
		var jsonParts []string
		for _, param := range params {
			jsonParts = append(jsonParts, fmt.Sprintf("'%s', %s", param.Name, param.Name))
		}
		argsExpr = fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(jsonParts, ", "))
	} else {
		argsExpr = "JSON_OBJECT()"
	}

	objectType := "procedure"
	if isFunction {
		objectType = "function"
	}

	callText := fmt.Sprintf("CALL __record_ftrace_entry('%s', '%s', '%s', %s)", i.filename, functionName, objectType, argsExpr)

	return &sqlflowparser.GenericStmt{
		BaseStatement: sqlflowparser.BaseStatement{
			Pos:   sqlflowparser.Position{Line: 1, Column: 1},
			Label: "",
			Text:  callText,
		},
	}
}

// createExitTracingCall creates a tracing call for function exit
func (i *ASTInstrumenter) createExitTracingCall(functionName string, returnValue string, isFunction bool, returnType string) sqlflowparser.StatementAST {
	objectType := "procedure"
	if isFunction {
		objectType = "function"
	}

	// Handle different data types based on function return type:
	// - For BLOB/BINARY types: encode as base64 and wrap in JSON_QUOTE
	// - For other types: use JSON_QUOTE directly
	var jsonReturnValue string
	if isBinaryType(returnType) {
		// For BLOB/BINARY types, always use base64 encoding
		jsonReturnValue = fmt.Sprintf("JSON_QUOTE(CASE WHEN %s IS NULL THEN 'NULL' ELSE CONCAT('base64:', TO_BASE64(%s)) END)", returnValue, returnValue)
	} else {
		// For other types, use JSON_QUOTE with safe string conversion
		jsonReturnValue = fmt.Sprintf("JSON_QUOTE(CASE WHEN %s IS NULL THEN 'NULL' ELSE CAST(%s AS CHAR) END)", returnValue, returnValue)
	}

	callText := fmt.Sprintf("CALL __record_ftrace_exit('%s', '%s', '%s', %s)", i.filename, functionName, objectType, jsonReturnValue)

	return &sqlflowparser.GenericStmt{
		BaseStatement: sqlflowparser.BaseStatement{
			Pos:   sqlflowparser.Position{Line: 1, Column: 1},
			Label: "",
			Text:  callText,
		},
	}
}

// createProcedureExitTracingCall creates a tracing call for procedure exit with OUT parameters
func (i *ASTInstrumenter) createProcedureExitTracingCall(functionName string, params []sqlflowparser.Parameter) sqlflowparser.StatementAST {
	// Build output parameters JSON using MySQL JSON_OBJECT function
	var outParts []string
	for _, param := range params {
		if param.Mode == "OUT" || param.Mode == "INOUT" {
			outParts = append(outParts, fmt.Sprintf("'%s', %s", param.Name, param.Name))
		}
	}

	var outExpr string
	if len(outParts) > 0 {
		outExpr = fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(outParts, ", "))
	} else {
		outExpr = "JSON_OBJECT()"
	}

	callText := fmt.Sprintf("CALL __record_ftrace_exit('%s', '%s', 'procedure', %s)", i.filename, functionName, outExpr)

	return &sqlflowparser.GenericStmt{
		BaseStatement: sqlflowparser.BaseStatement{
			Pos:   sqlflowparser.Position{Line: 1, Column: 1},
			Label: "",
			Text:  callText,
		},
	}
}

// extractReturnValue extracts the return value from a RETURN statement
func (i *ASTInstrumenter) extractReturnValue(returnText string) string {
	// Simple extraction - remove "RETURN " prefix and ";" suffix
	text := strings.TrimSpace(returnText)
	if strings.HasPrefix(strings.ToUpper(text), "RETURN ") {
		text = strings.TrimSpace(text[7:]) // Remove "RETURN "
	}
	if strings.HasSuffix(text, ";") {
		text = text[:len(text)-1] // Remove ";"
	}

	if text == "" {
		return "NULL"
	}

	// Wrap in COALESCE to handle NULL values safely
	return fmt.Sprintf("COALESCE(%s, 'NULL')", text)
}
