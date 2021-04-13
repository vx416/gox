package container

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
)

type DB struct {
	Host     string
	Port     int32
	Username string
	Password string
	DBName   string
}

func (builder *Builder) RunPg(name string, dbName string, port ...string) (*DB, error) {
	container, err := builder.FindContainer(name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		builder.containerIDs[container.ID] = true
		return &DB{
			Port:     int32(container.Ports[0].PublicPort),
			DBName:   dbName,
			Host:     "localhost",
			Username: "test",
			Password: "test",
		}, nil
	}

	options := &dockertest.RunOptions{
		Repository: "postgres", Tag: "12.3-alpine", Name: name,
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=" + dbName,
		},
	}
	if len(port) == 1 {
		options.PortBindings = make(map[dc.Port][]dc.PortBinding)
		options.PortBindings[dc.Port("5432/tcp")] = []dc.PortBinding{
			{
				HostPort: port[0],
			}}
	}

	resource, err := builder.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	builder.containerIDs[resource.Container.ID] = true

	dbPort, err := strconv.ParseInt(resource.GetPort("5432/tcp"), 10, 64)
	if err != nil {
		return nil, err
	}
	pg := &DB{
		Port:     int32(dbPort),
		DBName:   dbName,
		Host:     "localhost",
		Username: "test",
		Password: "test",
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pg.Username, pg.Password, pg.Host, pg.Port, pg.DBName)
	err = builder.Retry(func() error {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			fmt.Println("testing, err", err)
			return err
		}
		defer db.Close()

		return db.Ping()
	})
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func (builder *Builder) RunMysql(name string, dbName string, port ...string) (*DB, error) {
	container, err := builder.FindContainer(name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		builder.containerIDs[container.ID] = true
		return &DB{
			Port:     int32(container.Ports[0].PublicPort),
			DBName:   dbName,
			Host:     "localhost",
			Username: "test",
			Password: "test",
		}, nil
	}

	options := &dockertest.RunOptions{
		Repository: "mysql", Tag: "8", Name: name,
		Env: []string{
			"MYSQL_USER=test",
			"MYSQL_PASSWORD=test",
			"MYSQL_ROOT_PASSWORD=test",
			"MYSQL_DATABASE=" + dbName,
		},
	}
	if len(port) == 1 {
		options.PortBindings = make(map[dc.Port][]dc.PortBinding)
		options.PortBindings[dc.Port("3306/tcp")] = []dc.PortBinding{
			{
				HostPort: port[0],
			}}
	}

	resource, err := builder.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	builder.containerIDs[resource.Container.ID] = true

	dbPort, err := strconv.ParseInt(resource.GetPort("3306/tcp"), 10, 64)
	if err != nil {
		return nil, err
	}
	db := &DB{
		Port:     int32(dbPort),
		DBName:   dbName,
		Host:     "localhost",
		Username: "test",
		Password: "test",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		db.Username, db.Password, db.Host, db.Port, db.DBName)

	err = builder.Retry(func() error {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
