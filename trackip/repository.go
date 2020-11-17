package trackip

import (
	"context"

	"github.com/disaster37/rancher-track-ip/model"
)

type ElasticsearchRepository interface {
	Index(ctx context.Context, container *model.Container, id ...string) error
}

type RancherRepository interface {
	GetContainers(ctx context.Context) (listContainers []*model.Container, err error)
}
