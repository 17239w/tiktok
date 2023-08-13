package models

// StatusCodeResponse：通用响应
type StatusCodeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
