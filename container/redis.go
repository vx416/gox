package container

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest"
)

type Redis struct {
	Addr string
}

func (builder *Builder) RunRedis(name string, port ...string) (*Redis, error) {
	container, err := builder.FindContainer(name)
	if err != nil {
		return nil, err
	}

	if container != nil {
		builder.containerIDs[container.ID] = true
		r := &Redis{
			Addr: fmt.Sprintf("localhost:%d", container.Ports[0].PublicPort),
		}
		return r, nil
	}

	resource, err := builder.RunWithOptions(&dockertest.RunOptions{Repository: "redis", Tag: "6.0.9-alpine", Name: name})
	if err != nil {
		return nil, err
	}

	builder.containerIDs[resource.Container.ID] = true
	r := &Redis{
		Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
	}

	err = builder.Retry(func() error {
		ctx := context.Background()
		client := redis.NewClient(&redis.Options{
			Addr: r.Addr,
		})

		return client.Ping(ctx).Err()
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}
