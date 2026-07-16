package model

type NotificationRequest struct {
	UserID   int      `json:"userId" validate:"gt=0"`
	Message  string   `json:"message" validate:"required,min=1,max=125"`
	Channels []string `json:"channels" validate:"required,min=1,dive,oneof=email sms push"`
}

type NotificationResult struct {
	Channel string `json:"channel"`
	Status  string `json:"status"`
}

type NotificationResponse struct {
	Results []NotificationResult `json:"results"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
