package electrs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"

	testData "github.com/keep-network/keep-core/internal/testdata/electrs"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestConnect(t *testing.T) {
	testData := map[string]struct {
		config        Config
		expectedURL   string
		expectedError error
	}{
		"URL with no trailing slash": {
			config: Config{
				URL:            "example.org/api",
				RequestTimeout: DefaultRequestTimeout,
				RetryTimeout:   DefaultRetryTimeout,
			},
			expectedURL: "example.org/api",
		},
		"URL with trailing slash": {
			config: Config{
				URL:            "example.org/api/",
				RequestTimeout: DefaultRequestTimeout,
				RetryTimeout:   DefaultRetryTimeout,
			},
			expectedURL: "example.org/api",
		},
		"URL with protocol": {
			config: Config{
				URL:            "https://example.org",
				RequestTimeout: DefaultRequestTimeout,
				RetryTimeout:   DefaultRetryTimeout,
			},
			expectedURL: "https://example.org",
		},
		"URL with IP and port": {
			config: Config{
				URL:            "45.85.96.45:8596/api",
				RequestTimeout: DefaultRequestTimeout,
				RetryTimeout:   DefaultRetryTimeout,
			},
			expectedURL: "45.85.96.45:8596/api",
		},
		"non-default request timeout": {
			config: Config{
				URL:            testAPIURL,
				RequestTimeout: 5 * time.Minute,
				RetryTimeout:   DefaultRetryTimeout,
			},
			expectedURL: testAPIURL,
		},
		"non-default retry timeout": {
			config: Config{
				URL:            testAPIURL,
				RequestTimeout: DefaultRequestTimeout,
				RetryTimeout:   3 * time.Hour,
			},
			expectedURL: testAPIURL,
		},
		"URL not set": {
			config: Config{
				URL: "",
			},
			expectedError: fmt.Errorf("URL not set"),
		},
	}
	for testName, testData := range testData {
		t.Run(testName, func(t *testing.T) {
			connection, err := Connect(testData.config)
			if !reflect.DeepEqual(testData.expectedError, err) {
				t.Errorf(
					"unexpected error\nexpected: %v\nactual:   %v\n",
					testData.expectedError,
					err,
				)
			}

			if testData.expectedError == nil {
				testutils.AssertStringsEqual(
					t,
					"URL",
					testData.expectedURL,
					connection.(*Connection).url,
				)

				testutils.AssertIntsEqual(
					t,
					"RetryTimeout",
					int(testData.config.RetryTimeout),
					int(connection.(*Connection).retryTimeout),
				)

				testutils.AssertIntsEqual(
					t,
					"RequestTimeout",
					int(testData.config.RequestTimeout),
					int(connection.(*Connection).client.(*http.Client).Timeout),
				)
			}
		})
	}
}

func TestGetTransaction(t *testing.T) {
	transactionHashString := "c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479"
	transactionHash, err := bitcoin.NewHashFromString(
		transactionHashString,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	httpClientMock := newHTTPClientMock()
	httpClientMock.mockGet(
		fmt.Sprintf("%s/tx/%s", testAPIURL, transactionHash.String(bitcoin.ReversedByteOrder)),
		200,
		testData.Tx,
	)

	electrs := newMockedConnection(httpClientMock)

	expectedResult := bitcoinTestTx(t)

	result, err := electrs.GetTransaction(transactionHash)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(result, expectedResult); diff != nil {
		t.Errorf("compare failed: %v", diff)
	}
}

func TestGetTransactionConfirmations(t *testing.T) {
	currentBlock := 2403554
	expectedResult := uint(268506)

	transactionHashString := "c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479"
	transactionHash, err := bitcoin.NewHashFromString(
		transactionHashString,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	httpClientMock := newHTTPClientMock()
	httpClientMock.mockGet(
		fmt.Sprintf("%s/tx/%s", testAPIURL, transactionHash.String(bitcoin.ReversedByteOrder)),
		200,
		testData.Tx,
	)
	httpClientMock.mockGet(
		fmt.Sprintf("%s/blocks/tip/height", testAPIURL),
		200,
		fmt.Sprint(currentBlock),
	)

	electrs := newMockedConnection(httpClientMock)

	result, err := electrs.GetTransactionConfirmations(transactionHash)
	if err != nil {
		t.Fatal(err)
	}

	if result != expectedResult {
		t.Errorf(
			"invalid result\nexpected: %v\nactual:   %v",
			expectedResult,
			result,
		)
	}
}

func TestBroadcastTransaction(t *testing.T) {
	bitcoinTx := bitcoinTestTx(t)

	mockClient := newHTTPClientMock()
	mockClient.mockPost(
		fmt.Sprintf("%s/tx", testAPIURL),
		string(bitcoinTx.Serialize()),
		200,
		"fake-tx-id",
		t,
	)
	electrs := newMockedConnection(mockClient)

	if err := electrs.BroadcastTransaction(bitcoinTx); err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}
}

func TestGetLatestBlockHeight(t *testing.T) {
	expectedResult := uint(189645)

	mockClient := newHTTPClientMock()
	mockClient.mockGet(
		fmt.Sprintf("%s/blocks/tip/height", testAPIURL),
		200,
		fmt.Sprint(expectedResult),
	)

	electrs := newMockedConnection(mockClient)

	result, err := electrs.GetLatestBlockHeight()
	if err != nil {
		t.Fatal(err)
	}

	if result != expectedResult {
		t.Errorf(
			"invalid result\nexpected: %v\nactual:   %v",
			expectedResult,
			result,
		)
	}
}

func TestGetBlockHeader(t *testing.T) {
	blockHeight := uint(2135502)
	blockHash := "000000000000002af10911b8db32ed34dc6ea6515f84af5f7b82973c9a839e6d"

	mockClient := newHTTPClientMock()
	mockClient.mockGet(
		fmt.Sprintf("%s/block-height/%d", testAPIURL, blockHeight),
		200,
		blockHash,
	)
	mockClient.mockGet(
		fmt.Sprintf("%s/block/%s", testAPIURL, blockHash),
		200,
		testData.Block,
	)

	electrs := newMockedConnection(mockClient)

	previousBlockHeaderHash, err := bitcoin.NewHashFromString(
		"000000000066450030efdf72f233ed2495547a32295deea1e2f3a16b1e50a3a5",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	merkleRootHash, err := bitcoin.NewHashFromString(
		"1251774996b446f85462d5433f7a3e384ac1569072e617ab31e86da31c247de2",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := &bitcoin.BlockHeader{
		Version:                 536870916,
		PreviousBlockHeaderHash: previousBlockHeaderHash,
		MerkleRootHash:          merkleRootHash,
		Time:                    1641914003,
		Bits:                    436256810,
		Nonce:                   778087099,
	}

	result, err := electrs.GetBlockHeader(blockHeight)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(result, expectedResult); diff != nil {
		t.Errorf("compare failed: %v", diff)
	}
}

const testAPIURL = "example.org/api"

func newMockedConnection(client *httpClientMock) *Connection {
	electrs := &Connection{
		url:          testAPIURL,
		retryTimeout: 5 * time.Second,
		client:       client,
	}

	return electrs
}

func newHTTPClientMock() *httpClientMock {
	return &httpClientMock{
		getMocks:  make(map[string]func() (*http.Response, error)),
		postMocks: make(map[string]func(requestBody io.Reader) (*http.Response, error)),
	}
}

type httpClientMock struct {
	getMocks  map[string]func() (*http.Response, error)
	postMocks map[string]func(requestBody io.Reader) (*http.Response, error)
}

func (m httpClientMock) Get(url string) (*http.Response, error) {
	mock, ok := m.getMocks[url]
	if !ok {
		return nil, fmt.Errorf("mocked get request not registered for url [%s]", url)
	}

	return mock()
}

func (m httpClientMock) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	// Define TeeReader to read from body io.Reader multiple times.
	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)

	bodyBytes, err := ioutil.ReadAll(tee)
	if err != nil {
		return nil, fmt.Errorf("failed to read the request body: [%w]", err)
	}

	mock, ok := m.postMocks[fmt.Sprintf("%s-%s", url, bodyBytes)]
	if !ok {
		return nil, fmt.Errorf(
			"mocked post request not registered for url [%s] and body [%s]",
			url,
			bodyBytes,
		)
	}

	return mock(&buf)
}

func (m httpClientMock) mockGet(expectedURL string, responseStatusCode int, responseBody string) {
	m.getMocks[expectedURL] = func() (*http.Response, error) {
		return mockResponse(responseStatusCode, responseBody), nil
	}
}

func (m httpClientMock) mockPost(expectedURL string, expectedRequestBody string, responseStatusCode int, responseBody string, t *testing.T) {
	m.postMocks[fmt.Sprintf("%s-%s", expectedURL, expectedRequestBody)] = func(body io.Reader) (*http.Response, error) {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if string(bodyBytes) != expectedRequestBody {
			t.Fatalf(
				"unexpected request body\nexpected: %s\nactual:   %s",
				expectedRequestBody,
				bodyBytes,
			)
		}

		return mockResponse(responseStatusCode, responseBody), nil
	}
}

func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}
}

func bitcoinTestTx(t *testing.T) *bitcoin.Transaction {
	prevTxHash, err := bitcoin.NewHashFromString(
		"e788a344a86f7e369511fe37ebd1d74686dde694ee99d06db5db3d4a14719b1d",
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	return &bitcoin.Transaction{
		Version: 1,
		Inputs: []*bitcoin.TransactionInput{
			{
				Outpoint: &bitcoin.TransactionOutpoint{
					TransactionHash: prevTxHash,
					OutputIndex:     1,
				},
				SignatureScript: []byte("47304402206f8553c07bcdc0c3b906311888103d623ca9096ca0b28b7d04650a029a01fcf9022064cda02e39e65ace712029845cfcf58d1b59617d753c3fd3556f3551b609bbb00121039d61d62dcd048d3f8550d22eb90b4af908db60231d117aeede04e7bc11907bfa"),
				Sequence:        4294967295,
			},
		},
		Outputs: []*bitcoin.TransactionOutput{
			{
				PublicKeyScript: []byte("a9143ec459d0f3c29286ae5df5fcc421e2786024277e87"),
				Value:           20000,
			},
			{
				PublicKeyScript: []byte("0014e257eccafbc07c381642ce6e7e55120fb077fbed"),
				Value:           1360550,
			},
		},
		Locktime: 0,
	}
}
