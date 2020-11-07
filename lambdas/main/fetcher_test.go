package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func Test_found_url_should_return_ok(t *testing.T) {
	url := Url{
		ShortenedURL: "shortened-id",
		OriginalURL:  "www.original-url.com",
	}
	result := &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"shortened_url": {
				S: aws.String(url.ShortenedURL),
			},
			"original_url": {
				S: aws.String(url.OriginalURL),
			},
		},
	}

	repository := new(MockRepository)
	repository.On("GetItem", mock.Anything, mock.Anything).Return(result, nil)
	dynamo = repository

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"shortened_url": url.ShortenedURL,
		},
	}

	response, err := handler(request)

	repository.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	responseURL := Url{}
	json.Unmarshal([]byte(response.Body), &responseURL)
	assert.Equal(t, url.OriginalURL, responseURL.OriginalURL)
	assert.Equal(t, url.ShortenedURL, responseURL.ShortenedURL)
}

func Test_not_found_url_should_return_not_found(t *testing.T) {
	url := Url{
		ShortenedURL: "shortened-id",
		OriginalURL:  "www.original-url.com",
	}
	result := &dynamodb.GetItemOutput{Item: nil}

	repository := new(MockRepository)
	repository.On("GetItem", mock.Anything, mock.Anything).Return(result, nil)
	dynamo = repository

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"shortened_url": url.ShortenedURL,
		},
	}

	response, err := handler(request)

	repository.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Empty(t, response.Body)
}

func Test_error_fetching_database_should_return_internal_server_error(t *testing.T) {
	url := Url{
		ShortenedURL: "shortened-id",
		OriginalURL:  "www.original-url.com",
	}
	result := &dynamodb.GetItemOutput{}

	expectedError := errors.New("error")
	repository := new(MockRepository)
	repository.On("GetItem", mock.Anything, mock.Anything).Return(result, expectedError)
	dynamo = repository

	request := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"shortened_url": url.ShortenedURL,
		},
	}

	response, err := handler(request)

	repository.AssertExpectations(t)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Empty(t, response.Body)
}
