package db

import (
	"context"
	"loficode/internal/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

func (db *Db) UpsertPosts(ctx context.Context, posts []model.Post) {
	for _, post := range posts {
		tags := append(post.Tags, "all")
		for _, tag := range tags {
			db.upsertPost(ctx, tag, post)
			log.Debug().Str("slug", post.Slug).Str("tag", tag).Msg("Upserted post")
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
		log.Fatal().Str("slug", post.Slug).Str("tag", tag).Msg("Failed to upsert post")
	}
}

func (db *Db) GetPostsByTag(ctx context.Context, tag string, cursor string) ([]model.Post, *string, error) {
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
