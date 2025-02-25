// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dell/gounity/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EmptyMockBody struct{}

// MockBody is a mock implementation of the io.ReadCloser interface
type MockBody struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m *MockBody) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

func (m *MockBody) Close() error {
	return m.CloseFunc()
}

func TestDoAndGetResponseBody(t *testing.T) {
	// Create a mock client
	c := &client{}
	c.SetToken("token")
	token := c.GetToken()
	c = &client{
		host:     "https://example.com",
		http:     http.DefaultClient,
		showHTTP: true,
		token:    token,
	}
	ctx := context.Background()
	// Create a mock request body
	body := &MockBody{
		ReadFunc: func(_ []byte) (n int, err error) {
			return 0, io.EOF
		},
		CloseFunc: func() error {
			return nil
		},
	}
	_, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com/api/v1/endpoint", nil)
	require.NoError(t, err)

	// Create a mock server for get method
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	c.host = server.URL
	res, err := c.DoAndGetResponseBody(ctx, http.MethodGet, "api/v1/endpoint", nil, body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	//  Create a mock server for delete method
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	c.host = server.URL
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "some_token_value",
		"User-Agent":    "Go-Client/1.0",
		"Accept":        "application/json",
	}
	res, err = c.DoAndGetResponseBody(ctx, http.MethodDelete, "api/v1/endpoint", headers, body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// with empty body
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	c.host = server.URL
	body1 := EmptyMockBody{}
	res, err = c.DoAndGetResponseBody(ctx, http.MethodGet, "api/v1/endpoint", headers, body1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	headers = map[string]string{
		"Authorization": "some_token_value",
		"Accept":        "application/json",
	}
	res, err = c.DoAndGetResponseBody(ctx, http.MethodGet, "api/v1/endpoint", headers, body1)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = c.DoAndGetResponseBody(ctx, http.MethodGet, "/api/v1/endpoint", nil, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDoWithHeaders(t *testing.T) {
	// Create a mock client
	c := &client{
		host: "https://example.com",
		http: http.DefaultClient,
	}

	ctx := context.Background()
	// Create a mock request body
	body := &MockBody{
		ReadFunc: func(_ []byte) (n int, err error) {
			return 0, io.EOF
		},
		CloseFunc: func() error {
			return nil
		},
	}

	// Create a mock request
	_, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com/api/v1/endpoint", nil)
	require.NoError(t, err)

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint", r.URL.String())

		// Write the mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	}))
	defer server.Close()

	// Set the client's host to the mock server's URL
	c.host = server.URL

	// Create a mock response object
	var responseData map[string]string

	// Call the function
	err = c.DoWithHeaders(ctx, http.MethodGet, "api/v1/endpoint", nil, body, &responseData)
	// Assert that there was no error
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"message": "ok"}, responseData)

	// for 401 response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	c.host = server.URL
	err = c.DoWithHeaders(ctx, http.MethodGet, "api/v1/endpoint", nil, body, &responseData)
	errorContent := types.ErrorContent{
		Message: []types.ErrorMessage{
			{
				EnUS: "Unauthorized",
			},
		},
		HTTPStatusCode: 401,
		ErrorCode:      0,
	}
	expectedError := types.Error{
		ErrorContent: errorContent,
	}
	assert.Equal(t, &expectedError, err)

	// for 400 response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	c.host = server.URL
	err = c.DoWithHeaders(ctx, http.MethodGet, "api/v1/endpoint", nil, body, &responseData)
	errorContent = types.ErrorContent{
		Message:        nil,
		HTTPStatusCode: 0,
		ErrorCode:      0,
	}
	expectedError = types.Error{
		ErrorContent: errorContent,
	}
	assert.Equal(t, &expectedError, err)
}

func TestNew(t *testing.T) {
	ctx := context.Background()

	opts := ClientOptions{
		Insecure: false,
		Timeout:  0,
		ShowHTTP: false,
	}

	debug := false
	_, err := New(ctx, "", opts, debug)
	assert.Equal(t, errors.New("missing endpoint"), err)

	host := "https://example.com"
	_, err = New(ctx, host, opts, debug)
	assert.Equal(t, nil, err)

	opts = ClientOptions{
		Insecure: true,
		Timeout:  10,
		ShowHTTP: true,
	}
	_, err = New(ctx, host, opts, debug)
	assert.Equal(t, nil, err)
}

func TestDoLog(t *testing.T) {
	c := &client{
		debug: true,
	}
	// Create a mock logger
	var loggedMessage string
	logger := func(args ...interface{}) {
		loggedMessage = args[0].(string)
	}

	c.doLog(logger, "Test message")
	assert.Equal(t, "Test message", loggedMessage)
}

func TestClientGet(t *testing.T) {
	c := &client{
		http: http.DefaultClient,
		host: "https://example.com",
	}
	err := c.Get(context.Background(), c.host, nil, nil)
	assert.Error(t, err)
}

func TestClientPost(t *testing.T) {
	c := &client{
		http: http.DefaultClient,
		host: "https://example.com",
	}
	err := c.Post(context.Background(), c.host, nil, nil, nil)
	assert.Error(t, err)
}

func TestClientPut(t *testing.T) {
	c := &client{
		http: http.DefaultClient,
		host: "https://example.com",
	}
	err := c.Put(context.Background(), c.host, nil, nil, nil)
	assert.Error(t, err)
}

func TestClientDelete(t *testing.T) {
	c := &client{
		http: http.DefaultClient,
		host: "https://example.com",
	}
	err := c.Delete(context.Background(), c.host, nil, nil)
	assert.Error(t, err)
}

func TestClientDo(t *testing.T) {
	c := &client{
		http: http.DefaultClient,
		host: "https://example.com",
	}
	err := c.Do(context.Background(), http.MethodGet, c.host, nil, nil)
	assert.Error(t, err)
}
