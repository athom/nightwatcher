package nightwatcher

import (
	"context"

	"github.com/olivere/elastic"
)

func NewElasticsearchReporter(esURLs []string, index string) (r *ElasticsearchReporter, err error) {
	var c *elastic.Client
	c, err = elastic.NewClient(
		elastic.SetURL(esURLs...),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		return
	}

	r = &ElasticsearchReporter{
		client: c,
		index:  index,
	}
	return
}

type ElasticsearchReporter struct {
	client *elastic.Client
	index  string
}

func (this ElasticsearchReporter) Output(aim *Aim) (err error) {
	ctx := context.Background()

	var exists bool
	exists, err = this.client.IndexExists(this.index).Do(ctx)
	if err != nil {
		return
	}
	if !exists {
		_, err = this.client.CreateIndex(this.index).Do(ctx)
		if err != nil {
			return
		}
	}

	_, err = this.client.Index().Index(this.index).Type("nightwatcher").BodyJson(aim.ToJson()).Do(ctx)
	if err != nil {
		return
	}
	return
}
