package container

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
)

type Redis struct {
	Host string
	Port int64
}

func (builder *Builder) RunRedis(name string, port ...string) (*Redis, error) {
	container, err := builder.FindContainer(name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		builder.containerIDs[container.ID] = true
		r := &Redis{
			Host: "localhost",
			Port: container.Ports[0].PublicPort,
		}
		return r, nil
	}
	options := &dockertest.RunOptions{Repository: "redis", Tag: "6.0.9-alpine", Name: name}
	if len(port) == 1 {
		options.PortBindings = make(map[dc.Port][]dc.PortBinding)
		options.PortBindings[dc.Port("6379/tcp")] = []dc.PortBinding{
			{
				HostPort: port[0],
			}}
	}

	resource, err := builder.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	builder.containerIDs[resource.Container.ID] = true
	redisPort, err := strconv.ParseInt(resource.GetPort("6379/tcp"), 10, 64)
	if err != nil {
		return nil, err
	}

	r := &Redis{
		Host: "localhost",
		Port: redisPort,
	}

	err = builder.Retry(func() error {
		ctx := context.Background()
		client := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", r.Host, r.Port),
		})

		return client.Ping(ctx).Err()
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
