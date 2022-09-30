package mock

import (
	"encoding/json"
	"errors"
	"net/http"
)

// MustMarshal helper function that wraps json.Marshal
func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)

	if err == nil {
		return b
	}

	panic(err)
}

// WriteError helper function to write errors to HTTP handlers
func WriteError(
	w http.ResponseWriter,
	httpStatus int,
	msg string,
) {
	w.WriteHeader(httpStatus)

	w.Write(MustMarshal(errors.New(msg)))
}
