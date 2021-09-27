package api_client_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type APIClient struct {
	URL        string
	httpClient *http.Client
}

func NewAPIClient(url string, timeout time.Duration) APIClient {
	return APIClient{
		URL: url,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (apiClient APIClient) ToUpper(input string) (string, error) {
	req, err := http.NewRequest("GET", apiClient.URL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Set("input", input)
	req.URL.RawQuery = q.Encode()

	resp, err := apiClient.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func TestValid(t *testing.T) {
	input := "expected"

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(strings.ToUpper(input)))
		}),
	)

	defer s.Close()

	apiClient := NewAPIClient(s.URL, 5*time.Second)

	actual, err := apiClient.ToUpper(input)
	assert.NoError(t, err, "error")
	assert.Equal(t, strings.ToUpper(input), actual, "actual")
}

// You need to have the http server running
func TestRootURL(t *testing.T) {
	apiClient := NewAPIClient("http://localhost:8080", 5*time.Second)
	actual, err := apiClient.ToUpper("hello")
	assert.NoError(t, err)
	assert.Equal(t, actual, "<h1>Hello, Welcome to my HTTP server!</h1>")
}
