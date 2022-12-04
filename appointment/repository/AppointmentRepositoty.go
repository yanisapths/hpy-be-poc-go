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

func FetchAppointment(id, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (*domain.Appointment, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(appointmentDetailsTable),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(domain.ErrorFailedToFetchRecord)
	}

	item := new(domain.Appointment)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(domain.ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchAppointments(appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (*[]domain.Appointment, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(appointmentDetailsTable),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(domain.ErrorFailedToFetchRecord)
	}
	item := new([]domain.Appointment)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateAppointment(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (
	*domain.Appointment,
	error,
) {
	var a domain.Appointment

	if err := json.Unmarshal([]byte(req.Body), &a); err != nil {
		return nil, errors.New(domain.ErrorInvalidAppointmentData)
	}

	currentAppointment, _ := FetchAppointment(a.AppointmentId, appointmentDetailsTable, dynaClient)
	if currentAppointment != nil && len(currentAppointment.AppointmentId) != 0 {
		return nil, errors.New(domain.ErrorInvalidAppointmentData)
	}

	av, err := dynamodbattribute.MarshalMap(a)

	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(appointmentDetailsTable),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotDynamoPutItem)
	}
	return &a, nil
}

func UpdateAppointment(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (
	*domain.Appointment,
	error,
) {
	var a domain.Appointment
	if err := json.Unmarshal([]byte(req.Body), &a); err != nil {
		return nil, errors.New(domain.ErrorInvalidAppointmentId)
	}

	currentAppointment, _ := FetchAppointment(a.AppointmentId, appointmentDetailsTable, dynaClient)
	if currentAppointment != nil && len(currentAppointment.AppointmentId) == 0 {
		return nil, errors.New(domain.ErrorAppointmentDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(a)
	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(appointmentDetailsTable),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(domain.ErrorCouldNotDynamoPutItem)
	}
	return &a, nil
}

func DeleteAppointment(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) error {

	id := req.QueryStringParameters["id"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(appointmentDetailsTable),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(domain.ErrorCouldNotDeleteItem)
	}

	return nil
}
