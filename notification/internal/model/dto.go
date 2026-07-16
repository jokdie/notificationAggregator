package model

type NotificationRequest struct {
	UserID   int      `json:"userId" validate:"gt=0"`
	Message  string   `json:"message" validate:"required,min=1,max=125"`
	Channels []string `json:"channels validate "`
}

type NotificationResponse struct {
	Results []struct {
		Channel string `json:"channel"`
		Status  string `json:"status"`
	} `json:"results"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
