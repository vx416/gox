package container

import (
	"strings"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

type Builder struct {
	*dockertest.Pool
	containerIDs map[string]bool
}

func (builder *Builder) RemoveByID(containerID string) error {
	return builder.Client.RemoveContainer(docker.RemoveContainerOptions{ID: containerID, Force: true, RemoveVolumes: true})
}

func (builder *Builder) FindContainer(containerName string) (*docker.APIContainers, error) {
	containers, err := builder.Client.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": {containerName},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return nil, nil
	}
	if strings.Contains(containers[0].Status, "Exited") {
		if err := builder.RemoveByID(containers[0].ID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	return &containers[0], nil
}

func (builder *Builder) PruneAll() error {
	var err error
	for id := range builder.containerIDs {
		purErr := builder.RemoveByID(id)
		if purErr != nil {
			err = purErr
		}
	}
	return err
}

func NewConBuilder() (*Builder, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	return &Builder{
		Pool:         pool,
		containerIDs: make(map[string]bool),
	}, nil
}
