package database_test

// don't change order
import _ "github.com/kosotd/go-service-base/testcontainers"
import _ "github.com/kosotd/go-service-base/config"
import _ "github.com/kosotd/go-service-base/database"
import "github.com/kosotd/go-service-base/database"
import "github.com/kosotd/go-service-base/testcontainers"
// don't change order

import (
	"context"
	"github.com/kosotd/go-service-base/utils"
	"gotest.tools/assert"
	"testing"
)

func TestDatabases(t *testing.T) {
	defer testcontainers.Close()

	dbType, err := database.GetDatabaseType("PostgresDb")
	assert.NilError(t, err)
	assert.Equal(t, dbType, "postgres")

	conn, err := database.GetConnection("PostgresDb")
	assert.NilError(t, err)
	defer utils.CloseSafe(conn)
	err = conn.PingContext(context.Background())
	assert.NilError(t, err)

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
