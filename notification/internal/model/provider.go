package model

type ProviderRequest struct {
	UserID  int    `json:"userId"`
	Message string `json:"message"`
}

type ProviderErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
