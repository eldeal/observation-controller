package observations

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrInvalidObservationAPIResponse struct {
	actualCode int
	uri        string
	body       string
}

// Error should be called by the user to print out the stringified version of the error
func (e ErrInvalidObservationAPIResponse) Error() string {
	return fmt.Sprintf("invalid response: %d from observation api: %s, body: %s",
		e.actualCode,
		e.uri,
		e.body,
	)
}

// Code returns the status code received from observation api if an error is returned
func (e ErrInvalidObservationAPIResponse) Code() int {
	return e.actualCode
}

// NewDatasetAPIResponse creates an error response, optionally adding body to e when status is 404
func NewObservationAPIResponse(resp *http.Response, uri string) (e *ErrInvalidObservationAPIResponse) {
	e = &ErrInvalidObservationAPIResponse{
		actualCode: resp.StatusCode,
		uri:        uri,
	}
	if resp.StatusCode == http.StatusNotFound {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			e.body = "Client failed to read ObservationAPI body"
			return
		}
		defer closeResponseBody(nil, resp)

		e.body = string(b)
	}
	return
}
