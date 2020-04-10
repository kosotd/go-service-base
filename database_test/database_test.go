package database_test

// don't change order
import _ "github.com/kosotd/go-service-base/testcontainers"
import _ "github.com/kosotd/go-service-base/config"
import _ "github.com/kosotd/go-service-base/database"
import "github.com/kosotd/go-service-base/testcontainers"
import "github.com/kosotd/go-service-base/database"

// don't change order

import (
	"context"
	"database/sql"
	"github.com/cridenour/go-postgis"
	"github.com/kosotd/go-service-base/utils"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestDatabases(t *testing.T) {
	defer testcontainers.Close()

	dbType := database.MustGetDatabaseType("PostgresDb")
	assert.Equal(t, dbType, "postgres")

	conn := database.MustGetConnection("PostgresDb")
	defer utils.CloseSafe(conn)
	err := conn.PingContext(context.Background())
	assert.NilError(t, err)

	dbType, err = database.GetDatabaseType("PostgresDb")
	assert.NilError(t, err)
	assert.Equal(t, dbType, "postgres")

	conn, err = database.GetConnection("PostgresDb")
	assert.NilError(t, err)
	defer utils.CloseSafe(conn)
	err = conn.PingContext(context.Background())
	assert.NilError(t, err)

	prepareDb(conn, t)

	values := database.MustGetValues(conn, "select * from public.test")
	assert.Equal(t, values[0][0].(int64), int64(1))
	assert.Equal(t, values[0][1].(string), "text")

	values, err = database.GetValues(conn, "select * from public.test")
	assert.NilError(t, err)
	assert.Equal(t, values[0][0].(int64), int64(1))
	assert.Equal(t, values[0][1].(string), "text")
	tm, err := time.Parse("2006-01-02", "2016-06-22")
	assert.NilError(t, err)
	assert.Equal(t, values[0][2].(time.Time).Nanosecond(), tm.Nanosecond())
	var point postgis.Point
	err = point.Scan(values[0][3])
	assert.NilError(t, err)
	assert.Equal(t, point.X, 10.0)
	assert.Equal(t, point.Y, 20.0)

	dbType, err = database.GetDatabaseType("OracleDb")
	assert.NilError(t, err)
	assert.Equal(t, dbType, "oracle")

	conn, err = database.GetConnection("OracleDb")
	assert.NilError(t, err)
	defer utils.CloseSafe(conn)
	err = conn.PingContext(context.Background())
	assert.NilError(t, err)

	dbType, err = database.GetDatabaseType("PrestoDb")
	assert.NilError(t, err)
	assert.Equal(t, dbType, "presto")

	conn, err = database.GetConnection("PrestoDb")
	assert.NilError(t, err)
	defer utils.CloseSafe(conn)
	err = conn.PingContext(context.Background())
	assert.NilError(t, err)
}

func prepareDb(conn *sql.Conn, t *testing.T) {
	query := `
	create table public.test(
		id int,
		text text,
		date timestamp,
		geom geometry
	)
	`
	_, err := conn.ExecContext(context.Background(), query)
	assert.NilError(t, err)

	query = `insert into public.test(id, text, date, geom) values(1, 'text', '2016-06-22', ST_GeomFromText('POINT(10 20)'))`
	_, err = conn.ExecContext(context.Background(), query)
	assert.NilError(t, err)
}
