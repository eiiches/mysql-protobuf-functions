package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/onsi/gomega/types"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/gomega"
)

var db *sql.DB

func TestMain(m *testing.M) {
	dataSourceName := flag.String("database", "", "Database connection string. Example: user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	flag.Parse()

	if *dataSourceName == "" {
		panic("Run with go test -args -database user:password@tcp(host:port)/dbname")
	}

	var err error
	db, err = sql.Open("mysql", *dataSourceName)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	m.Run()
}

type ExpressionTestContext struct {
	Expression string
	Args       []any
	T          *testing.T
}

func AssertThatExpression(t *testing.T, expression string, args ...any) *ExpressionTestContext {
	return &ExpressionTestContext{Expression: expression, Args: args, T: t}
}

func assertThatExpressionTo[T any](t *testing.T, matcher types.GomegaMatcher, expression string, args ...any) {
	g := NewWithT(t)

	stmt, err := db.Prepare("SELECT " + expression)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(stmt.Close()).To(Succeed())
	}()

	rows, err := stmt.Query(args...)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(rows.Close()).To(Succeed())
	}()

	if !rows.Next() { // no rows or an error
		g.Expect(rows.Err()).NotTo(HaveOccurred(), "Unexpected error while iterating over rows.")
		g.Fail("Expected one row, but got none.")
	}

	var actual T
	g.Expect(rows.Scan(&actual)).To(Succeed())

	g.Expect(actual).To(matcher)
}

func assertThatExpressionToFailWith(t *testing.T, matcher types.GomegaMatcher, expression string, args ...any) {
	g := NewWithT(t)

	stmt, err := db.Prepare("SELECT " + expression)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(stmt.Close()).To(Succeed())
	}()

	rows, err := stmt.Query(args...)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(rows.Close()).To(Succeed())
	}()

	if !rows.Next() { // no rows or an error
		if rows.Err() == nil {
			g.Fail("Expected an error, but got no rows.")
		} else {
			g.Expect(rows.Err()).To(matcher)
			return
		}
	}

	// If we got here, it means we got a row, which is unexpected.

	var actual interface{}
	g.Expect(rows.Scan(&actual)).To(Succeed())

	g.Fail(fmt.Sprintf("Expected an error, but got a row with value: %v", actual))
}

func runAssertThatExpressionIsEqualTo[T any](t *testing.T, method string, expected T, expression string, args ...any) {
	t.Run(fmt.Sprintf("AssertThatExpression(`%s`, %v...).%s(%v)", expression, args, method, expected), func(t *testing.T) {
		assertThatExpressionTo[T](t, Equal(expected), expression, args...)
	})
}

func (this *ExpressionTestContext) IsEqualToUint(expected uint64) {
	runAssertThatExpressionIsEqualTo[uint64](this.T, "IsEqualToUint", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToInt(expected int64) {
	runAssertThatExpressionIsEqualTo[int64](this.T, "IsEqualToInt", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToDouble(expected float64) {
	runAssertThatExpressionIsEqualTo[float64](this.T, "IsEqualToDouble", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualTo(expected interface{}) {
	runAssertThatExpressionIsEqualTo[interface{}](this.T, "IsEqualTo", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToString(expected string) {
	runAssertThatExpressionIsEqualTo[string](this.T, "IsEqualToString", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToBytes(expected []byte) {
	runAssertThatExpressionIsEqualTo[[]byte](this.T, "IsEqualToBytes", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToBool(expected bool) {
	runAssertThatExpressionIsEqualTo[bool](this.T, "IsEqualToBool", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsTrue() {
	this.T.Run(fmt.Sprintf("AssertThatExpression(`%s`, %v...).IsTrue()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[bool](this.T, BeTrue(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsFalse() {
	this.T.Run(fmt.Sprintf("AssertThatExpression(`%s`, %v...).IsFalse()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[bool](this.T, BeFalse(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsNull() {
	this.T.Run(fmt.Sprintf("AssertThatExpression(`%s`, %v...).IsNull()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[interface{}](this.T, BeNil(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) ToSucceed() {
	this.T.Run(fmt.Sprintf("AssertThatExpression(`%s`, %v...).ToSucceed()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[interface{}](t, SatisfyAny(BeNil(), Not(BeNil())), this.Expression, this.Args...)
	})
}

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

func (this *ExpressionTestContext) ToFailWithSignalException(state string, containsMessage string) {
	this.T.Run(fmt.Sprintf("AssertThatExpression(`%s`, %v...).ToFailWithSignalException(%v, %v)", this.Expression, this.Args, state, containsMessage), func(t *testing.T) {
		assertThatExpressionToFailWith(t, BeMySQLError(1644, state, ContainSubstring(containsMessage)), this.Expression, this.Args...)
	})
}
