package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

func (db *Db) CreateTable() {
	_, err := db.client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String("blog"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		var exists *types.ResourceInUseException
		if errors.As(err, &exists) {
			log.Info().Msg("Table already exists. Skipping creation.")
		} else {
			log.Fatal().Err(err).Msg("Failed to create table")
		}
	}

	_, err = db.client.UpdateTimeToLive(context.Background(), &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String("blog"),
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			AttributeName: aws.String("ttl"),
			Enabled:       aws.Bool(true),
		},
	})
	if err != nil {
		log.Info().Msg("TTL is already enabled. Skipping update.")
	}
}
