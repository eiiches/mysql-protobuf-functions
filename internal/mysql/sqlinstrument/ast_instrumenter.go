package sqlinstrument

import (
	"fmt"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlsplitter"
)

// ASTInstrumenter adds coverage instrumentation using AST reconstruction
type ASTInstrumenter struct {
	filename string
	codegen  *CodeGenerator
}

// NewASTInstrumenter creates a new AST-based instrumenter
func NewASTInstrumenter(filename string) *ASTInstrumenter {
	return &ASTInstrumenter{
		filename: filename,
		codegen:  NewCodeGenerator(),
	}
}

// InstrumentSQL adds coverage instrumentation to SQL content using AST reconstruction
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

// instrumentStatementWithAST instruments a statement using AST reconstruction
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

// instrumentCreateFunctionAST instruments a CREATE FUNCTION statement's AST
func (i *ASTInstrumenter) instrumentCreateFunctionAST(stmt sqlsplitter.Statement, node *sqlflowparser.CreateFunctionStmt) *sqlflowparser.CreateFunctionStmt {
	instrumentedBody := i.instrumentStatementList(stmt, node.Name, node.Body)
	return &sqlflowparser.CreateFunctionStmt{
		BaseStatement: node.BaseStatement,
		Name:          node.Name,
		Parameters:    node.Parameters,
		ReturnType:    node.ReturnType,
		Body:          instrumentedBody,
	}
}

// instrumentCreateProcedureAST instruments a CREATE PROCEDURE statement's AST
func (i *ASTInstrumenter) instrumentCreateProcedureAST(stmt sqlsplitter.Statement, node *sqlflowparser.CreateProcedureStmt) *sqlflowparser.CreateProcedureStmt {
	instrumentedBody := i.instrumentStatementList(stmt, node.Name, node.Body)
	return &sqlflowparser.CreateProcedureStmt{
		BaseStatement: node.BaseStatement,
		Name:          node.Name,
		Parameters:    node.Parameters,
		Body:          instrumentedBody,
	}
}

// instrumentStatementList instruments a list of statements
func (i *ASTInstrumenter) instrumentStatementList(stmt sqlsplitter.Statement, functionName string, statements []sqlflowparser.StatementAST) []sqlflowparser.StatementAST {
	var result []sqlflowparser.StatementAST

	for _, astStmt := range statements {
		// Add the original statement (potentially modified)
		instrumentedStmt := i.instrumentSingleStatement(stmt, functionName, astStmt)

		// If this statement should be instrumented, add a coverage call before it
		if i.shouldInstrumentStatement(astStmt) {
			originalLineNo := stmt.LineNo + astStmt.GetPosition().Line - 1
			coverageCallText := fmt.Sprintf("CALL __record_coverage('%s', '%s', %d)", i.filename, functionName, originalLineNo)
			coverageCall := &sqlflowparser.GenericStmt{
				BaseStatement: sqlflowparser.BaseStatement{
					Pos:   astStmt.GetPosition(),
					Label: "",
					Text:  coverageCallText,
				},
			}
			result = append(result, coverageCall)
		}

		result = append(result, instrumentedStmt)
	}

	return result
}

// instrumentSingleStatement instruments a single statement
func (i *ASTInstrumenter) instrumentSingleStatement(stmt sqlsplitter.Statement, functionName string, astStmt sqlflowparser.StatementAST) sqlflowparser.StatementAST {
	switch s := astStmt.(type) {
	case *sqlflowparser.BeginStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body)
		return &sqlflowparser.BeginStmt{
			BaseStatement: s.BaseStatement,
			Body:          instrumentedBody,
		}
	case *sqlflowparser.IfStmt:
		instrumentedThen := i.instrumentStatementList(stmt, functionName, s.Then)
		instrumentedElse := i.instrumentStatementList(stmt, functionName, s.Else)

		var instrumentedElseIfs []sqlflowparser.ElseIfClause
		for _, elseif := range s.ElseIfs {
			instrumentedElseIfThen := i.instrumentStatementList(stmt, functionName, elseif.Then)
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
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body)
		return &sqlflowparser.WhileStmt{
			BaseStatement: s.BaseStatement,
			Condition:     s.Condition,
			Body:          instrumentedBody,
		}
	case *sqlflowparser.LoopStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body)
		return &sqlflowparser.LoopStmt{
			BaseStatement: s.BaseStatement,
			Body:          instrumentedBody,
		}
	case *sqlflowparser.RepeatStmt:
		instrumentedBody := i.instrumentStatementList(stmt, functionName, s.Body)
		return &sqlflowparser.RepeatStmt{
			BaseStatement: s.BaseStatement,
			Body:          instrumentedBody,
			Condition:     s.Condition,
		}
	case *sqlflowparser.CaseStmt:
		var instrumentedWhenClauses []sqlflowparser.WhenClause
		for _, when := range s.WhenClauses {
			instrumentedWhenThen := i.instrumentStatementList(stmt, functionName, when.Then)
			instrumentedWhenClauses = append(instrumentedWhenClauses, sqlflowparser.WhenClause{
				BaseStatement: when.BaseStatement,
				Condition:     when.Condition,
				Then:          instrumentedWhenThen,
			})
		}

		instrumentedElse := i.instrumentStatementList(stmt, functionName, s.Else)

		return &sqlflowparser.CaseStmt{
			BaseStatement: s.BaseStatement,
			Expression:    s.Expression,
			WhenClauses:   instrumentedWhenClauses,
			Else:          instrumentedElse,
		}
	default:
		// For statements without nested structure, return as-is
		return astStmt
	}
}

// shouldInstrumentStatement determines if an AST statement should be instrumented
func (i *ASTInstrumenter) shouldInstrumentStatement(stmt sqlflowparser.StatementAST) bool {
	switch stmt.(type) {
	case *sqlflowparser.DeclareStmt:
		return false
	case *sqlflowparser.GenericStmt:
		return true
	case *sqlflowparser.IfStmt, *sqlflowparser.WhileStmt, *sqlflowparser.LoopStmt,
		*sqlflowparser.RepeatStmt, *sqlflowparser.CaseStmt, *sqlflowparser.ReturnStmt,
		*sqlflowparser.LeaveStmt, *sqlflowparser.IterateStmt:
		// Control flow statements should be instrumented
		return true
	default:
		// Don't instrument other statement types (BeginStmt, etc.)
		return false
	}
}
