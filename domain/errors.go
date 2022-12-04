package domain

var (
	//Daycare Error
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
	ErrorMethodNotAllowed        = "method not allowed"

	//Appointment Error
	ErrorInvalidAppointmentData   = "invalid appointment data"
	ErrorInvalidAppointmentId     = "Appointment Id is already existed."
	ErrorAppointmentAlreadyExists = "Appointment already exists"
	ErrorAppointmentDoesNotExist  = "Appointment does not exist"
)
