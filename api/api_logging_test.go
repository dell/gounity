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
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"
)

// Custom writer that fails after a certain number of writes
type errorWriter struct {
	maxWrites int
	writes    int
}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	if ew.writes >= ew.maxWrites {
		return 0, errors.New("write error")
	}
	ew.writes++
	return len(p), nil
}

func TestIsBinOctetBody(t *testing.T) {
	// Test case: header with correct content type
	header := http.Header{}
	header.Set(HeaderKeyContentType, headerValContentTypeBinaryOctetStream)
	if !isBinOctetBody(header) {
		t.Errorf("isBinOctetBody() = false, want true")
	}

	// Test case: header with incorrect content type
	header = http.Header{}
	header.Set(HeaderKeyContentType, "application/json")
	if isBinOctetBody(header) {
		t.Errorf("isBinOctetBody() = true, want false")
	}

	// Test case: header with empty content type
	header = http.Header{}
	if isBinOctetBody(header) {
		t.Errorf("isBinOctetBody() = true, want false")
	}
}

func TestLogRequest(t *testing.T) {
	// Test case: Normal request
	t.Run("Normal request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		logRequest(context.TODO(), req, nil)
	})

	// Test case: Request with binary octet stream
	t.Run("Request with binary octet stream", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(HeaderKeyContentType, headerValContentTypeBinaryOctetStream)

		logRequest(context.TODO(), req, nil)
	})

	// Test case: Request with body
	t.Run("Request with body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer([]byte("test")))
		if err != nil {
			t.Fatal(err)
		}

		logRequest(context.TODO(), req, nil)
	})

	// Test case: Request with body and binary octet stream
	t.Run("Request with body and binary octet stream", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://example.com", bytes.NewBuffer([]byte("test")))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(HeaderKeyContentType, headerValContentTypeBinaryOctetStream)

		logRequest(context.TODO(), req, nil)
	})
}

func TestLogResponse(t *testing.T) {
	// Test case: Response with valid headers
	t.Run("ValidHeaders", func(_ *testing.T) {
		res := &http.Response{
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		}
		logResponse(context.Background(), res, nil)
		// Add assertions to check if the response is logged correctly
	})

	// Test case: Response with empty headers
	t.Run("EmptyHeaders", func(_ *testing.T) {
		res := &http.Response{
			Header: http.Header{},
		}
		logResponse(context.Background(), res, nil)
		// Add assertions to check if the response is logged correctly
	})

	// Test case: Response with binary octet stream
	t.Run("BinOctetStream", func(_ *testing.T) {
		res := &http.Response{
			Header: http.Header{
				"Content-Type": []string{"binary/octet-stream"},
			},
		}
		logResponse(context.Background(), res, nil)
		// Add assertions to check if the response is logged correctly
	})

	// Test case: Response with indentation error
	t.Run("IndentationError", func(_ *testing.T) {
		res := &http.Response{
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		}
		logResponse(context.Background(), res, nil)
		// Add assertions to check if the indentation error is logged correctly
	})
}

func TestWriteIndentedN(t *testing.T) {
	// Test case: empty input
	var buf bytes.Buffer
	err := WriteIndentedN(&buf, []byte{}, 4)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if buf.String() != "" {
		t.Errorf("Expected empty output, got %q", buf.String())
	}

	// Test case: single line
	buf.Reset()
	err = WriteIndentedN(&buf, []byte("Hello, world!"), 4)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := "    Hello, world!"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}

	// Test case: multiple lines
	buf.Reset()
	err = WriteIndentedN(&buf, []byte("Hello\nworld!"), 4)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected = "    Hello\n    world!"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}

	// Error conditions using errorWriter
	// Test case: Error at initial indent
	ew := &errorWriter{maxWrites: 0}
	err = WriteIndentedN(ew, []byte("Hello, world!"), 4)
	if err == nil || err.Error() != "write error" {
		t.Errorf("Expected write error, got %v", err)
	}

	// Test case: Error in writing line content
	ew = &errorWriter{maxWrites: 4}
	err = WriteIndentedN(ew, []byte("Hello, world!"), 4)
	if err == nil || err.Error() != "write error" {
		t.Errorf("Expected write error, got %v", err)
	}

	// Test case: Error in writing newline
	ew = &errorWriter{maxWrites: 9} // Enough for "    Hello"
	err = WriteIndentedN(ew, []byte("Hello\nworld!"), 4)
	if err == nil || err.Error() != "write error" {
		t.Errorf("Expected write error, got %v", err)
	}
}
