package models

type ResponseModel struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type LoginResponse struct {
	ResponseModel
	Token string `json:"token"`
}
