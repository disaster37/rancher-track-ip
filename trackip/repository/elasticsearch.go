package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/disaster37/rancher-track-ip/model"
	"github.com/disaster37/rancher-track-ip/trackip"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type elasticsearchRepository struct {
	Conn      *elastic.Client
	IndexName string
}

// NewElasticsearchRepository will create an object that implement RepositoryElasticsearch interface
func NewElasticsearchRepository(conn *elastic.Client, index string) trackip.ElasticsearchRepository {
	return &elasticsearchRepository{
		Conn:      conn,
		IndexName: index,
	}
}

// Create or update document on Elasticsearch
func (h *elasticsearchRepository) Index(ctx context.Context, data *model.Container, id ...string) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	log.Debugf("Data: %+v", data)

	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(dataJson)

	finalID := ""

	if len(id) > 0 {
		finalID = id[0]
	}

	res, err := h.Conn.Index(
		h.IndexName,
		b,
		h.Conn.Index.WithDocumentID(finalID),
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	log.Debugf("Response: %s", res.String())

	return nil
}
