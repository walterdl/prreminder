package sfn

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

var cache *sfn.Client

func New() (*sfn.Client, error) {
	if cache != nil {
		return cache, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	cache = sfn.NewFromConfig(cfg)
	return cache, nil
}
