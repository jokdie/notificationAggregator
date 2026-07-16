package model

type ProviderRequest struct {
	UserID  int    `json:"userId" validate:"gt=0"`
	Message string `json:"message" validate:"required,min=1,max=125"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
