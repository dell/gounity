/*
 Copyright © 2019-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
)

// Header Key constants
const (
	HeaderKeyAccept                       = "Accept"
	HeaderKeyContentType                  = "Content-Type"
	HeaderValContentTypeJSON              = "application/json"
	headerValContentTypeBinaryOctetStream = "binary/octet-stream"
	HeaderEMCCSRFToken                    = "EMC-CSRF-TOKEN" // #nosec G101
)

var (
	errNewClient       = errors.New("missing endpoint")
	errSysCerts        = errors.New("unable to initialize certificate pool from system")
	systemCertPoolFunc = x509.SystemCertPool
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

	// Timeout specifies a time limit for requests made by this client.
	Timeout time.Duration

	// ShowHTTP is a flag that indicates whether or not HTTP requests and
	// responses should be logged to stdout
	ShowHTTP bool
}

// New returns a new API client.
func New(_ context.Context, host string, opts ClientOptions, debug bool) (Client, error) {
	if host == "" {
		return nil, errNewClient
	}

	host = strings.Replace(host, "/api", "", 1)

	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: nil})

	c := &client{
		http:  &http.Client{},
		host:  host,
		debug: debug,
	}

	if opts.Timeout != 0 {
		c.http.Timeout = opts.Timeout
	}

	if opts.Insecure { // #nosec G402
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				CipherSuites:       util.GetSecuredCipherSuites(),
			},
		}
	} else {
		pool, err := systemCertPoolFunc()
		if err != nil {
			return nil, errSysCerts
		}
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: false,
				CipherSuites:       util.GetSecuredCipherSuites(),
				MinVersion:         tls.VersionTLS12,
			},
		}
	}
	c.http.Jar = cookieJar
	if opts.ShowHTTP {
		c.showHTTP = true
	}
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
	log := util.GetRunIDLogger(ctx)
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

	log.Debugf("Response code:%d for url: %s", res.StatusCode, uri)
	return res, err
}

func (c *client) DoWithHeaders(ctx context.Context, method, uri string, headers map[string]string, body, resp interface{}) error {
	log := util.GetRunIDLogger(ctx)
	if body != nil {
		data, _ := json.Marshal(body)
		strBody := strings.ReplaceAll(string(data), "\"", "")
		log.Debugf("Request Body: %s", strBody)
	}
	res, err := c.DoAndGetResponseBody(ctx, method, uri, headers, body)
	if err != nil {
		return fmt.Errorf("Error while receiving response for url: %s error: %v", uri, err)
	}
	defer res.Body.Close()

	// parse the response
	switch {
	case res == nil:
		return fmt.Errorf("Nil Response received for url: %s", uri)
	case res.StatusCode >= 200 && res.StatusCode <= 299:
		dec := json.NewDecoder(res.Body)
		if resp != nil {
			if err = dec.Decode(resp); err != nil && err != io.EOF {
				c.doLog(log.WithError(err).Error, fmt.Sprintf("Unable to decode response into %+v", resp))
				return err
			}
		}
	case res.StatusCode == 401:
		jsonError := &types.Error{}
		if err := json.NewDecoder(res.Body).Decode(jsonError); err != nil {
			jsonError.ErrorContent.HTTPStatusCode = res.StatusCode
			jsonError.ErrorContent.Message = append(jsonError.ErrorContent.Message, types.ErrorMessage{EnUS: http.StatusText(res.StatusCode)})
			return jsonError
		}

		jsonError.ErrorContent.HTTPStatusCode = res.StatusCode
		jsonError.ErrorContent.Message = append(jsonError.ErrorContent.Message, types.ErrorMessage{EnUS: string(res.Status)})
		return jsonError
	default:
		log.Debugf("Invalid Response received Body: %s error: %v", body, err)
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
	log := util.GetRunIDLogger(ctx)
	jsonError := &types.Error{}
	err := json.NewDecoder(r.Body).Decode(jsonError)
	if err != nil && err != io.EOF {
		return fmt.Errorf("ParseJSONError: %v", err)
	}
	_, err = json.Marshal(jsonError)
	if err != nil {
		log.Error("ParseJSONError marshal error", err)
	}

	return jsonError
}

func (c *client) doLog(l func(args ...interface{}), msg string) {
	if c.debug {
		l(msg)
	}
}
