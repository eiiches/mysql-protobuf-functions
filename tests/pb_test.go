package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gfloat"
	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gmysql"
	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gproto"
	"github.com/eiiches/mysql-protobuf-functions/internal/moresql"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	db         *sql.DB
	iterations = 100
)

func TestMain(m *testing.M) {
	dataSourceName := flag.String("database", "", "Database connection string. Example: user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	flag.IntVar(&iterations, "fuzz-iterations", 100, "Number of iterations for fuzz/random testing")
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

func marshalArgs(args []any) []any {
	marshaledArgs := make([]any, len(args))
	for i, arg := range args {
		if protoMsg, ok := arg.(proto.Message); ok {
			protoBytes, err := proto.Marshal(protoMsg)
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal proto message: %v", err))
			}
			marshaledArgs[i] = protoBytes
		} else {
			marshaledArgs[i] = arg
		}
	}
	return marshaledArgs
}

func formatArguments(args ...any) string {
	var formatted []string
	for _, arg := range args {
		switch v := arg.(type) {
		case proto.Message:
			fullName := string(v.ProtoReflect().Descriptor().FullName())
			jsonValue, err := protojson.MarshalOptions{Indent: "", Multiline: false}.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal proto message to JSON: %v", err))
			}
			formatted = append(formatted, fmt.Sprintf("&%s%s", fullName, string(jsonValue)))
		case []byte:
			formatted = append(formatted, "0x"+hex.EncodeToString(v))
		case protoreflect.Name:
			formatted = append(formatted, fmt.Sprintf("`%s`", string(v)))
		case protoreflect.FullName:
			formatted = append(formatted, fmt.Sprintf("`%s`", string(v)))
		case string:
			formatted = append(formatted, fmt.Sprintf("`%s`", v))
		case nil:
			formatted = append(formatted, "nil")
		case bool:
			formatted = append(formatted, fmt.Sprintf("%t", v))
		case int, int8, int16, int32, int64:
			formatted = append(formatted, fmt.Sprintf("%d", v))
		case uint, uint8, uint16, uint32, uint64:
			formatted = append(formatted, fmt.Sprintf("%d", v))
		case float32, float64:
			formatted = append(formatted, fmt.Sprintf("%g", v))
		default:
			formatted = append(formatted, fmt.Sprintf("%T(%v)", v, v))
		}
	}
	return strings.Join(formatted, ",")
}

func expandPlaceholders(expression string, args ...any) string {
	result := expression
	argIndex := 0

	// Replace ? placeholders with actual values
	for strings.Contains(result, "?") && argIndex < len(args) {
		arg := args[argIndex]
		var replacement string

		switch v := arg.(type) {
		case proto.Message:
			// Marshal proto message to binary and format as MySQL hex literal
			protoBytes, err := proto.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal proto message: %v", err))
			}
			replacement = fmt.Sprintf("_binary X'%s'", strings.ToUpper(hex.EncodeToString(protoBytes)))
		case []byte:
			replacement = fmt.Sprintf("_binary X'%s'", strings.ToUpper(hex.EncodeToString(v)))
		case string:
			replacement = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
		case nil:
			replacement = "NULL"
		case bool:
			if v {
				replacement = "TRUE"
			} else {
				replacement = "FALSE"
			}
		case int, int8, int16, int32, int64:
			replacement = fmt.Sprintf("%d", v)
		case uint, uint8, uint16, uint32, uint64:
			replacement = fmt.Sprintf("%d", v)
		case float32, float64:
			replacement = fmt.Sprintf("%g", v)
		default:
			replacement = fmt.Sprintf("%v", v)
		}

		// Replace the first occurrence of ?
		pos := strings.Index(result, "?")
		if pos != -1 {
			result = result[:pos] + replacement + result[pos+1:]
		}
		argIndex++
	}

	return result
}

func assertThatExpressionTo[T any](t *testing.T, matcher types.GomegaMatcher, expression string, args ...any) {
	t.Logf("Executing SQL: %s", expandPlaceholders("SELECT "+expression+";", args...))

	g := NewWithT(t)
	g.THelper()

	stmt, err := db.Prepare("SELECT " + expression)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(stmt.Close()).To(Succeed())
	}()

	rows, err := stmt.Query(marshalArgs(args)...)
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
	g.THelper()

	stmt, err := db.Prepare("SELECT " + expression)
	g.Expect(err).NotTo(HaveOccurred())
	defer func() {
		g.Expect(stmt.Close()).To(Succeed())
	}()

	rows, err := stmt.Query(marshalArgs(args)...)
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
	runFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).%s(%s)", expression, formatArguments(args...), method, formatArguments(expected)), func(t *testing.T) {
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

func (this *ExpressionTestContext) IsNegativeZero() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsNegativeZero()", this.Expression, formatArguments(this.Args...)), func(t *testing.T) {
		assertThatExpressionTo[float64](t, gfloat.BeNegativeZero(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsPositiveZero() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsPositiveZero()", this.Expression, formatArguments(this.Args...)), func(t *testing.T) {
		assertThatExpressionTo[float64](t, gfloat.BePositiveZero(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsEqualTo(expected interface{}) {
	if expected == nil {
		panic("Expected value cannot be a nil interface. Use IsNull() or typed nil instead.")
	}
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsEqualTo(%s)", this.Expression, formatArguments(this.Args...), formatArguments(expected)), func(t *testing.T) {
		g := NewWithT(t)
		g.THelper()

		stmt, err := db.Prepare("SELECT " + this.Expression)
		g.Expect(err).NotTo(HaveOccurred())
		defer func() {
			g.Expect(stmt.Close()).To(Succeed())
		}()

		rows, err := stmt.Query(marshalArgs(this.Args)...)
		g.Expect(err).NotTo(HaveOccurred())
		defer func() {
			g.Expect(rows.Close()).To(Succeed())
		}()

		if !rows.Next() { // no rows or an error
			g.Expect(rows.Err()).NotTo(HaveOccurred(), "Unexpected error while iterating over rows.")
			g.Fail("Expected one row, but got none.")
		}

		actual := reflect.New(reflect.ValueOf(expected).Type())
		g.Expect(rows.Scan(actual.Interface())).To(Succeed())

		g.Expect(actual.Elem().Interface()).To(Equal(expected))
	})
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
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsTrue()", this.Expression, formatArguments(this.Args...)), func(t *testing.T) {
		assertThatExpressionTo[bool](this.T, BeTrue(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsFalse() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsFalse()", this.Expression, formatArguments(this.Args...)), func(t *testing.T) {
		assertThatExpressionTo[bool](this.T, BeFalse(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsNull() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsNull()", this.Expression, formatArguments(this.Args...)), func(t *testing.T) {
		assertThatExpressionTo[interface{}](this.T, BeNil(), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsEqualToJsonString(expectedJson string) {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsEqualToJsonString(%s)", this.Expression, formatArguments(this.Args...), formatArguments(expectedJson)), func(t *testing.T) {
		assertThatExpressionTo[string](t, MatchJSON(expectedJson), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsEqualToJson(expectedJson interface{}) {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsEqualToJson(%s)", this.Expression, formatArguments(this.Args...), formatArguments(expectedJson)), func(t *testing.T) {
		jsonBytes := lo.Must(json.Marshal(expectedJson))
		assertThatExpressionTo[string](t, MatchJSON(string(jsonBytes)), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) ToSucceed() {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).ToSucceed()", this.Expression, formatArguments(this.Args...)), func(t *testing.T) {
		assertThatExpressionTo[interface{}](t, SatisfyAny(BeNil(), Not(BeNil())), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) ToFailWithSignalException(state string, containsMessage string) {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).ToFailWithSignalException(%s)", this.Expression, formatArguments(this.Args...), formatArguments(state, containsMessage)), func(t *testing.T) {
		assertThatExpressionToFailWith(t, gmysql.BeMySQLError(1644, state, ContainSubstring(containsMessage)), this.Expression, this.Args...)
	})
}

func (this *ExpressionTestContext) IsEqualToProto(expectedProto proto.Message) {
	this.RunFn(fmt.Sprintf("RunTestThatExpression(`%s`,%s).IsEqualToProto(%s)", this.Expression, formatArguments(this.Args...), formatArguments(expectedProto)), func(t *testing.T) {
		assertThatExpressionTo[[]byte](t, gproto.EqualProto(expectedProto), this.Expression, this.Args...)
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
	this.RunFn(fmt.Sprintf("RunTestThatStatement(`%s`,%s).ShouldReturnSingleRow(%v)", this.Statement, formatArguments(this.Args...), matcher), func(t *testing.T) {
		g := NewWithT(this.T)

		stmt, err := db.Prepare(this.Statement)
		g.Expect(err).NotTo(HaveOccurred())
		defer func() {
			g.Expect(stmt.Close()).To(Succeed())
		}()

		rows, err := stmt.Query(marshalArgs(this.Args)...)
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

func GivenFieldDefinitionsWithExtraFields(t *testing.T, fieldDefinition string, fn func(messageType protoreflect.MessageType)) {
	t.Helper()
	GivenFieldDefinitions(t, fmt.Sprintf("int32 v00 = 100; fixed32 v01 = 101; string v02 = 102; fixed64 v05 = 105; %s; int32 v10 = 110; fixed32 v111 = 111; string v12 = 112; fixed64 v15 = 115;", fieldDefinition), fn)
}

func GivenFieldDefinitions(t *testing.T, fieldDefinition string, fn func(messageType protoreflect.MessageType)) {
	t.Helper()
	support := testutils.NewProtoTestSupport(t, map[string]string{
		"test.proto": fmt.Sprintf(dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|  %s;
			|}
			|message MessageType {
			|    int32 value = 1;
			|}
			|enum EnumType {
			|    ENUM_TYPE_UNSPECIFIED = 0;
			|    ENUM_TYPE_ONE = 1;
			|}
		`), fieldDefinition),
	})
	fn(support.GetMessageType("Test"))
}

func FormatPackedOption(usePacked string) string {
	if usePacked == "" {
		return ""
	}
	return fmt.Sprintf(" [packed = %s]", usePacked)
}
