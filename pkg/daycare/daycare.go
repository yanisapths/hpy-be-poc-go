package daycare

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/yanisapths/happyelders-api/pkg/validators"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidDaycareData      = "invalid daycare data"
	ErrorInvalidEmail            = "invalid email"
	ErrorInvalidDaycareName      = "Daycare name is already taken."
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item"
	ErrorDaycareAlreadyExists    = "daycare.Daycare already exists"
	ErrorDaycareDoesNotExist     = "daycare.Daycare does not exist"
)

type Daycare struct {
	DaycareName string `json:"name"`
	Address     string `json:"address"`
	Owner       string `json:"owner"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

func FetchDaycare(name, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Daycare, error) {
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
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Daycare)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchDaycares(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Daycare, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]Daycare)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateDaycare(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Daycare,
	error,
) {
	var d Daycare

	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDaycareData)
	}

	if !validators.IsEmailValid(d.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentDaycare, _ := FetchDaycare(d.DaycareName, tableName, dynaClient)
	if currentDaycare != nil && len(currentDaycare.DaycareName) != 0 {
		return nil, errors.New(ErrorDaycareAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(d)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}

func UpdateDaycare(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Daycare,
	error,
) {
	var d Daycare
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDaycareName)
	}

	currentDaycare, _ := FetchDaycare(d.DaycareName, tableName, dynaClient)
	if currentDaycare != nil && len(currentDaycare.DaycareName) == 0 {
		return nil, errors.New(ErrorDaycareDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(d)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
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
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
