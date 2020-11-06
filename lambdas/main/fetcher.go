package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func handler(cts context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}
