package goarubacloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"strings"

	"github.com/hashicorp/logutils"
)

const (
	libraryVersion    = "0.4.0"
	apiServerEnvName  = "ARUBACLOUD_APISERVER"
	logLevelEnvName   = "ARUBACLOUD_LOG"
	apiServerBasePath = "/WsEndUser/v2.9/WsEndUser.svc/json"
	userAgent         = "goarubacloud/" + libraryVersion
	mediaType         = "application/json"
)

func init() {
	logLevel := os.Getenv(logLevelEnvName)
	if logLevel == "" {
		logLevel = "INFO"
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

// Client manages communication with Arubacloud API.
type Client struct {
	// HTTP client used to communicate with the Arubacloud API.
	client *http.Client

	// Data center region
	Datacenter DataCenterRegion

	// Base URL for API requests.
	BaseURL *url.URL

	// Control panel username
	Username string

	// Control panel password
	Password string

	// User agent for client
	UserAgent string

	// Services used for communicating with the API
	DataCenters        DataCentersService
	Hypervisors        HypervisorsService
	CloudServers       CloudServersService
	CloudServerActions CloudServerActionsService
	ScheduledTasks     ScheduledTasksService
	Snapshots          SnapshotsService
	PurchasedIPs       PurchasedIPsService
	VLANs              VLANsService

	// Optional function called after every successful request made to the Arubacloud API
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// Response is a Arubacloud API response. This wraps the standard http.Response returned from Arubacloud.
type Response struct {
	*http.Response
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	Success bool

	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"ResultMessage"`

	// ResultCode returned from the API
	ResultCode int `json:"ResultCode"`
}

// NewClient returns a new Arubacloud API client.
func NewClient(datacenter DataCenterRegion, username, password string) *Client {
	apiServerHost := os.Getenv(apiServerEnvName)
	if apiServerHost == "" {
		apiServerHost = fmt.Sprintf("https://api.dc%d.computing.cloud.it", datacenter)
	}

	apiServerBaseUrl := fmt.Sprintf("%s%s", apiServerHost, apiServerBasePath)
	httpClient := http.DefaultClient
	baseURL, _ := url.Parse(apiServerBaseUrl)

	log.Printf("[DEBUG] Base URL: %s\n", baseURL)

	client := &Client{client: httpClient,
		BaseURL:   baseURL,
		Username:  username,
		Password:  password,
		UserAgent: userAgent}

	client.DataCenters = &DataCentersServiceOp{client: client}
	client.Hypervisors = &HypervisorsServiceOp{client: client}
	client.CloudServers = &CloudServersServiceOp{client: client}
	client.CloudServerActions = &CloudServerActionsServiceOp{client: client}
	client.ScheduledTasks = &ScheduledTasksServiceOp{client: client}
	client.Snapshots = &SnapshotsServiceOp{client: client}
	client.PurchasedIPs = &PurchasedIPsServiceOp{client: client}
	client.VLANs = &VLANsServiceOp{client: client}

	return client
}

// NewRequest creates an API request
func (c *Client) NewRequest(action string, body interface{}) (*http.Request, error) {
	callUrl := fmt.Sprintf("%s/%s", c.BaseURL.String(), action)

	requestMap := map[string]interface{}{
		"ApplicationId": action,
		"RequestId":     action,
		"SessionId":     action,
		"Username":      c.Username,
		"Password":      c.Password,
	}

	var buffer []byte
	var err error
	if body != nil {
		buffer, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	var bodyMap map[string]interface{}

	if len(buffer) > 0 {
		err := json.Unmarshal(buffer, &bodyMap)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range bodyMap {
		requestMap[k] = v
	}

	buffer, err = json.Marshal(requestMap)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(buffer)

	log.Printf("[DEBUG] Request: %s\n", bodyBuffer.String())
	req, err := http.NewRequest("POST", callUrl, bodyBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// OnRequestCompleted sets the API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Response:%s\n", bytes.NewBuffer(data).String())

	err = CheckResponse(resp, data)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err := io.Copy(w, bytes.NewBuffer(data))
			if err != nil {
				return nil, err
			}
		} else {
			err := json.NewDecoder(bytes.NewBuffer(data)).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%s. Result code: %d", r.Message, r.ResultCode)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response, data []byte) error {
	errorResponse := &ErrorResponse{Response: r}
	if len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}
	if errorResponse.Success == false {
		message := errorResponse.Message
		trimPosition := strings.Index(message, "\r")
		errorResponse.Message = message[0:trimPosition]
		return errorResponse
	}

	return nil
}
