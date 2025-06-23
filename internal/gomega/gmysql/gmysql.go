package gmysql

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/onsi/gomega/types"
)

func BeMySQLError(number uint16, sqlState string, messageMatcher types.GomegaMatcher) types.GomegaMatcher {
	return &MySQLErrorMatcher{
		Number:         number,
		SQLState:       sqlState,
		MessageMatcher: messageMatcher,
	}
}

type MySQLErrorMatcher struct {
	Number         uint16
	SQLState       string
	MessageMatcher types.GomegaMatcher
}

func (m *MySQLErrorMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, nil
	}
	mysqlError, ok := actual.(*mysql.MySQLError)
	if !ok {
		return false, nil
	}
	if mysqlError.Number != m.Number {
		return false, nil
	}
	if string(mysqlError.SQLState[:]) != m.SQLState {
		return false, nil
	}
	if matches, err := m.MessageMatcher.Match(mysqlError.Message); err != nil {
		return false, fmt.Errorf("error matching message: %w", err)
	} else if !matches {
		return false, nil
	}
	return true, nil
}

func (m *MySQLErrorMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto match\n\t%#v", actual, m)
}

func (m *MySQLErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to match\n\t%#v", actual, m)
}

func (m *MySQLErrorMatcher) String() string {
	return fmt.Sprintf("MySQLErrorMatcher(Number: %d, SQLState: %s, MessageMatcher: %s)", m.Number, m.SQLState, m.MessageMatcher)
}
