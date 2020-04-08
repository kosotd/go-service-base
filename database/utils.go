package database

import (
	"context"
	"database/sql"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
)

func GetValues(db *sql.Conn, query string) ([][]interface{}, error) {
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, errors.Wrapf(err, "database.GetValues -> db.QueryContext(%s)", query)
	}
	defer utils.CloseSafe(rows)

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrapf(err, "database.GetValues -> rows.Columns")
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
			return nil, errors.Wrapf(err, "database.GetValues -> rows.Scan")
		}

		result = append(result, values)
	}

	return result, nil
}
