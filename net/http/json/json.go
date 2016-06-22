package json

import (
	"encoding/json"
	"net/http"
)

// APIResponse is a generic API response type.
type APIResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  interface{}            `json:"result"`
	Data    map[string]interface{} `json:"data"` // Generic extra data to be sent along in response
}

// WriteResponse writes a marshalled reponses using JSON.
func WriteResponse(w http.ResponseWriter, code int, resp interface{}) error {
	j, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(j)
	return err
}
