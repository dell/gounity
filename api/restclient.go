package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	HeaderKeyAccept                       = "Accept"
	HeaderKeyContentType                  = "Content-Type"
	HeaderValContentTypeJSON              = "application/json"
	headerValContentTypeBinaryOctetStream = "binary/octet-stream"
	HeaderEMCCSRFToken                    = "EMC-CSRF-TOKEN"
)

var (
	errNewClient = errors.New("missing endpoint")
	errSysCerts  = errors.New("unable to initialize certificate pool from system")
)

// Client Interface defines the methods.
type Client interface {
	// Get sends an HTTP request using the GET method to the API.
	Get(
		ctx context.Context,
		path string,
		headers map[string]string,
		resp interface{}) error

	// Post sends an HTTP request using the POST method to the API.
	Post(
		ctx context.Context,
		path string,
		headers map[string]string,
		body, resp interface{}) error

	// Put sends an HTTP request using the PUT method to the API.
	Put(
		ctx context.Context,
		path string,
		headers map[string]string,
		body, resp interface{}) error

	// Delete sends an HTTP request using the DELETE method to the API.
	Delete(
		ctx context.Context,
		path string,
		headers map[string]string,
		resp interface{}) error

	// ParseJSONError parses the JSON in r into an error object
	ParseJSONError(ctx context.Context, r *http.Response) error

	// DoWithHeaders sends HTTP request using the given method & headers.
	DoWithHeaders(
		ctx context.Context,
		method, uri string,
		headers map[string]string,
		body, resp interface{}) error

	// DoandGetResponseBody sends an HTTP request to the API and returns
	// the raw response body
	DoAndGetResponseBody(
		ctx context.Context,
		method, path string,
		headers map[string]string,
		body interface{}) (*http.Response, error)

	// SetToken sets the Auth token for the HTTP client
	SetToken(token string)

	// GetToken gets the Auth token for the HTTP client
	GetToken() string
}

type client struct {
	http     *http.Client
	host     string
	token    string
	showHTTP bool
	debug    bool
}

// ClientOptions are options for the API client.
type ClientOptions struct {
	// Insecure is a flag that indicates whether or not to supress SSL errors.
	Insecure bool

	// UseCerts is a flag that indicates whether system certs should be loaded
	UseCerts bool

	// Timeout specifies a time limit for requests made by this client.
	Timeout time.Duration

	// ShowHTTP is a flag that indicates whether or not HTTP requests and
	// responses should be logged to stdout
	ShowHTTP bool
}

//New returns a new API client.
func New(ctx context.Context, host string, opts ClientOptions, debug bool) (Client, error) {
	if host == "" {
		return nil, errNewClient
	}

	host = strings.Replace(host, "/api", "", 1)

	cookieJar, _ := cookiejar.New(nil)

	c := &client{
		http: &http.Client{},
		host: host,
	}

	if opts.Timeout != 0 {
		c.http.Timeout = opts.Timeout
	}

	if opts.Insecure {
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if opts.UseCerts {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, errSysCerts
		}
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: opts.Insecure,
			},
		}
	}
	c.http.Jar = cookieJar
	if opts.ShowHTTP {
		c.showHTTP = true
	}

	c.debug = debug

	return c, nil
}

// Makes a GET call to the Unity REST API Server with the given path & headers
func (c *client) Get(ctx context.Context, path string, headers map[string]string, resp interface{}) error {
	return c.DoWithHeaders(ctx, http.MethodGet, path, headers, nil, resp)
}

// Makes a POST call to the Unity REST API Server with the given path & headers
func (c *client) Post(ctx context.Context, path string, headers map[string]string, body, resp interface{}) error {
	return c.DoWithHeaders(ctx, http.MethodPost, path, headers, body, resp)
}

// Makes a PUT call to the Unity REST API Server with the given path & headers
func (c *client) Put(ctx context.Context, path string, headers map[string]string, body, resp interface{}) error {
	return c.DoWithHeaders(ctx, http.MethodPut, path, headers, body, resp)
}

// Makes a DELETE call to the Unity REST API Server with the given path & headers
func (c *client) Delete(ctx context.Context, path string, headers map[string]string, resp interface{}) error {
	return c.DoWithHeaders(ctx, http.MethodDelete, path, headers, nil, resp)
}

// Makes a PUT call to the Unity REST API Server with the given path & headers
func (c *client) Do(ctx context.Context, method, path string, body, resp interface{}) error {
	return c.DoWithHeaders(ctx, method, path, nil, body, resp)
}

func beginsWithSlash(s string) bool {
	return s[0] == '/'
}

func endsWithSlash(s string) bool {
	return s[len(s)-1] == '/'
}

func (c *client) DoAndGetResponseBody(ctx context.Context, method, uri string, headers map[string]string, body interface{}) (*http.Response, error) {
	log := util.GetRunIdLogger(ctx)
	var (
		err                error
		req                *http.Request
		res                *http.Response
		ubf                = &bytes.Buffer{}
		luri               = len(uri)
		hostEndsWithSlash  = endsWithSlash(c.host)
		uriBeginsWithSlash = beginsWithSlash(uri)
	)

	ubf.WriteString(c.host)

	if !hostEndsWithSlash && (luri > 0) {
		ubf.WriteString("/")
	}

	if luri > 0 {
		if uriBeginsWithSlash {
			ubf.WriteString(uri[1:])
		} else {
			ubf.WriteString(uri)
		}
	}

	u, err := url.Parse(ubf.String())
	if err != nil {
		return nil, err
	}

	var isContentTypeSet bool
	// marshal the message body (assumes json format)
	if r, ok := body.(io.ReadCloser); ok {
		req, err = http.NewRequest(method, u.String(), r)
		if err != nil {
			log.Errorf("Error while making new Request: %v", err)
			return nil, err
		}
		defer r.Close()
		if v, ok := headers[HeaderKeyContentType]; ok {
			req.Header.Set(HeaderKeyContentType, v)
		} else {
			req.Header.Set(HeaderKeyContentType, headerValContentTypeBinaryOctetStream)
		}
		isContentTypeSet = true
	} else if body != nil {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		if err = enc.Encode(body); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			log.Errorf("Error while making new Request: %v", err)
			return nil, err
		}
		if v, ok := headers[HeaderKeyContentType]; ok {
			req.Header.Set(HeaderKeyContentType, v)
		} else {
			req.Header.Set(HeaderKeyContentType, HeaderValContentTypeJSON)
		}
		isContentTypeSet = true
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			log.Errorf("Error while making new Request: %v", err)
			return nil, err
		}
	}

	if !isContentTypeSet {
		isContentTypeSet = req.Header.Get(HeaderKeyContentType) != ""
	}

	// add headers to the request
	for header, value := range headers {
		if header == HeaderKeyContentType && isContentTypeSet {
			continue
		}
		req.Header.Add(header, value)
	}

	// set the auth token for POST and DELETE methods only
	if (method == "POST" || method == "DELETE") && c.token != "" {
		req.Header.Set(HeaderEMCCSRFToken, c.token)
	}

	if c.showHTTP {
		logRequest(ctx, req, c.doLog)
	}

	// send the request
	req = req.WithContext(ctx)
	if res, err = c.http.Do(req); err != nil {
		return nil, err
	}

	if c.showHTTP {
		logResponse(ctx, res, c.doLog)
	}

	log.Infof("Response code:%d for url: %s", res.StatusCode, uri)
	return res, err
}

func (c *client) DoWithHeaders(ctx context.Context, method, uri string, headers map[string]string, body, resp interface{}) error {
	log := util.GetRunIdLogger(ctx)
	if body != nil {
		data, _ := json.Marshal(body)
		strBody := strings.ReplaceAll(string(data), "\"", "")
		log.Debugf("Request Body: %s", strBody)
	}
	res, err := c.DoAndGetResponseBody(ctx, method, uri, headers, body)
	if err != nil {
		log.Errorf("Error while receiving response for url: %s error: %v", uri, err)
		return err
	}
	defer res.Body.Close()

	// parse the response
	switch {
	case res == nil:
		log.Errorf("Nil Response received for url: %s", uri)
		return errors.New("nil Response received")
	case res.StatusCode >= 200 && res.StatusCode <= 299:
		dec := json.NewDecoder(res.Body)
		if resp != nil {
			if err = dec.Decode(resp); err != nil && err != io.EOF {
				log.Errorf("Decode error %v", err)
				c.doLog(log.WithError(err).Error, fmt.Sprintf("Unable to decode response into %+v", resp))
				return err
			}
		}
	case res.StatusCode == 401:
		jsonError := &types.Error{}
		if err := json.NewDecoder(res.Body).Decode(jsonError); err != nil {
			jsonError.ErrorContent.HTTPStatusCode = res.StatusCode
			jsonError.ErrorContent.Message = append(jsonError.ErrorContent.Message, types.ErrorMessage{http.StatusText(res.StatusCode)})
			return jsonError
		}

		jsonError.ErrorContent.HTTPStatusCode = res.StatusCode
		jsonError.ErrorContent.Message = append(jsonError.ErrorContent.Message, types.ErrorMessage{string(res.Status)})
		return jsonError
	default:
		c.doLog(log.WithError(err).Error, fmt.Sprintf("Invalid Response received Body: %s", body))
		return c.ParseJSONError(ctx, res)
	}
	return nil
}

func (c *client) SetToken(token string) {
	c.token = token
}

func (c *client) GetToken() string {
	return c.token
}

func (c *client) ParseJSONError(ctx context.Context, r *http.Response) error {
	log := util.GetRunIdLogger(ctx)
	jsonError := &types.Error{}
	err := json.NewDecoder(r.Body).Decode(jsonError)
	if err != nil && err != io.EOF {
		log.Error("ParseJSONError:", err)
		return err
	}
	data, err := json.Marshal(jsonError)
	if err != nil {
		log.Error("ParseJSONError marshal error", err)
	}
	log.WithError(err).Errorf("Json error response: %s", strings.ReplaceAll(string(data), "\"", ""))

	return jsonError
}

func (c *client) doLog(
	l func(args ...interface{}),
	msg string) {

	if c.debug {
		l(msg)
	}
}
