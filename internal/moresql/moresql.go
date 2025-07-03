package moresql

import (
	"database/sql"
	"reflect"

	"github.com/eiiches/mysql-protobuf-functions/internal/caseconv"
)

func ScanStruct[T any](rows *sql.Rows) (*T, error) {
	structPointer := new(T)
	structPointerValue := reflect.ValueOf(structPointer)
	structValue := structPointerValue.Elem()
	structType := structValue.Type()
	if structType.Kind() != reflect.Struct {
		panic("ScanStruct only works with struct types")
	}

	fieldAddrValues := []any{}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for _, column := range columns {
		fieldName := caseconv.SnakeToUpperCamel(column)

		fieldValue := structValue.FieldByName(fieldName)
		if !fieldValue.IsValid() {
			panic("ScanStruct: field " + fieldName + " not found in struct " + structType.Name())
		}

		fieldAddrValues = append(fieldAddrValues, fieldValue.Addr().Interface())
	}

	if err := rows.Scan(fieldAddrValues...); err != nil {
		return nil, err
	}

	return structPointer, nil
}
