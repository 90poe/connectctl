// NOTE: Code originally from https://github.com/go-kafka/connect

package connect

import (
	"fmt"
	"net/http"
)

// APIError holds information returned from a Kafka Connect API instance about
// why an API call failed.
type APIError struct {
	Code     int            `json:"error_code"`
	Message  string         `json:"message"`
	Response *http.Response // HTTP response that caused this error
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v (HTTP %d)", e.Message, e.Code)
}

func IsAPIError(err error) bool {
	_, ok := err.(APIError)
	return ok
}

func IsNotFound(err error) bool {
	apiErr, ok := err.(APIError)
	if !ok {
		return false
	}

	return apiErr.Code == 404
}
