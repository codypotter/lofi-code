package db

import (
	"context"
	"errors"
	"log"
	"time"

	loficodeconfig "loficode/internal/config"
	"loficode/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Db struct {
	client *dynamodb.Client
}

func New(ctx context.Context, lc *loficodeconfig.Config) *Db {
	return &Db{
		client: dynamodb.NewFromConfig(lc.AwsConfig, func(o *dynamodb.Options) {
			if lc.Environment == "development" {
				o.BaseEndpoint = aws.String("http://localhost:8000")
			}
		}),
	}
}

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
			log.Println("Table already exists. Skipping creation.")
		} else {
			log.Fatalf("failed to create table: %v", err)
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
		log.Println("Failed to update TTL:", err)
	}
}

func (db *Db) UpsertPosts(posts []model.Post) {
	ctx := context.Background()

	for _, post := range posts {
		tags := append(post.Tags, "all")
		for _, tag := range tags {
			db.upsertPost(ctx, tag, post)
			log.Println("Upserted post:", post.Slug, "for tag:", tag)
		}
	}
}

func (db *Db) upsertPost(ctx context.Context, tag string, post model.Post) {
	pk := "TAG#" + tag
	sk := "POST#" + post.Date.Format(time.RFC3339) + "#" + post.Slug

	item := map[string]types.AttributeValue{
		"pk":             &types.AttributeValueMemberS{Value: pk},
		"sk":             &types.AttributeValueMemberS{Value: sk},
		"slug":           &types.AttributeValueMemberS{Value: post.Slug},
		"title":          &types.AttributeValueMemberS{Value: post.Title},
		"summary":        &types.AttributeValueMemberS{Value: post.Summary},
		"date":           &types.AttributeValueMemberS{Value: post.Date.Format(time.RFC3339)},
		"tags":           &types.AttributeValueMemberSS{Value: post.Tags},
		"openGraphImage": &types.AttributeValueMemberS{Value: post.OpenGraphImage},
	}

	_, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("blog"),
		Item:      item,
	})
	if err != nil {
		log.Fatalf("failed to upsert post %s for tag %s: %v", post.Slug, tag, err)
	}
}

func (db *Db) GetPostsByTag(tag string, cursor string) ([]model.Post, *string, error) {
	ctx := context.Background()
	var posts []model.Post
	var nextCursor *string

	pk := "TAG#" + tag
	params := &dynamodb.QueryInput{
		TableName:              aws.String("blog"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	}

	if cursor != "" {
		params.ExclusiveStartKey = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: cursor},
		}
	}

	result, err := db.client.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	for _, item := range result.Items {
		post := model.Post{
			Slug:           item["slug"].(*types.AttributeValueMemberS).Value,
			Title:          item["title"].(*types.AttributeValueMemberS).Value,
			Summary:        item["summary"].(*types.AttributeValueMemberS).Value,
			Date:           time.Time{},
			Tags:           item["tags"].(*types.AttributeValueMemberSS).Value,
			OpenGraphImage: item["openGraphImage"].(*types.AttributeValueMemberS).Value,
		}
		dateStr := item["date"].(*types.AttributeValueMemberS).Value
		post.Date, _ = time.Parse(time.RFC3339, dateStr)
		posts = append(posts, post)
	}

	if result.LastEvaluatedKey != nil {
		if skAttr, ok := result.LastEvaluatedKey["sk"].(*types.AttributeValueMemberS); ok {
			nextCursor = aws.String(skAttr.Value)
		}
	}

	return posts, nextCursor, nil
}
