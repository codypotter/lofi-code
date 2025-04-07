package main

import (
	"context"
	"loficode/internal/application"
	"loficode/internal/config"
	"loficode/internal/logger"
	"loficode/internal/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/rs/zerolog/log"
)

// chiLambda is the global variable that holds the chi adapter
// for the AWS Lambda function. It is initialized in the init function
// once per cold start. This allows us to cache the router between
// warm invocations.
var chiLambda *chiadapter.ChiLambda

func init() {
	ctx := context.Background()
	cfg := config.New(ctx)

	logger.Configure(cfg.LogLevel, logger.LogFormatJson)

	app := application.New(ctx, cfg)
	r := router.NewRouter(app)

	chiLambda = chiadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("Handler invoked")
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	log.Debug().Msg("Starting Lambda function")
	lambda.Start(handler)
}
