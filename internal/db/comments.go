package db

import (
	"context"
	"fmt"
	"loficode/internal/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func (db *Db) GetCommentsBySlug(slug string) ([]model.Comment, error) {
	ctx := context.Background()
	var comments []model.Comment

	pk := "COMMENTS#" + slug

	result, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("blog"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		comment := model.Comment{
			Name:  item["name"].(*types.AttributeValueMemberS).Value,
			Email: item["email"].(*types.AttributeValueMemberS).Value,
			Text:  item["text"].(*types.AttributeValueMemberS).Value,
		}
		if dateStr, ok := item["date"].(*types.AttributeValueMemberS); ok {
			comment.Date, err = time.Parse(time.RFC3339, dateStr.Value)
			if err != nil {
				return nil, err
			}
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (db *Db) AddComment(slug string, comment model.Comment) error {
	ctx := context.Background()

	pk := "COMMENTS#" + slug
	sk := fmt.Sprintf("COMMENT#%s#%s", comment.Date.Format(time.RFC3339), uuid.NewString())

	_, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("blog"),
		Item: map[string]types.AttributeValue{
			"pk":    &types.AttributeValueMemberS{Value: pk},
			"sk":    &types.AttributeValueMemberS{Value: sk},
			"name":  &types.AttributeValueMemberS{Value: comment.Name},
			"email": &types.AttributeValueMemberS{Value: comment.Email},
			"text":  &types.AttributeValueMemberS{Value: comment.Text},
			"date":  &types.AttributeValueMemberS{Value: comment.Date.Format(time.RFC3339)},
		},
	})
	return err
}
