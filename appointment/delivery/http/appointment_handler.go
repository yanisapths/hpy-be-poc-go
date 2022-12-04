package http

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/yanisapths/happyelders-api/appointment/repository"
	"github.com/yanisapths/happyelders-api/domain"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetAppointment(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["id"]
	if len(id) > 0 {
		result, err := repository.FetchAppointment(id, appointmentDetailsTable, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := repository.FetchAppointments(appointmentDetailsTable, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func CreateAppointment(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := repository.CreateAppointment(req, appointmentDetailsTable, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

func UpdateAppointment(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	result, err := repository.UpdateAppointment(req, appointmentDetailsTable, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

func Delete(req events.APIGatewayProxyRequest, appointmentDetailsTable string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := repository.DeleteAppointment(req, appointmentDetailsTable, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, domain.ErrorMethodNotAllowed)
}
