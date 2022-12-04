package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	_appointmentHttp "github.com/yanisapths/happyelders-api/appointment/delivery/http"
	_daycareHttp "github.com/yanisapths/happyelders-api/daycare/delivery/http"
	"github.com/yanisapths/happyelders-api/domain"
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
const appointmentDetailsTable = "appointmentDetails"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/appointment":
		switch req.HTTPMethod {
		case "GET":
			return _appointmentHttp.GetAppointment(req, appointmentDetailsTable, dynaClient)
		case "POST":
			return _appointmentHttp.CreateAppointment(req, appointmentDetailsTable, dynaClient)
		case "PUT":
			return _appointmentHttp.UpdateAppointment(req, appointmentDetailsTable, dynaClient)
		case "DELETE":
			return _appointmentHttp.Delete(req, appointmentDetailsTable, dynaClient)
		}
		return _appointmentHttp.UnhandledMethod()
	case "/daycare":
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
	return apiResponse(http.StatusMethodNotAllowed, domain.ErrorMethodNotAllowed)
}

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
