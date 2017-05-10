package telesign

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sauliuspr/go-telesign/auth"
)

const (
	libraryVersion = "1"
	defaultBaseURL = "https://rest-ww.telesign.com"
	userAgent      = "TelesigSDK/GO v" + libraryVersion
)

type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
	common    service // Reuse of one struct

	// Telesign Services for talking to different parts of REST API
	Ping         *PingService
	SMSVerify    *SMSVerifyService
	GetStatus    *GetStatusService
	PhoneIDScore *PhoneIDScoreService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Ping = (*PingService)(&c.common)
	c.SMSVerify = (*SMSVerifyService)(&c.common)
	c.GetStatus = (*GetStatusService)(&c.common)
	c.PhoneIDScore = (*PhoneIDScoreService)(&c.common)

	return c
}

func (c *Client) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.BaseURL.ResolveReference(rel).String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

type ErrorSource struct {
	Parameter string `json:"parameter"`
}

type Error struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	Response *http.Response
	Errors   []Error
}

type responseWrapper struct {
	*http.Response
	Data interface{} `json:"data"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("ERROR: %v %v: %d %+v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Errors)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		defer r.Body.Close()
		// Ignore errors here so we always pass out an ErrorResponse
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func (c *Client) Do(ctx context.Context, req *http.Request, resource interface{}) (*http.Response, error) {
	ctx, req = withContext(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	if err := CheckResponse(resp); err != nil {
		return resp, err
	}
	// If Request Body exist preserve it
	var zdata []byte
	if resp.Body != nil {
		zdata, _ = ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(zdata))
	}

	if err := json.NewDecoder(resp.Body).Decode(&resource); err != nil {
		return resp, err
	}

	return resp, nil
}

func withContext(ctx context.Context, req *http.Request) (context.Context, *http.Request) {
	return ctx, req.WithContext(ctx)
}

type APIAuthTransport struct {
	CustomerID string
	APIKey     string // API Secret Token

	Transport http.RoundTripper
}

func (t *APIAuthTransport) Client() *http.Client {
	return &http.Client{
		Transport: t,
	}
}

func (t *APIAuthTransport) getTransport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *APIAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.getTransport().RoundTrip(auth.GenerateTelesignHeaders(req, t.CustomerID, t.APIKey))
}
