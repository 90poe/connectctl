package connect

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsRetryable(t *testing.T) {

	testCases := []struct {
		name   string
		req    *http.Request
		resp   *http.Response
		result bool
	}{
		{
			name:   "401",
			result: true,
			resp: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{ "error_code" : 409 }`)),
				StatusCode: http.StatusConflict,
			},
			req: &http.Request{},
		},
		{
			name:   "404",
			result: false,
			resp: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{ "error_code" : 404 }`)),
				StatusCode: http.StatusNotFound,
			},
			req: &http.Request{},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			err := buildError(tc.req, tc.resp)
			require.True(t, IsAPIError(err))

			result := IsRetryable(err)
			require.Equal(t, tc.result, result)

		})
	}
}

func Test_IsNotFound(t *testing.T) {

	testCases := []struct {
		name   string
		req    *http.Request
		resp   *http.Response
		result bool
	}{
		{
			name:   "404",
			result: true,
			resp: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{ "error_code" : 404 }`)),
				StatusCode: http.StatusNotFound,
			},
			req: &http.Request{},
		},
		{
			name: "401",
			resp: &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{ "error_code" : 409 }`)),
				StatusCode: http.StatusConflict,
			},
			req: &http.Request{},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			err := buildError(tc.req, tc.resp)
			require.True(t, IsAPIError(err))

			result := IsNotFound(err)
			require.Equal(t, tc.result, result)

		})
	}
}
