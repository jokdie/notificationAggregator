package model

type ProviderRequest struct {
	UserID  int    `json:"userId" validate:"gt=0"`
	Message string `json:"message" validate:"required,min=1,max=125"`
}

type ProviderResponse struct {
	Channel string `json:"channel"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
