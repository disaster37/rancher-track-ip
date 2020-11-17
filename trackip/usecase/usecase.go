package usecase

import (
	"context"
	"time"

	"github.com/disaster37/rancher-track-ip/trackip"
)

type trackIPUsecase struct {
	ElasticsearchRepo trackip.ElasticsearchRepository
	RancherRepo       trackip.RancherRepository
}

func NewTrackIPUsecase(elasticsearchRepo trackip.ElasticsearchRepository, rancherRepo trackip.RancherRepository) trackip.Usecase {
	return &trackIPUsecase{
		ElasticsearchRepo: elasticsearchRepo,
		RancherRepo:       rancherRepo,
	}
}

func (h *trackIPUsecase) TrackContainers(ctx context.Context, loopIntervalSecond int64) error {

	for {
		time.Sleep(time.Duration(loopIntervalSecond) * time.Second)

		containers, err := h.RancherRepo.GetContainers(ctx)
		if err != nil {
			return err
		}

		for _, container := range containers {
			err = h.ElasticsearchRepo.Index(ctx, container, container.ID)
			if err != nil {
				return err
			}
		}
	}
}
