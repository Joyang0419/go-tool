package response

type TResponse struct {
	// Success indicates whether the request is successful.
	Success bool `json:"success"`
	// ErrorCode indicates the error code of the request.
	ErrorCode int `json:"errorCode"`
	// Message indicates the error message of the request.
	Message string `json:"message"`
	// Data indicates the data of the request.
	Data interface{} `json:"data"`
}
