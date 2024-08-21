package manyrowsclient

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestNewClient_BlankParams(t *testing.T) {
	client := NewClient("", "")
	require.NotNil(t, client)
	require.Equal(t, "", client.apiKey)
	require.Equal(t, "", client.baseURL)
	require.Equal(t, false, client.acceptGzip)
}

func TestNewClient_WithParams(t *testing.T) {
	client := NewClient(" https://someurl.com ", "abcd")
	require.NotNil(t, client)
	require.Equal(t, "abcd", client.apiKey)
	require.Equal(t, "https://someurl.com", client.baseURL)
	require.Equal(t, false, client.acceptGzip)
}

func TestNewClient_DefaultClient(t *testing.T) {
	client := NewClient("", "")
	require.Equal(t, client.httpClient, defaultClient())
}

func TestClient_GetHttpClient(t *testing.T) {
	client := NewClient("", "")
	require.Equal(t, client.httpClient, defaultClient())
}

func TestNewClient_WithHttpClient(t *testing.T) {
	someClient := &http.Client{
		Timeout: time.Second * 1,
	}
	client := NewClient("", "", WithHTTPClient(someClient))
	require.Equal(t, client.httpClient, someClient)
}

func TestNewClient_WithAcceptGzip(t *testing.T) {
	client := NewClient("", "", WithAcceptGzip())
	require.Equal(t, client.acceptGzip, true)
}

func TestClientOptionsOverride(t *testing.T) {
	client := NewClient("bbb", "aaa")
	overrideOptions := clientOptions{apiKey: "a", baseURL: "b"}
	options := client.getOptions(overrideOptions)
	require.Equal(t, "a", options.apiKey)
	require.Equal(t, "b", options.baseURL)
}

func TestClientOptionsFallback(t *testing.T) {
	client := NewClient("bbb", "aaa")
	overrideOptions := clientOptions{apiKey: "", baseURL: ""}
	options := client.getOptions(overrideOptions)
	require.Equal(t, "aaa", options.apiKey)
	require.Equal(t, "bbb", options.baseURL)
}
