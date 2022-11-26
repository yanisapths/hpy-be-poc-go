package repository

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/yanisapths/happyelders-api/domain"
)

func FetchDaycare(name, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*domain.Daycare, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(domain.ErrorFailedToFetchRecord)
	}

	item := new(domain.Daycare)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(domain.ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchDaycares(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]domain.Daycare, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(domain.ErrorFailedToFetchRecord)
	}
	item := new([]domain.Daycare)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateDaycare(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*domain.Daycare,
	error,
) {
	var d domain.Daycare

	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(domain.ErrorInvalidDaycareData)
	}

	currentDaycare, _ := FetchDaycare(d.DaycareName, tableName, dynaClient)
	if currentDaycare != nil && len(currentDaycare.DaycareName) != 0 {
		return nil, errors.New(domain.ErrorDaycareAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(d)

	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}

func UpdateDaycare(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*domain.Daycare,
	error,
) {
	var d domain.Daycare
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(domain.ErrorInvalidDaycareName)
	}

	currentDaycare, _ := FetchDaycare(d.DaycareName, tableName, dynaClient)
	if currentDaycare != nil && len(currentDaycare.DaycareName) == 0 {
		return nil, errors.New(domain.ErrorDaycareDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}

func DeleteDaycare(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	name := req.QueryStringParameters["name"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(domain.ErrorCouldNotDeleteItem)
	}

	return nil
}
