package model

import (
	"encoding/json"
	"time"
)

type Container struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip,omitempty"`
	Hostname  string    `json:"hostname"`
	StartedAt time.Time `json:"started_at"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	Project   string    `json:"project,omitempty"`
	Image     string    `json:"image"`
}

func (h *Container) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}
