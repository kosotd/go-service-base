package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kosotd/go-service-base/config"
	"github.com/kosotd/go-service-base/database/domain"
	"github.com/kosotd/go-service-base/utils"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	_ "github.com/prestodb/presto-go-client/presto"
	_ "gopkg.in/goracle.v2"
	"strings"
	"sync"
	"time"
)

var connections sync.Map
var dbTypes sync.Map

func init() {
	databases := config.Databases()
	for _, database := range databases {
		conn, err := parseUrl(database)
		if err != nil {
			utils.LogError(errors.Wrapf(err, "database.init -> parseUrl(%s)", database).Error())
			continue
		}
		db, err := sql.Open(conn.Driver, conn.Url)
		if err != nil {
			utils.LogError(errors.Wrapf(err, "database.init -> sql.Open(%s)", conn.Url).Error())
			continue
		}
		err = pingTimeout(db, conn)
		if err != nil {
			_ = db.Close()
			utils.LogError(errors.Wrapf(err, "database.init -> pingTimeout").Error())
			continue
		}
		utils.LogInfo("%s connected successfully", conn.Name)
		name := strings.ToLower(strings.Trim(conn.Name, " "))
		connections.Store(name, db)
		dbTypes.Store(name, strings.ToLower(strings.Trim(conn.DbType, " ")))
	}
}

func pingTimeout(db *sql.DB, conn domain.Connection) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- db.Ping()
	}()

	select {
	case <-time.After(5 * time.Second):
		return errors.Errorf("ping %s database timeout", conn.Name)
	case err := <-errChan:
		return err
	}
}

func GetConnection(name string) (*sql.Conn, error) {
	if value, ok := connections.Load(strings.ToLower(strings.Trim(name, " "))); ok {
		db := value.(*sql.DB)
		conn, err := db.Conn(context.Background())
		if err != nil {
			return nil, errors.Wrapf(err, "database.GetConnection -> db.Conn")
		}
		return conn, nil
	}
	return nil, errors.New(fmt.Sprintf("connection with name %s not found", name))
}

func GetDatabaseType(name string) (string, error) {
	if value, ok := dbTypes.Load(strings.ToLower(strings.Trim(name, " "))); ok {
		return value.(string), nil
	}
	return "", errors.New(fmt.Sprintf("connection with name %s not found", name))
}

func Close() {
	connections.Range(func(key, value interface{}) bool {
		db := value.(*sql.DB)
		if db != nil {
			_ = db.Close()
		}
		return true
	})
}
