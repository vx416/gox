package container

import (
	"strconv"

	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
)

type PG struct {
	Addr     string
	Port     int32
	Username string
	Password string
	DBName   string
}

func (builder *Builder) RunPg(name string, dbName string, port ...string) (*PG, error) {
	container, err := builder.FindContainer(name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		builder.containerIDs[container.ID] = true
		return &PG{
			Port:     int32(container.Ports[0].PublicPort),
			DBName:   dbName,
			Addr:     "localhost",
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
			dc.PortBinding{
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
	return &PG{
		Port:     int32(dbPort),
		DBName:   dbName,
		Addr:     "localhost",
		Username: "test",
		Password: "test",
	}, nil
}
