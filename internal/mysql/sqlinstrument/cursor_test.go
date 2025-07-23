package sqlinstrument

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestDeclareCursorNotInstrumented(t *testing.T) {
	g := NewWithT(t)

	input := `DELIMITER $$
CREATE PROCEDURE test_proc()
BEGIN
	DECLARE done INT DEFAULT FALSE;
	DECLARE element_cursor CURSOR FOR
		SELECT i, n, t, v
		FROM JSON_TABLE(
			JSON_EXTRACT(data, '$.*[*]'),
			'$[*]' COLUMNS (
				i INT PATH '$.i',
				n INT PATH '$.n'
			)
		) jt
		ORDER BY i;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
	SET message = 'test';
END$$
DELIMITER ;`

	instrumenter := NewInstrumenter("test.sql")
	result, err := instrumenter.InstrumentSQL([]byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	// Should not instrument the SELECT inside DECLARE CURSOR
	g.Expect(result).ToNot(ContainSubstring("CALL __record_coverage('test.sql', 'test_proc', 6); SELECT"))
	g.Expect(result).ToNot(ContainSubstring("CALL __record_coverage('test.sql', 'test_proc', 7); FROM"))

	// Should instrument the SET statement
	g.Expect(result).To(ContainSubstring("CALL __record_coverage('test.sql', 'test_proc', 16);"))
	g.Expect(result).To(ContainSubstring("SET message = 'test'"))

	// Print result for debugging
	t.Logf("Result:\n%s", result)
}
