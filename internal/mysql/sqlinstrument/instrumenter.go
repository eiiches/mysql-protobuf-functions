package sqlinstrument

// Instrumenter adds coverage instrumentation to SQL statements
type Instrumenter struct {
	astInstrumenter *ASTInstrumenter
}

// NewInstrumenter creates a new instrumenter for the given filename
func NewInstrumenter(filename string) *Instrumenter {
	return &Instrumenter{
		astInstrumenter: NewASTInstrumenter(filename),
	}
}

// InstrumentSQL adds coverage instrumentation to SQL content
func (i *Instrumenter) InstrumentSQL(content []byte) (string, error) {
	return i.astInstrumenter.InstrumentSQL(content)
}
