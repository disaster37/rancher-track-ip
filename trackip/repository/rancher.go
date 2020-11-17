package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/rancher-track-ip/model"
	"github.com/disaster37/rancher-track-ip/trackip"
	rancherClient "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

type rancherRepository struct {
	Conn *rancherClient.RancherClient
}

func NewRancherRepository(conn *rancherClient.RancherClient) trackip.RancherRepository {
	return &rancherRepository{
		Conn: conn,
	}
}

func (h *rancherRepository) GetContainers(ctx context.Context) (listContainers []*model.Container, err error) {

	listContainers = make([]*model.Container, 0)

	containers, err := h.Conn.Container.List(nil)
	if err != nil {
		return
	}

	for containers != nil {
		log.Debugf("Found %d containers", len(containers.Data))

		for _, container := range containers.Data {
			containerInfo := &model.Container{
				ID:       fmt.Sprintf("%s_%s", container.Uuid, container.ExternalId),
				IP:       container.Ip,
				Hostname: container.Hostname,
				Status:   container.State,
				Name:     container.Name,
				Project:  fmt.Sprintf("%s/%s", container.StackId, container.ServiceId),
				Image:    container.ImageUuid,
			}

			containerInfo.StartedAt, err = time.Parse(time.RFC3339, container.Created)
			if err != nil {
				return nil, err
			}

			listContainers = append(listContainers, containerInfo)
		}

		containers, err = containers.Next()
		if err != nil {
			return nil, err
		}
	}

	return

}
