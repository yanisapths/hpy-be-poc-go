package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	_daycareHttp "github.com/yanisapths/happyelders-api/daycare/delivery/http"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})
	if err != nil {
		return
	}

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableName = "daycareDetails"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return _daycareHttp.GetDaycare(req, tableName, dynaClient)
	case "POST":
		return _daycareHttp.CreateDaycare(req, tableName, dynaClient)
	case "PUT":
		return _daycareHttp.UpdateDaycare(req, tableName, dynaClient)
	case "DELETE":
		return _daycareHttp.Delete(req, tableName, dynaClient)
	}
	return _daycareHttp.UnhandledMethod()
}
