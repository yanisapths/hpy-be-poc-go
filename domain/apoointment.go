package domain

// Appointment type
type Appointment struct {
	AppointmentId string `json:"id"`
	CustomerName  string `json:"customerName"`
	Date          string `json:"date"`
	PhoneNumber   string `json:"phoneNumber"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	Name          string `json:"name"`
	Status        string `json:"status"`
}
