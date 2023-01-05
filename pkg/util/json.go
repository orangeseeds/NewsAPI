package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// helpers
func DecodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func encodeBody(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func Respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, data)
	}
}

func RespondErr(w http.ResponseWriter, status int, data any) {
	var resp any
	resp, ok := data.(error)
	if ok {
		resp = fmt.Sprint(resp)
	} else {
		resp = data
	}

	Respond(w, status, map[string]interface{}{
		"data": []any{},
		"error": map[string]any{
			"message": resp,
		},
		"success": false,
	})
}

func RespondHTTPErr(w http.ResponseWriter, status int) {
	RespondErr(w, status, http.StatusText(status))
}