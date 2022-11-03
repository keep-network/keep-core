package electrs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/keep-network/keep-common/pkg/wrappers"
)

func (c *Connection) httpGet(resource string) ([]byte, error) {
	return c.httpRequest(
		func() (*http.Response, error) {
			return c.client.Get(
				fmt.Sprintf("%s/%s", c.url, resource),
			)
		},
	)
}

func (c *Connection) httpPost(resource string, requestBody []byte) ([]byte, error) {
	return c.httpRequest(
		func() (*http.Response, error) {
			return c.client.Post(
				fmt.Sprintf("%s/%s", c.url, resource),
				"text/plain",
				bytes.NewReader(requestBody),
			)
		},
	)
}

func (c *Connection) httpRequest(
	request func() (*http.Response, error),
) ([]byte, error) {
	var responseReader io.ReadCloser

	err := wrappers.DoWithDefaultRetry(c.retryTimeout, func(ctx context.Context) error {
		responseReader = nil

		resp, err := request()
		if err != nil {
			return fmt.Errorf(
				"failed to submit a request: [%w]",
				err,
			)
		}

		if resp.StatusCode != 200 {
			responseBody, _ := io.ReadAll(resp.Body)
			// Ignore error from body reading as returning the response body
			// is just a nice-to-have.

			return fmt.Errorf(
				"response status is not 200 - status: [%s], payload: [%s]",
				resp.Status,
				responseBody,
			)
		}

		responseReader = resp.Body

		return nil
	})
	if err != nil {
		return []byte{}, fmt.Errorf("request failed: [%w]", err)
	}

	responseBody, err := io.ReadAll(responseReader)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: [%w]", err)
	}

	return responseBody, nil
}
