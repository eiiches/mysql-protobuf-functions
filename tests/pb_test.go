package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gmysql"
	"github.com/eiiches/mysql-protobuf-functions/internal/moresql"
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

type CallTestContext struct {
	CallExpression string
	Args           []any
	T              *testing.T
}

func AssertThatCall(t *testing.T, callExpression string, args ...any) *CallTestContext {
	return &CallTestContext{CallExpression: callExpression, Args: args, T: t}
}

func (this *CallTestContext) ShouldSucceed() {
	g := NewWithT(this.T)

	stmt, err := db.Prepare("CALL " + this.CallExpression)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(stmt.Close()).To(Succeed())
	}()

	rows, err := stmt.Query(this.Args...)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(rows.Close()).To(Succeed())
	}()

	for rows.Next() {
		// do nothing
	}

	g.Expect(rows.Err()).NotTo(HaveOccurred(), "Unexpected error while iterating over rows.")
}

type ExpressionTestContext struct {
	Expression string
	Args       []any
	T          *testing.T
	RunFn      func(name string, testFn func(t *testing.T)) bool
}

func AssertThatExpression(t *testing.T, expression string, args ...any) *ExpressionTestContext {
	return &ExpressionTestContext{Expression: expression, Args: args, T: t, RunFn: func(name string, testFn func(t *testing.T)) bool {
		testFn(t)
		return true // unused
	}}
}

func RunTestThatExpression(t *testing.T, expression string, args ...any) *ExpressionTestContext {
	return AssertThatExpression(t, expression, args...).AsSubTest()
}

func (this *ExpressionTestContext) AsSubTest() *ExpressionTestContext {
	return &ExpressionTestContext{
		Expression: this.Expression,
		Args:       this.Args,
		T:          this.T,
		RunFn:      this.T.Run,
	}
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

func runAssertThatExpressionIsEqualTo[T any](runFn func(name string, fn func(t *testing.T)) bool, method string, expected T, expression string, args ...any) {
	runFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).%s(%v)", expression, args, method, expected), func(t *testing.T) {
		assertThatExpressionTo[T](t, Equal(expected), expression, args...)
	})
}

func (this *ExpressionTestContext) IsEqualToUint(expected uint64) {
	runAssertThatExpressionIsEqualTo[uint64](this.RunFn, "IsEqualToUint", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToInt(expected int64) {
	runAssertThatExpressionIsEqualTo[int64](this.RunFn, "IsEqualToInt", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToDouble(expected float64) {
	runAssertThatExpressionIsEqualTo[float64](this.RunFn, "IsEqualToDouble", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToFloat(expected float32) {
	runAssertThatExpressionIsEqualTo[float32](this.RunFn, "IsEqualToFloat", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualTo(expected interface{}) {
	runAssertThatExpressionIsEqualTo[interface{}](this.RunFn, "IsEqualTo", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToString(expected string) {
	runAssertThatExpressionIsEqualTo[string](this.RunFn, "IsEqualToString", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToBytes(expected []byte) {
	runAssertThatExpressionIsEqualTo[[]byte](this.RunFn, "IsEqualToBytes", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsEqualToBool(expected bool) {
	runAssertThatExpressionIsEqualTo[bool](this.RunFn, "IsEqualToBool", expected, this.Expression, this.Args...)
}

func (this *ExpressionTestContext) IsTrue() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).IsTrue()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[bool](this.T, BeTrue(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsFalse() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).IsFalse()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[bool](this.T, BeFalse(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsNull() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).IsNull()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[interface{}](this.T, BeNil(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsEqualToJson(expectedJson string) {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).IsEqualToJson(%+v)", this.Expression, this.Args, expectedJson), func(t *testing.T) {
		assertThatExpressionTo[string](t, MatchJSON(expectedJson), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) ToSucceed() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).ToSucceed()", this.Expression, this.Args), func(t *testing.T) {
		assertThatExpressionTo[interface{}](t, SatisfyAny(BeNil(), Not(BeNil())), this.Expression, this.Args...)
	})
}
func (this *ExpressionTestContext) ToFailWithSignalException(state string, containsMessage string) {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`, %v...).ToFailWithSignalException(%v, %v)", this.Expression, this.Args, state, containsMessage), func(t *testing.T) {
		assertThatExpressionToFailWith(t, gmysql.BeMySQLError(1644, state, ContainSubstring(containsMessage)), this.Expression, this.Args...)
	})
}

type StatementTestContext[T any] struct {
	Statement string
	Args      []any
	T         *testing.T
	RunFn     func(name string, testFn func(t *testing.T)) bool
}

func AssertThatStatement[T any](t *testing.T, statement string, args ...any) *StatementTestContext[T] {
	return &StatementTestContext[T]{Statement: statement, Args: args, T: t, RunFn: func(name string, testFn func(t *testing.T)) bool {
		testFn(t)
		return true // unused
	}}
}

func RunTestThatStatement[T any](t *testing.T, statement string, args ...any) *StatementTestContext[T] {
	return AssertThatStatement[T](t, statement, args...).AsSubTest()
}

func (this *StatementTestContext[T]) AsSubTest() *StatementTestContext[T] {
	return &StatementTestContext[T]{
		Statement: this.Statement,
		Args:      this.Args,
		T:         this.T,
		RunFn:     this.T.Run,
	}
}

func (this *StatementTestContext[T]) ShouldReturnSingleRow(matcher types.GomegaMatcher) {
	this.RunFn(fmt.Sprintf("RunTestThatStatement(`%s`, %v...).ShouldReturnSingleRow(%v)", this.Statement, this.Args, matcher), func(t *testing.T) {
		g := NewWithT(this.T)

		stmt, err := db.Prepare(this.Statement)
		g.Expect(err).NotTo(HaveOccurred())
		defer func() {
			g.Expect(stmt.Close()).To(Succeed())
		}()

		rows, err := stmt.Query(this.Args...)
		g.Expect(err).NotTo(HaveOccurred())
		defer func() {
			g.Expect(rows.Close()).To(Succeed())
		}()

		if !rows.Next() { // no rows or an error
			g.Expect(rows.Err()).NotTo(HaveOccurred(), "Unexpected error while iterating over rows.")
			g.Fail("Expected one row, but got none.")
		}

		row, err := moresql.ScanStruct[T](rows)
		g.Expect(err).NotTo(HaveOccurred(), "Unexpected error while scanning row.")
		g.Expect(row).To(matcher)

		g.Expect(rows.Next()).To(BeFalse(), "Expected no more rows, but got another row.")
		g.Expect(rows.Err()).NotTo(HaveOccurred(), "Unexpected error after scanning a row.")
	})
}
