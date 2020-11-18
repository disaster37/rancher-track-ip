package model

import (
	"encoding/json"
	"time"
)

type Container struct {
	ID         string    `json:"container.id"`
	IP         string    `json:"container.ip,omitempty"`
	Hostname   string    `json:"host.name"`
	HostIP     string    `json:"host.ip"`
	StartedAt  time.Time `json:"container.start"`
	FinishedAt time.Time `json:"container.stop"`
	Status     string    `json:"container.status"`
	Name       string    `json:"container.name"`
	Project    string    `json:"container.labels.project,omitempty"`
	Image      string    `json:"container.image.name"`
	Platform   string    `json:"container.platform"`
}

func (h *Container) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}
