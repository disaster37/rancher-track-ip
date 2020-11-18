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
				ID:      fmt.Sprintf("%s_%s_%s", container.Id, container.Uuid, container.ExternalId),
				IP:      container.PrimaryIpAddress,
				Status:  container.State,
				Name:    container.Name,
				Project: fmt.Sprintf("%s/%s", container.AccountId, container.Labels["io.rancher.project_service.name"]),
				Image:   container.ImageUuid,
			}

			containerInfo.StartedAt, err = time.Parse(time.RFC3339, container.Created)
			if err != nil {
				return nil, err
			}

			// Get host info
			hosts := &rancherClient.HostCollection{}
			err = h.Conn.GetLink(container.Resource, "hosts", hosts)
			if err != nil {
				return nil, err
			}
			if len(hosts.Data) > 0 {
				containerInfo.Hostname = hosts.Data[0].Name
				containerInfo.HostIP = hosts.Data[0].AgentIpAddress
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
