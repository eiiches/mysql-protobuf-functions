package sqlftrace

import (
	"fmt"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
)

// CodeGenerator converts AST nodes back to SQL text for function tracing
type CodeGenerator struct {
	indent string
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		indent: "\t",
	}
}

// GenerateSQL converts an AST node to SQL text
func (cg *CodeGenerator) GenerateSQL(node sqlflowparser.AST) string {
	return cg.generateStatement(node, 0)
}

// generateStatement generates SQL for a statement AST node
func (cg *CodeGenerator) generateStatement(node sqlflowparser.AST, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	switch stmt := node.(type) {
	case *sqlflowparser.CreateFunctionStmt:
		return cg.generateCreateFunction(*stmt, indentLevel)
	case *sqlflowparser.CreateProcedureStmt:
		return cg.generateCreateProcedure(*stmt, indentLevel)
	case *sqlflowparser.BeginStmt:
		return cg.generateBegin(*stmt, indentLevel)
	case *sqlflowparser.IfStmt:
		return cg.generateIf(*stmt, indentLevel)
	case *sqlflowparser.WhileStmt:
		return cg.generateWhile(*stmt, indentLevel)
	case *sqlflowparser.LoopStmt:
		return cg.generateLoop(*stmt, indentLevel)
	case *sqlflowparser.RepeatStmt:
		return cg.generateRepeat(*stmt, indentLevel)
	case *sqlflowparser.CaseStmt:
		return cg.generateCase(*stmt, indentLevel)
	case *sqlflowparser.ReturnStmt:
		result := cg.addLabelPrefix(stmt, stmt.Text)
		if !strings.HasSuffix(result, ";") {
			result += ";"
		}
		return indent + result
	case *sqlflowparser.LeaveStmt:
		// LEAVE statements should not have their own labels, they reference other labels
		// The stmt.Label field should be empty, and we use stmt.Text for the target label
		if !strings.HasSuffix(stmt.Text, ";") {
			return indent + stmt.Text + ";"
		}
		return indent + stmt.Text
	case *sqlflowparser.IterateStmt:
		// ITERATE statements should not have their own labels, they reference other labels
		// The stmt.Label field should be empty, and we use stmt.Text for the target label
		if !strings.HasSuffix(stmt.Text, ";") {
			return indent + stmt.Text + ";"
		}
		return indent + stmt.Text
	case *sqlflowparser.DeclareStmt:
		result := cg.addLabelPrefix(stmt, stmt.Text)
		if !strings.HasSuffix(result, ";") {
			result += ";"
		}
		return indent + result
	case *sqlflowparser.GenericStmt:
		result := cg.addLabelPrefix(stmt, stmt.Text)
		if !strings.HasSuffix(result, ";") {
			result += ";"
		}
		return indent + result
	default:
		return indent + "-- Unknown statement type"
	}
}

// addLabelPrefix adds a label prefix to a statement if it has one
func (cg *CodeGenerator) addLabelPrefix(stmt sqlflowparser.StatementAST, content string) string {
	if label := stmt.GetLabel(); label != "" {
		return label + ": " + content
	}
	return content
}

// generateCreateFunction generates CREATE FUNCTION statement
func (cg *CodeGenerator) generateCreateFunction(stmt sqlflowparser.CreateFunctionStmt, indentLevel int) string {
	var params []string
	for _, param := range stmt.Parameters {
		paramStr := param.Name + " " + param.Type
		if param.Mode != "" {
			paramStr = param.Mode + " " + paramStr
		}
		params = append(params, paramStr)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("CREATE FUNCTION %s(%s) RETURNS %s DETERMINISTIC\n",
		stmt.Name, strings.Join(params, ", "), stmt.ReturnType))

	for _, bodyStmt := range stmt.Body {
		result.WriteString(cg.generateStatement(bodyStmt, indentLevel))
		result.WriteString("\n")
	}

	return strings.TrimRight(result.String(), "\n")
}

// generateCreateProcedure generates CREATE PROCEDURE statement
func (cg *CodeGenerator) generateCreateProcedure(stmt sqlflowparser.CreateProcedureStmt, indentLevel int) string {
	var params []string
	for _, param := range stmt.Parameters {
		paramStr := param.Name + " " + param.Type
		if param.Mode != "" {
			paramStr = param.Mode + " " + paramStr
		}
		params = append(params, paramStr)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("CREATE PROCEDURE %s(%s)\n",
		stmt.Name, strings.Join(params, ", ")))

	for _, bodyStmt := range stmt.Body {
		result.WriteString(cg.generateStatement(bodyStmt, indentLevel))
		result.WriteString("\n")
	}

	return strings.TrimRight(result.String(), "\n")
}

// generateBegin generates BEGIN...END block
func (cg *CodeGenerator) generateBegin(stmt sqlflowparser.BeginStmt, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	var result strings.Builder
	if stmt.GetLabel() != "" {
		result.WriteString(indent + stmt.GetLabel() + ": ")
	}
	result.WriteString("BEGIN\n")

	for _, bodyStmt := range stmt.Body {
		result.WriteString(cg.generateStatement(bodyStmt, indentLevel+1))
		result.WriteString("\n")
	}

	result.WriteString(indent + "END")
	return result.String()
}

// generateIf generates IF statement
func (cg *CodeGenerator) generateIf(stmt sqlflowparser.IfStmt, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	var result strings.Builder
	if stmt.GetLabel() != "" {
		result.WriteString(indent + stmt.GetLabel() + ": ")
	} else {
		result.WriteString(indent)
	}
	result.WriteString("IF " + stmt.Condition + " THEN\n")

	for _, thenStmt := range stmt.Then {
		result.WriteString(cg.generateStatement(thenStmt, indentLevel+1))
		result.WriteString("\n")
	}

	for _, elseif := range stmt.ElseIfs {
		result.WriteString(indent + "ELSEIF " + elseif.Condition + " THEN\n")
		for _, elseifStmt := range elseif.Then {
			result.WriteString(cg.generateStatement(elseifStmt, indentLevel+1))
			result.WriteString("\n")
		}
	}

	if len(stmt.Else) > 0 {
		result.WriteString(indent + "ELSE\n")
		for _, elseStmt := range stmt.Else {
			result.WriteString(cg.generateStatement(elseStmt, indentLevel+1))
			result.WriteString("\n")
		}
	}

	result.WriteString(indent + "END IF;")
	return result.String()
}

// generateWhile generates WHILE statement
func (cg *CodeGenerator) generateWhile(stmt sqlflowparser.WhileStmt, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	var result strings.Builder
	if stmt.GetLabel() != "" {
		result.WriteString(indent + stmt.GetLabel() + ": ")
	}
	result.WriteString("WHILE " + stmt.Condition + " DO\n")

	for _, bodyStmt := range stmt.Body {
		result.WriteString(cg.generateStatement(bodyStmt, indentLevel+1))
		result.WriteString("\n")
	}

	result.WriteString(indent + "END WHILE")
	if stmt.GetLabel() != "" {
		result.WriteString(" " + stmt.GetLabel())
	}
	result.WriteString(";")
	return result.String()
}

// generateLoop generates LOOP statement
func (cg *CodeGenerator) generateLoop(stmt sqlflowparser.LoopStmt, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	var result strings.Builder
	if stmt.Label != "" {
		result.WriteString(indent + stmt.Label + ": ")
	}
	result.WriteString("LOOP\n")

	for _, bodyStmt := range stmt.Body {
		result.WriteString(cg.generateStatement(bodyStmt, indentLevel+1))
		result.WriteString("\n")
	}

	result.WriteString(indent + "END LOOP")
	if stmt.Label != "" {
		result.WriteString(" " + stmt.Label)
	}
	result.WriteString(";")
	return result.String()
}

// generateRepeat generates REPEAT statement
func (cg *CodeGenerator) generateRepeat(stmt sqlflowparser.RepeatStmt, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	var result strings.Builder
	if stmt.GetLabel() != "" {
		result.WriteString(indent + stmt.GetLabel() + ": ")
	}
	result.WriteString("REPEAT\n")

	for _, bodyStmt := range stmt.Body {
		result.WriteString(cg.generateStatement(bodyStmt, indentLevel+1))
		result.WriteString("\n")
	}

	result.WriteString(indent + "UNTIL " + stmt.Condition + " END REPEAT;")
	return result.String()
}

// generateCase generates CASE statement
func (cg *CodeGenerator) generateCase(stmt sqlflowparser.CaseStmt, indentLevel int) string {
	indent := strings.Repeat(cg.indent, indentLevel)

	var result strings.Builder
	if stmt.GetLabel() != "" {
		result.WriteString(indent + stmt.GetLabel() + ": ")
	}
	if stmt.Expression != "" {
		result.WriteString("CASE " + stmt.Expression + "\n")
	} else {
		result.WriteString("CASE\n")
	}

	for _, when := range stmt.WhenClauses {
		result.WriteString(indent + cg.indent + "WHEN " + when.Condition + " THEN\n")
		for _, thenStmt := range when.Then {
			result.WriteString(cg.generateStatement(thenStmt, indentLevel+2))
			result.WriteString("\n")
		}
	}

	if len(stmt.Else) > 0 {
		result.WriteString(indent + cg.indent + "ELSE\n")
		for _, elseStmt := range stmt.Else {
			result.WriteString(cg.generateStatement(elseStmt, indentLevel+2))
			result.WriteString("\n")
		}
	}

	result.WriteString(indent + "END CASE;")
	return result.String()
}
