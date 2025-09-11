package sqlftrace

// Instrumenter adds function tracing instrumentation to SQL statements
type Instrumenter struct {
	astInstrumenter *ASTInstrumenter
}

// NewInstrumenter creates a new instrumenter for the given filename
func NewInstrumenter(filename string) *Instrumenter {
	return &Instrumenter{
		astInstrumenter: NewASTInstrumenter(filename),
	}
}

// SetStatementTracing enables or disables statement-level tracing
func (i *Instrumenter) SetStatementTracing(enabled bool) {
	i.astInstrumenter.SetStatementTracing(enabled)
}

// InstrumentSQL adds function tracing instrumentation to SQL content
func (i *Instrumenter) InstrumentSQL(content []byte) (string, error) {
	return i.astInstrumenter.InstrumentSQL(content)
}
