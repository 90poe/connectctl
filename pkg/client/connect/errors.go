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

// IsAPIError indicates if the error is an struct of type APIError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// IsNotFound indicates if the error represents an HTTP 404 status code
func IsNotFound(err error) bool {
	apiErr, ok := err.(*APIError)
	if !ok {
		return false
	}
	return apiErr.Code == http.StatusNotFound
}

// IsRetryable indicates if the error could be retryed.
// See https://github.com/apache/kafka/blob/master/connect/runtime/src/main/java/org/apache/kafka/connect/runtime/rest/resources/ConnectorsResource.java#L299-L325
func IsRetryable(err error) bool {
	apiErr, ok := err.(*APIError)
	if !ok {
		return false
	}
	return apiErr.Code == http.StatusConflict
}
