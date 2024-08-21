package manyrowsclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

const apiVersion = "v1"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	acceptGzip bool
}

func NewClient(baseURL, apiKey string, options ...func(*Client)) *Client {
	baseURL = strings.TrimSpace(baseURL)
	client := &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
	for _, option := range options {
		option(client)
	}
	if client.httpClient == nil {
		client.httpClient = defaultClient()
	}
	return client
}

func WithHTTPClient(httpClient *http.Client) func(client *Client) {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func WithAcceptGzip() func(client *Client) {
	return func(client *Client) {
		client.acceptGzip = true
	}
}

func defaultClient() *http.Client {
	var transport = &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	}
	var client = &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
	return client
}

func (client *Client) GetHttpClient() *http.Client {
	return client.httpClient
}

func (client *Client) Query(projectPath, kindPath string, req QueryRequest) (*QueryResponse, error) {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/entities/%s/query", options.baseURL, apiVersion, projectPath, kindPath)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res := QueryResponse{}
	_, err, _ = client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client *Client) GetOne(projectPath, kindPath string, req GetOneRequest) (*Entity, error) {
	options := client.getOptions(req.clientOptions)
	en := Entity{}
	u := fmt.Sprintf("%s/%s/%s/entities/%s/%s", options.baseURL, apiVersion, projectPath, kindPath, req.ID.String())
	_, err, _ := client.doRequest("GET", u, nil, options.apiKey, &en)
	if err != nil {
		return nil, err
	}
	return &en, nil
}

func (client *Client) Create(projectPath, kindPath string, req CreateEntityRequest) (*CreateEntityResponse, error, *ErrorInfo) {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/entities/%s", options.baseURL, apiVersion, projectPath, kindPath)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err, nil
	}
	resp, err, errInfo := client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, nil)
	if err != nil {
		return nil, err, errInfo
	}
	loc := resp.Header.Get("Location")
	entityID := resp.Header.Get("EntityID")
	id, err := uuid.FromString(entityID)
	if err != nil {
		return nil, err, nil
	}
	c := CreateEntityResponse{ID: id, Location: loc}
	return &c, nil, nil
}

func (client *Client) DeleteOne(projectPath, kindPath string, req DeleteOneRequest) error {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/entities/%s/%s", options.baseURL, apiVersion, projectPath, kindPath, req.ID.String())
	_, err, _ := client.doRequest("DELETE", u, nil, options.apiKey, nil)
	return err
}

func (client *Client) Update(projectPath, kindPath string, entityID uuid.UUID, req UpdateRequest) (error, *ErrorInfo) {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/entities/%s/%s", options.baseURL, apiVersion, projectPath, kindPath, entityID)
	payload, err := json.Marshal(req)
	if err != nil {
		return err, nil
	}
	_, err, errorInfo := client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, nil)
	return err, errorInfo
}

func (client *Client) DeleteEntities(projectPath, kindPath string, req DeleteRequest) error {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/entities/%s/delete", options.baseURL, apiVersion, projectPath, kindPath)
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err, _ = client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, nil)
	return err
}

func (client *Client) QueryCollection(projectPath string, collectionId uuid.UUID, req QueryRequest) (*QueryResponse, error) {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/collections/%s/query", options.baseURL, apiVersion, projectPath, collectionId)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res := QueryResponse{}
	_, err, _ = client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (client *Client) CreateCollectionItem(projectPath string, collectionId uuid.UUID, req CreateCollectionItemRequest) error {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/collections/%s", options.baseURL, apiVersion, projectPath, collectionId)
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err, _ = client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, nil)
	return err
}

func (client *Client) MoveCollectionItem(projectPath string, collectionId uuid.UUID, req MoveCollectionItemRequest) error {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/collections/%s/move", options.baseURL, apiVersion, projectPath, collectionId)
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err, _ = client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, nil)
	return err
}

func (client *Client) DeleteCollectionItems(projectPath string, collectionId uuid.UUID, req DeleteCollectionItemsRequest) error {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/%s/%s/collections/%s/remove", options.baseURL, apiVersion, projectPath, collectionId)
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err, _ = client.doRequest("POST", u, bytes.NewReader(payload), options.apiKey, nil)
	return err
}

func (client *Client) DeleteProject(req DeleteOneRequest) error {
	options := client.getOptions(req.clientOptions)
	u := fmt.Sprintf("%s/s/projects/%s", options.baseURL, req.ID)
	_, err, _ := client.doRequest("DELETE", u, nil, options.apiKey, nil)
	return err
}

func (client *Client) getOptions(opts clientOptions) clientOptions {
	baseURL := client.baseURL
	apiKey := client.apiKey
	if opts.baseURL != "" {
		baseURL = opts.baseURL
	}
	if opts.apiKey != "" {
		apiKey = opts.apiKey
	}
	return clientOptions{baseURL: baseURL, apiKey: apiKey}
}

func (client *Client) doRequest(method, url string, body io.Reader, apiKey string, respObj interface{}) (*http.Response, error, *ErrorInfo) {
	rq, err := http.NewRequest(method, url, body)
	errInfo := ErrorInfo{}
	if err != nil {
		return nil, err, nil
	}
	rq.Header.Set("ManyRowsAuthToken", apiKey)
	rq.Header.Set("Content-Type", "application/json")
	if client.acceptGzip {
		rq.Header.Set("Accept-Encoding", "gzip, deflate")
	}
	resp, err := client.httpClient.Do(rq)
	if err != nil {
		return nil, err, &errInfo
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		err = extractJSON(client.acceptGzip, resp, &errInfo)
		if err != nil {
			return resp, err, nil
		}
		errInfo.HttpCode = resp.StatusCode
		return resp, fmt.Errorf("status code was %d", resp.StatusCode), &errInfo
	} else if resp.StatusCode >= 500 || resp.StatusCode < 200 {
		errInfo.HttpCode = resp.StatusCode
		return resp, fmt.Errorf("unexpected server error: status code was %d", resp.StatusCode), &errInfo
	}
	if respObj != nil {
		err = extractJSON(client.acceptGzip, resp, respObj)
		if err != nil {
			return resp, err, nil
		}
	}
	return resp, nil, nil
}

func extractJSON(gzipped bool, resp *http.Response, respObj interface{}) error {
	var reader io.Reader
	var err error
	if gzipped {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
	} else {
		reader = resp.Body
	}
	dec := json.NewDecoder(reader)
	err = dec.Decode(&respObj)
	return err
}
