package models

type UpdateCustomerInput struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	FunnelID    *uint   `json:"funnel_id"`
	FunnelStage string  `json:"funnel_stage"`
}
