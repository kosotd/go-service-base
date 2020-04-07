package database

import (
	"gotest.tools/assert"
	"testing"
)

func TestParseUrl(t *testing.T) {
	connection, err := parseUrl("OracleDb;Oracle:user/password@127.0.0.1:1521/ORCL")
	assert.NilError(t, err)
	assert.Equal(t, connection.Driver, "goracle")
	assert.Equal(t, connection.DbType, "Oracle")
	assert.Equal(t, connection.Name, "OracleDb")
	assert.Equal(t, connection.Url, "user/password@127.0.0.1:1521/ORCL")

	connection, err = parseUrl("PostgresDb;poSTgres:user/password@127.0.0.1:5432/database")
	assert.NilError(t, err)
	assert.Equal(t, connection.Driver, "postgres")
	assert.Equal(t, connection.DbType, "poSTgres")
	assert.Equal(t, connection.Name, "PostgresDb")
	assert.Equal(t, connection.Url, "host=127.0.0.1 port=5432 user=user password=password dbname=database sslmode=disable")

	connection, err = parseUrl("PrestoDb;presto:user/password@127.0.0.1:8080")
	assert.NilError(t, err)
	assert.Equal(t, connection.Driver, "presto")
	assert.Equal(t, connection.DbType, "presto")
	assert.Equal(t, connection.Name, "PrestoDb")
	assert.Equal(t, connection.Url, "http://user@127.0.0.1:8080")

	_, err = parseUrl("PrestoDb;presto:user/password@localhost:8080")
	assert.NilError(t, err)

	_, err = parseUrl("PrestoDb;presto:user/password@127.0.0.1")
	assert.Error(t, err, "url PrestoDb;presto:user/password@127.0.0.1 doesn't match with pattern")

	_, err = parseUrl("OracleDb;oracle1:user/password@127.0.0.1:1521/ORCL")
	assert.Error(t, err, "database.parseUrl -> getDriverName(oracle1): unknown database type oracle1")
}