package domain

// Daycare type
type Daycare struct {
	DaycareName    string `json:"name"`
	Address        string `json:"address"`
	Owner          string `json:"owner"`
	PhoneNumber    string `json:"phoneNumber"`
	Email          string `json:"email"`
	ImageUrl       string `json:"imageUrl"`
	ApprovalStatus string `json:"approvalStatus"`
}
