package database

import (
	"context"
	"database/sql"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
)

func GetValues(dbName string, query string, args ...interface{}) ([][]interface{}, error) {
	conn, err := GetConnection(dbName)
	if err != nil {
		return nil, errors.Wrapf(err, "database.GetValues -> GetConnection(%s)", dbName)
	}
	defer utils.CloseSafe(conn)
	return GetValuesConn(conn, query, args...)
}

func MustGetValues(dbName string, query string, args ...interface{}) [][]interface{} {
	if values, err := GetValues(dbName, query, args...); err != nil {
		panic(err)
	} else {
		return values
	}
}

func GetValuesConn(conn *sql.Conn, query string, args ...interface{}) ([][]interface{}, error) {
	rows, err := conn.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "database.GetValuesConn -> conn.QueryContext(%s)", query)
	}
	defer utils.CloseSafe(rows)

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrapf(err, "database.GetValuesConn -> rows.Columns")
	}

	ptrs := make([]interface{}, len(columns))
	result := make([][]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := 0; i < len(values); i++ {
			ptrs[i] = &values[i]
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			return nil, errors.Wrapf(err, "database.GetValuesConn -> rows.Scan")
		}

		result = append(result, values)
	}

	return result, nil
}

func MustGetValuesConn(conn *sql.Conn, query string, args ...interface{}) [][]interface{} {
	if values, err := GetValuesConn(conn, query, args...); err != nil {
		panic(err)
	} else {
		return values
	}
}

func ScanStruct(object interface{}, dbName string, query string, args ...interface{}) error {
	db, err := getDB(dbName)
	if err != nil {
		return errors.Wrapf(err, "database.ScanStruct -> getDB(%s)", dbName)
	}
	err = db.Get(object, query, args...)
	if err != nil {
		return errors.Wrapf(err, "database.ScanStruct -> db.Get(%s)", query)
	}
	return nil
}

func MustScanStruct(object interface{}, dbName string, query string, args ...interface{}) {
	if err := ScanStruct(object, dbName, query, args...); err != nil {
		panic(err)
	}
}

func ScanStructSlice(slice interface{}, dbName string, query string, args ...interface{}) error {
	db, err := getDB(dbName)
	if err != nil {
		return errors.Wrapf(err, "database.ScanStructSlice -> getDB(%s)", dbName)
	}
	err = db.Select(slice, query, args...)
	if err != nil {
		return errors.Wrapf(err, "database.ScanStructSlice -> db.Select(%s)", query)
	}
	return nil
}

func MustScanStructSlice(slice interface{}, dbName string, query string, args ...interface{}) {
	if err := ScanStructSlice(slice, dbName, query, args...); err != nil {
		panic(err)
	}
}
