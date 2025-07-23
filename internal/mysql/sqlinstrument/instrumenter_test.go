package sqlinstrument

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestInstrumentSimpleFunction(t *testing.T) {
	g := NewWithT(t)

	input := `DELIMITER $$
CREATE FUNCTION test_func(x INT) RETURNS INT DETERMINISTIC
BEGIN
	RETURN x * 2;
END$$
DELIMITER ;`

	instrumenter := NewInstrumenter("test.sql")
	result, err := instrumenter.InstrumentSQL([]byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 4);"))
	g.Expect(result).To(ContainSubstring("RETURN x * 2"))
}

func TestInstrumentFunctionWithIfStatement(t *testing.T) {
	g := NewWithT(t)

	input := `DELIMITER $$
CREATE FUNCTION test_func(x INT) RETURNS INT DETERMINISTIC
BEGIN
	IF x > 0 THEN
		RETURN x;
	ELSE
		RETURN 0;
	END IF;
END$$
DELIMITER ;`

	instrumenter := NewInstrumenter("test.sql")
	result, err := instrumenter.InstrumentSQL([]byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 4);"))
	g.Expect(result).To(ContainSubstring("IF x > 0 THEN"))
	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 5);"))
	g.Expect(result).To(ContainSubstring("RETURN x"))
	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 7);"))
	g.Expect(result).To(ContainSubstring("RETURN 0"))
}

func TestInstrumentProcedure(t *testing.T) {
	g := NewWithT(t)

	input := `DELIMITER $$
CREATE PROCEDURE test_proc(IN x INT, OUT result INT)
BEGIN
	SET result = x * 2;
END$$
DELIMITER ;`

	instrumenter := NewInstrumenter("test.sql")
	result, err := instrumenter.InstrumentSQL([]byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_proc', 4);"))
	g.Expect(result).To(ContainSubstring("SET result = x * 2"))
}

func TestInstrumentDelimiterStatements(t *testing.T) {
	g := NewWithT(t)

	input := `DELIMITER $$
CREATE FUNCTION test_func(x INT) RETURNS INT DETERMINISTIC
BEGIN
	RETURN x;
END$$
DELIMITER ;`

	instrumenter := NewInstrumenter("test.sql")
	result, err := instrumenter.InstrumentSQL([]byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(result).To(ContainSubstring("DELIMITER $$"))
	g.Expect(result).To(ContainSubstring("DELIMITER ;"))
	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 4);"))
	g.Expect(result).To(ContainSubstring("RETURN x"))
}

func TestSkipNonExecutableLines(t *testing.T) {
	g := NewWithT(t)

	input := `DELIMITER $$
CREATE FUNCTION test_func(x INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE temp INT;
	-- This is a comment
	SET temp = x;
	RETURN temp;
END$$
DELIMITER ;`

	instrumenter := NewInstrumenter("test.sql")
	result, err := instrumenter.InstrumentSQL([]byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	// Should not instrument DECLARE or comments
	g.Expect(result).ToNot(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 4); DECLARE"))
	g.Expect(result).ToNot(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 5); --"))

	// Should instrument SET and RETURN
	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 6);"))
	g.Expect(result).To(ContainSubstring("SET temp = x"))
	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_func', 7);"))
	g.Expect(result).To(ContainSubstring("RETURN temp"))
}
