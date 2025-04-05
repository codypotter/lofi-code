package db

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func (db *Db) NewVerificationToken(email string, ttlSeconds int64) (string, error) {
	ctx := context.Background()

	token := generateToken()

	_, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("blog"),
		Item: map[string]types.AttributeValue{
			"pk":  &types.AttributeValueMemberS{Value: "EMAIL_TOKEN#" + token},
			"sk":  &types.AttributeValueMemberS{Value: "EMAIL#" + email},
			"ttl": &types.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix()+ttlSeconds, 10)},
		},
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateToken() string {
	return uuid.New().String()
}

func (db *Db) VerifyEmail(token string) (string, error) {
	ctx := context.Background()
	pk := "EMAIL_TOKEN#" + token

	resp, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("blog"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return "", err
	}
	if len(resp.Items) == 0 {
		return "", errors.New("invalid or expired token")
	}

	item := resp.Items[0]
	sk := item["sk"].(*types.AttributeValueMemberS).Value
	email := strings.TrimPrefix(sk, "EMAIL#")

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("blog"),
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "MAILING_LIST"},
			"sk": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return "", err
	}

	_, _ = db.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("blog"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})

	return email, nil
}

func (db *Db) Subscribe(email string) error {
	ctx := context.Background()

	pk := "MAILING_LIST"
	sk := email

	_, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("blog"),
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
