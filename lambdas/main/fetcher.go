package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/eaneto/serverless-url-shortener/client"
	log "github.com/sirupsen/logrus"
)

type Url struct {
	ShortenedURL string `json:"shortened_url"`
	OriginalURL  string `json:"original_url"`
}

var dynamo client.Repository

func init() {
	session, _ := session.NewSession()
	dynamo = client.DynamoRepository{
		Client: dynamodb.New(session),
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	shortenedUrl := request.PathParameters["shortened_url"]
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("url"),
		Key: map[string]*dynamodb.AttributeValue{
			"shortened_url": {
				S: aws.String(shortenedUrl),
			},
		},
	})
	if err != nil {
		log.Error("Error reading table", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	url := Url{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &url)
	if err != nil {
		log.Error("Error unmarshaling found result", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	payload, err := json.Marshal(url)
	if err != nil {
		log.Error("Error marshaling response", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(payload),
	}, nil
}
