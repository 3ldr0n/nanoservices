package client

import "github.com/aws/aws-sdk-go/service/dynamodb"

var dynamoClient *dynamodb.DynamoDB

type Repository interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

type DynamoRepository struct {
	Client *dynamodb.DynamoDB
}

func (repository DynamoRepository) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return repository.Client.GetItem(input)
}
