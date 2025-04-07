package db

import (
	"context"

	loficodeconfig "loficode/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Db struct {
	client *dynamodb.Client
}

func New(ctx context.Context, lc *loficodeconfig.Config) *Db {
	return &Db{
		client: dynamodb.NewFromConfig(lc.AwsConfig),
	}
}

func NewDevelopment(ctx context.Context, lc *loficodeconfig.Config) *Db {
	return &Db{
		client: dynamodb.NewFromConfig(lc.AwsConfig, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://localhost:8000")
		}),
	}
}
