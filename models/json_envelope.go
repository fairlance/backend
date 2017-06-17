package models

type JSONEnvelope struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}
