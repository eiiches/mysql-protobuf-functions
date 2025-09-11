package sqlflowparser

// Position represents the position of a statement in the source code
type Position struct {
	Line   int // Line number (1-based)
	Column int // Column number (1-based)
	Offset int // Byte offset in the source
}

// AST represents the Abstract Syntax Tree for MySQL statements
type AST interface {
	astNode()
	GetPosition() Position
}

// Base AST node types
type (
	// StatementAST represents any SQL statement
	StatementAST interface {
		AST
		statementNode()

		// SetLabel sets the label for the statement
		SetLabel(label string)
		// GetLabel gets the label for the statement
		GetLabel() string
	}

	// BaseStatement contains common fields for all statements
	BaseStatement struct {
		Pos   Position
		Label string
		Text  string // Optional text content for statements that need it
	}
)

// Statement types
type (
	// CreateProcedureStmt represents CREATE PROCEDURE statement
	CreateProcedureStmt struct {
		BaseStatement
		Name       string
		Parameters []Parameter
		Body       []StatementAST
	}

	// CreateFunctionStmt represents CREATE FUNCTION statement
	CreateFunctionStmt struct {
		BaseStatement
		Name       string
		Parameters []Parameter
		ReturnType string
		Body       []StatementAST
	}

	// IfStmt represents IF statement
	IfStmt struct {
		BaseStatement
		Condition string
		Then      []StatementAST
		ElseIfs   []ElseIfClause
		Else      []StatementAST
	}

	// ElseIfClause represents ELSEIF clause
	ElseIfClause struct {
		BaseStatement
		Condition string
		Then      []StatementAST
	}

	// WhileStmt represents WHILE statement
	WhileStmt struct {
		BaseStatement
		Condition string
		Body      []StatementAST
	}

	// LoopStmt represents LOOP statement
	LoopStmt struct {
		BaseStatement
		Body []StatementAST
	}

	// RepeatStmt represents REPEAT statement
	RepeatStmt struct {
		BaseStatement
		Body      []StatementAST
		Condition string
	}

	// CaseStmt represents CASE statement
	CaseStmt struct {
		BaseStatement
		Expression  string
		WhenClauses []WhenClause
		Else        []StatementAST
	}

	// WhenClause represents WHEN clause in CASE statement
	WhenClause struct {
		BaseStatement
		Condition string
		Then      []StatementAST
	}

	// BeginStmt represents BEGIN...END block
	BeginStmt struct {
		BaseStatement
		Body []StatementAST
	}

	// LeaveStmt represents LEAVE statement
	LeaveStmt struct {
		BaseStatement
		// Label is inherited from BaseStatement
	}

	// IterateStmt represents ITERATE statement
	IterateStmt struct {
		BaseStatement
		// Label is inherited from BaseStatement
	}

	// ReturnStmt represents RETURN statement
	ReturnStmt struct {
		BaseStatement
		// Text is inherited from BaseStatement
	}

	// DeclareStmt represents DECLARE statement
	DeclareStmt struct {
		BaseStatement
		// Text is inherited from BaseStatement
	}

	// GenericStmt represents any other SQL statement (SELECT, INSERT, etc.)
	GenericStmt struct {
		BaseStatement
		// Text is inherited from BaseStatement
	}

	// SetVariableStmt represents SET variable statements (user variables, system variables)
	SetVariableStmt struct {
		BaseStatement
		Assignments []VariableAssignment // Multiple variables can be set in one statement
	}

	// Parameter represents procedure/function parameter
	Parameter struct {
		Name string
		Type string
		Mode string // IN, OUT, INOUT
	}

	// VariableAssignment represents a single variable assignment within SET
	VariableAssignment struct {
		// Explicit scope keyword (GLOBAL, SESSION, PERSIST, PERSIST_ONLY) or empty string
		ScopeKeyword string

		// Variable reference as written (e.g., "var", "@var", "@@var", "@@GLOBAL.var")
		VariableRef string

		// Assignment operator (usually "=", but could be ":=" for user variables)
		Operator string

		// Right-hand side expression
		Value string
	}
)

// Implement AST interface
func (BaseStatement) astNode() {}

// Implement StatementAST interface
func (BaseStatement) statementNode() {}

// Implement GetPosition method for all AST nodes
func (s BaseStatement) GetPosition() Position { return s.Pos }

func (s *BaseStatement) SetLabel(label string) {
	s.Label = label
}

func (s *BaseStatement) GetLabel() string {
	return s.Label
}
