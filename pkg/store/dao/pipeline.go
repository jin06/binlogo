package dao

import (
	"context"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	storeredis "github.com/jin06/binlogo/v2/pkg/store/redis"
)

// PipelinePrefix returns etcd prefix of pipeline info
func PipelinePrefix() string {
	return storeredis.Prefix() + "/pipeline"
}

func PipelinesKey() string {
	return storeredis.Prefix() + "/pipelines"
}

// GetPipeline get pipeline info from etcd
func GetPipeline(ctx context.Context, name string) (p *pipeline.Pipeline, err error) {
	return myDao.GetPipeline(ctx, name)
}

// CreatePipeline write pipeline info to etcd
func CreatePipeline(ctx context.Context, d *pipeline.Pipeline) (ok bool, err error) {
	cmd := storeredis.GetClient().HSet(ctx, PipelinesKey(), d.Name, d.Val())
	i, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return i > 0, nil
}

// UpdatePipeline update pipeline info in etcd
func UpdatePipeline(ctx context.Context, name string, opts ...pipeline.OptionPipeline) error {
	return myDao.UpdatePipeline(ctx, name, opts...)
}

// AllPipelines get all info of pipelines from etcd in array form
func AllPipelines(ctx context.Context) (list []*pipeline.Pipeline, err error) {
	return myDao.AllPipelines(ctx)
}

// AllPipelinesMap returns all pipelines from etcd in map form
func AllPipelinesMap(ctx context.Context) (mapping map[string]*pipeline.Pipeline, err error) {
	return myDao.AllPipelinesMap(ctx)
}

// DeletePipeline delete pipeline info by name
func DeletePipeline(ctx context.Context, name string) (err error) {
	return storeredis.Default.Delete(ctx, &pipeline.Pipeline{Name: name})
}

// DeleteCompletePipeline delete pipeline, contains pipeline info, pipeline position
func DeleteCompletePipeline(ctx context.Context, name string) (err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}
	err = DeletePipeline(ctx, name)
	if err != nil {
		return
	}
	err = DeletePosition(ctx, name)
	return
}
