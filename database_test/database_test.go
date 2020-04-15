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
	"gopkg.in/goracle.v2"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestDatabases(t *testing.T) {
	defer testcontainers.Close()

	dbType, err := database.GetDatabaseType("PostgresDb")
	assert.NilError(t, err)
	assert.Equal(t, dbType, "postgres")

	dbType = database.MustGetDatabaseType("PostgresDb")
	assert.Equal(t, dbType, "postgres")

	pgConn1, err := database.GetConnection("PostgresDb")
	assert.NilError(t, err)
	defer utils.CloseSafe(pgConn1)

	pgConn2 := database.MustGetConnection("PostgresDb")
	defer utils.CloseSafe(pgConn2)
	err = pgConn2.PingContext(context.Background())
	assert.NilError(t, err)
	prepareDb(pgConn2, t)

	values, err := database.GetValues("PostgresDb", "select * from public.test")
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

	values = database.MustGetValues("PostgresDb", "select * from public.test where $1=1", 1)
	assert.Equal(t, values[0][0].(int64), int64(1))
	assert.Equal(t, values[0][1].(string), "text")

	rows := make([]TestPostgres, 0)
	database.MustScanStructSlice(&rows, "PostgresDb", `SELECT * FROM public.test`)
	assertTestPostgres(t, rows[0])
	var text *string
	var date *time.Time
	var geom *postgis.Point
	assert.Equal(t, rows[1].Id, int64(2))
	assert.Equal(t, rows[1].Text, text)
	assert.Equal(t, rows[1].Date, date)
	assert.Equal(t, rows[1].Geom, geom)

	var test TestPostgres
	database.MustScanStruct(&test, "PostgresDb", `SELECT * FROM public.test WHERE id = 1`)
	assertTestPostgres(t, test)

	dbType = database.MustGetDatabaseType("OracleDb")
	assert.Equal(t, dbType, "oracle")

	oraConn := database.MustGetConnection("OracleDb")
	defer utils.CloseSafe(oraConn)
	err = oraConn.PingContext(context.Background())
	assert.NilError(t, err)

	prepareOracleDb(oraConn, t)

	values, err = database.GetValuesConn(oraConn, `SELECT * FROM test WHERE :1=1`, 1)
	assert.NilError(t, err)
	assert.Equal(t, values[0][0].(goracle.Number), goracle.Number("1"))
	assert.Equal(t, values[0][1].(string), "text")

	var testOracle TestOracle
	err = database.ScanStruct(&testOracle, "OracleDb", `SELECT * FROM test`)
	assert.NilError(t, err)
	assert.Equal(t, testOracle.Id, int64(1))
	assert.Equal(t, testOracle.Text, "text")

	testOracleSlice := make([]TestOracle, 0)
	err = database.ScanStructSlice(&testOracleSlice, "OracleDb", `SELECT * FROM test WHERE :1=1`, 1)
	assert.NilError(t, err)
	assert.Equal(t, testOracleSlice[0].Id, int64(1))
	assert.Equal(t, testOracleSlice[0].Text, "text")

	dbType = database.MustGetDatabaseType("PrestoDb")
	assert.Equal(t, dbType, "presto")

	values = database.MustGetValuesConn(database.MustGetConnection("PrestoDb"), "select * from postgresql.public.test")
	assert.NilError(t, err)
	assert.Equal(t, values[0][0].(int64), int64(1))
	assert.Equal(t, values[0][1].(string), "text")
	tm, err = time.Parse("2006-01-02", "2016-06-22")
	assert.NilError(t, err)
	assert.Equal(t, values[0][2].(time.Time).Nanosecond(), tm.Nanosecond())

	var test1 TestPostgres
	database.MustScanStruct(&test1, "PrestoDb", `select * from postgresql.public.test where id = ?`, 1)
	assert.Equal(t, test1.Id, int64(1))
	assert.Equal(t, *test1.Text, "text")
	tm, err = time.Parse("2006-01-02", "2016-06-22")
	assert.NilError(t, err)
	assert.Equal(t, test1.Date.Nanosecond(), tm.Nanosecond())
}

func assertTestPostgres(t *testing.T, obj TestPostgres) {
	assert.Equal(t, obj.Id, int64(1))
	assert.Equal(t, *obj.Text, "text")
	tm, err := time.Parse("2006-01-02", "2016-06-22")
	assert.NilError(t, err)
	assert.Equal(t, obj.Date.Nanosecond(), tm.Nanosecond())
	assert.Equal(t, obj.Geom.X, 10.0)
	assert.Equal(t, obj.Geom.Y, 20.0)
}

type TestOracle struct {
	Id   int64  `db:"ID"`
	Text string `db:"TEXT"`
}

type TestPostgres struct {
	Id   int64          `db:"id"`
	Text *string        `db:"text"`
	Date *time.Time     `db:"date"`
	Geom *postgis.Point `db:"geom"`
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

	query = `insert into public.test(id, text, date, geom) values(2, null, null, null)`
	_, err = conn.ExecContext(context.Background(), query)
	assert.NilError(t, err)
}

func prepareOracleDb(conn *sql.Conn, t *testing.T) {
	query := `
	create table test(
		id NUMBER,
		text VARCHAR2(255)
	)
	`
	_, err := conn.ExecContext(context.Background(), query)
	assert.NilError(t, err)

	query = `insert into test(id, text) values(1, 'text')`
	_, err = conn.ExecContext(context.Background(), query)
	assert.NilError(t, err)
}
