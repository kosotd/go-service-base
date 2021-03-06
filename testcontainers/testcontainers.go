package testcontainers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/docker/go-connections/nat"
	"github.com/kosotd/go-service-base/utils"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	_ "github.com/prestodb/presto-go-client/presto"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "gopkg.in/goracle.v2"
	"io/ioutil"
	"os"
	"time"
)

var postgres testcontainers.Container
var oracle testcontainers.Container
var presto testcontainers.Container
var network testcontainers.Network
var networkName string

func init() {
	provider, err := testcontainers.NewDockerProvider()
	utils.FailIfError(err)
	networkName = uuid.Generate().String()
	network, err = provider.CreateNetwork(context.Background(), testcontainers.NetworkRequest{Name: networkName})
	utils.FailIfError(err)

	postgresUrl := initPostgres()
	oracleUrl := initOracle()
	prestoUrl := initPresto()
	_ = ioutil.WriteFile("config.env", []byte(fmt.Sprintf(`DATABASES=["%s", "%s", "%s"]`, postgresUrl, oracleUrl, prestoUrl)), os.ModePerm)
}

func Close() {
	if postgres != nil {
		_ = postgres.Terminate(context.Background())
	}

	if oracle != nil {
		_ = oracle.Terminate(context.Background())
	}

	if presto != nil {
		_ = presto.Terminate(context.Background())
	}

	if network != nil {
		_ = network.Remove(context.Background())
	}
}

func initPostgres() string {
	ctx := context.Background()

	strPort := "5432/tcp"
	postgresEnv := make(map[string]string)
	postgresEnv["POSTGRES_PASSWORD"] = "postgres"
	postgresEnv["POSTGRES_USER"] = "postgres"
	postgresEnv["POSTGRES_DB"] = "postgres"
	postgres = initContainer("test_postgres", "mdillon/postgis:9.6", strPort, postgresEnv)

	host, err := postgres.Host(ctx)
	utils.FailIfError(err)
	mappedPort, err := postgres.MappedPort(ctx, nat.Port(strPort))
	utils.FailIfError(err)

	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, mappedPort.Port(), "postgres", "postgres", "postgres")
	db, err := sql.Open("postgres", url)
	utils.FailIfError(err)
	defer utils.CloseSafe(db)

	err = waitForPing(db)
	utils.FailIfError(err)

	return fmt.Sprintf("PostgresDb;poSTgres:postgres/postgres@%s:%s/postgres", host, mappedPort.Port())
}

func initOracle() string {
	ctx := context.Background()

	strPort := "1521/tcp"
	oracle = initContainer("test_oralce", "oracleinanutshell/oracle-xe-11g", "1521/tcp", make(map[string]string))

	host, err := oracle.Host(ctx)
	utils.FailIfError(err)
	mappedPort, err := oracle.MappedPort(ctx, nat.Port(strPort))
	utils.FailIfError(err)

	url := fmt.Sprintf("%s/%s@%s:%s/%s", "system", "oracle", host, mappedPort.Port(), "xe")
	db, err := sql.Open("goracle", url)
	utils.FailIfError(err)
	defer utils.CloseSafe(db)

	err = waitForPing(db)
	utils.FailIfError(err)

	return fmt.Sprintf("OracleDb;Oracle:system/oracle@%s:%s/xe", host, mappedPort.Port())
}

func initPresto() string {
	ctx := context.Background()

	connector := fmt.Sprintf(`connector.name=postgresql
connection-url=jdbc:postgresql://%s:%s/postgres
connection-user=postgres
connection-password=postgres`, "test_postgres", "5432")

	err := ioutil.WriteFile("./presto-image/catalog/postgresql.properties", []byte(connector), os.ModePerm)
	utils.FailIfError(err)

	strPort := "8080/tcp"
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Dockerfile: "Dockerfile",
			Context:    "./presto-image",
		},
		ExposedPorts: []string{strPort},
		Networks:     []string{networkName},
		WaitingFor:   wait.ForLog("======== SERVER STARTED ========"),
	}

	presto, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	utils.FailIfError(err)

	host, err := presto.Host(ctx)
	utils.FailIfError(err)
	mappedPort, err := presto.MappedPort(ctx, nat.Port(strPort))
	utils.FailIfError(err)

	url := fmt.Sprintf("http://user@%s:%s", host, mappedPort.Port())
	db, err := sql.Open("presto", url)
	utils.FailIfError(err)
	defer utils.CloseSafe(db)

	err = waitForPing(db)
	utils.FailIfError(err)

	time.Sleep(3 * time.Second)

	return fmt.Sprintf("PrestoDb;presto:user/password@%s:%s", host, mappedPort.Port())
}

func initContainer(name string, image string, port string, env map[string]string) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Name:         name,
		Image:        image,
		Env:          env,
		ExposedPorts: []string{port},
		Networks:     []string{networkName},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	utils.FailIfError(err)
	return container
}

func waitForPing(db *sql.DB) error {
	ping := make(chan bool, 1)
	go func() {
		for {
			if err := db.Ping(); err == nil {
				ping <- true
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {
	case <-ping:
		return nil
	case <-time.After(60 * time.Second):
		return errors.New("wait for ping timeout")
	}
}
