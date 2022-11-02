package electrs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/keep-network/keep-common/pkg/wrappers"
)

func (c *Connection) httpGet(resource string) (io.ReadCloser, error) {
	return c.httpRequest(
		func() (*http.Response, error) {
			return c.client.Get(
				fmt.Sprintf("%s/%s", c.url, resource),
			)
		},
	)
}

func (c *Connection) httpPost(resource string, requestBody []byte) (io.ReadCloser, error) {
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
) (io.ReadCloser, error) {
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
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Warnf(
					"failed to read the response: [%v]",
					err,
				)
			}

			return fmt.Errorf(
				"response status is not 200 - status: [%s], payload: [%s]",
				resp.Status,
				responseBody,
			)
		}

		responseReader = resp.Body

		return nil
	})

	return responseReader, err
}
