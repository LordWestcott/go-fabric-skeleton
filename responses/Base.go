package responses

import (
	"encoding/json"
	"net/http"
)

type ResponseBase struct {
	Success bool        `json:"success"`
	Error   string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponseBase(success bool, err string, data interface{}) *ResponseBase {
	return &ResponseBase{
		Success: success,
		Error:   err,
		Data:    data,
	}
}

func (r *ResponseBase) WriteJson(w http.ResponseWriter) {
	enc := json.NewEncoder(w)
	enc.Encode(r)
}
