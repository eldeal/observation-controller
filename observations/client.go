package observations

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ONSdigital/dp-api-clients-go/clientlog"
	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-observation-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

const service = "observations-api"

type Client struct {
	hcCli *healthcheck.Client
}

// NewAPIClient creates a new instance of Client with a given dataset api url and the relevant tokens
func NewAPIClient(datasetAPIURL string) *Client {
	return &Client{
		healthcheck.NewClient(service, datasetAPIURL),
	}
}

// NewWithHealthClient creates a new instance of Client,
// reusing the URL and Clienter from the provided health check client.
func NewWithHealthClient(hcCli *healthcheck.Client) *Client {
	return &Client{
		healthcheck.NewClientWithClienter(service, hcCli.URL, hcCli.Client),
	}
}

// Checker calls dataset api health endpoint and returns a check object to the caller.
func (c *Client) Checker(ctx context.Context, check *health.CheckState) error {
	return c.hcCli.Checker(ctx, check)
}

func (c *Client) do(ctx context.Context, uri string, values url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	if values != nil {
		req.URL.RawQuery = values.Encode()
	}

	return c.hcCli.Client.Do(ctx, req)
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(ctx context.Context, resp *http.Response) {
	if resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			log.Error(ctx, "error closing http response body", err)
		}
	}
}

// Get returns dataset level information for a given dataset id
func (c *Client) Get(ctx context.Context, datasetID, edition, version string, q url.Values) (m models.ObservationsDoc, err error) {
	uri := fmt.Sprintf("%s/datasets/%s/editions/%s/versions/%s/observations", c.hcCli.URL, datasetID, edition, version)

	clientlog.Do(ctx, "retrieving dataset", service, uri)

	resp, err := c.do(ctx, uri, q)
	if err != nil {
		return
	}

	defer closeResponseBody(ctx, resp)

	if resp.StatusCode != http.StatusOK {
		err = NewObservationAPIResponse(resp, uri)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &m)
	return
}
