package nats_connector

type NetworkResponse struct {
	Data     []byte `json:"data"`
	HTTPCode int    `json:"HTTPCode"`
}

func NewResponse(data []byte, code int) NetworkResponse {
	return NetworkResponse{data, code}
}
