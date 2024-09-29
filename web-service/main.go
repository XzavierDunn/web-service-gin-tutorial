package main

import (
	"context"
	"function/src/router"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	router := router.SetupRouter()
	ginLambda = ginadapter.New(router)
}

func main() {
	if os.Getenv("ENV") == "local" {
		router := router.SetupRouter()
		router.Run(":8080")
	} else {
		lambda.Start(handler)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
