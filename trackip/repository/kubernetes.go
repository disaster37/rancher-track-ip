package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/disaster37/rancher-track-ip/model"
	"github.com/disaster37/rancher-track-ip/trackip"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubernetesRepository struct {
	Conn kubernetes.Interface
}

func NewKubernetesRepository(conn kubernetes.Interface) trackip.RancherRepository {
	return &kubernetesRepository{
		Conn: conn,
	}
}

func (h *kubernetesRepository) GetContainers(ctx context.Context) (listContainers []*model.Container, err error) {

	listContainers = make([]*model.Container, 0)

	pods, err := h.Conn.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	log.Debugf("Found %d containers", len(pods.Items))

	for _, pod := range pods.Items {
		containerInfo := &model.Container{
			ID:        string(pod.UID),
			IP:        pod.Status.PodIP,
			HostIP:    pod.Status.HostIP,
			Hostname:  pod.Spec.NodeName,
			Status:    string(pod.Status.Phase),
			Name:      pod.Name,
			Project:   fmt.Sprintf("%s", pod.Namespace),
			StartedAt: pod.CreationTimestamp.Time,
			Platform:  "kubernetes",
		}

		if pod.DeletionTimestamp != nil {
			containerInfo.FinishedAt = pod.DeletionTimestamp.Time
		}

		// Get all pods images
		var image strings.Builder
		for _, container := range pod.Spec.InitContainers {
			image.WriteString(container.Image + "\n")
		}
		for _, container := range pod.Spec.Containers {
			image.WriteString(container.Image + "\n")
		}
		containerInfo.Image = image.String()

		listContainers = append(listContainers, containerInfo)
	}

	return

}
