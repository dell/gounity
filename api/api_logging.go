// Copyright © 2019-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

func isBinOctetBody(h http.Header) bool {
	return h.Get(HeaderKeyContentType) == headerValContentTypeBinaryOctetStream
}

type (
	dumpRequestFunc   func(req *http.Request, body bool) ([]byte, error)
	dumpResponseFunc  func(req *http.Response, body bool) ([]byte, error)
	writeIndentedFunc func(w io.Writer, b []byte) error
)

var (
	dumpRequest   dumpRequestFunc   = httputil.DumpRequest
	dumpResponse  dumpResponseFunc  = httputil.DumpResponse
	writeIndented writeIndentedFunc = WriteIndented
)

func logRequest(
	_ context.Context,
	req *http.Request,
	lf func(func(args ...interface{}), string),
) {
	w := &bytes.Buffer{}

	fmt.Fprintln(w)
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOUNITY HTTP REQUEST")
	fmt.Fprintln(w, " -------------------------")

	buf, err := dumpRequest(req, !isBinOctetBody(req.Header))
	if err != nil {
		return
	}

	err2 := writeIndented(w, buf)
	if err2 != nil {
		message := fmt.Sprintf("Indentation failed with error: %v", err2)
		log.Info(message)
	}
	fmt.Fprintln(w)

	lf(log.Debug, w.String())
}

func logResponse(
	_ context.Context,
	res *http.Response,
	lf func(func(args ...interface{}), string),
) {
	w := &bytes.Buffer{}

	fmt.Fprintln(w)
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOUNITY HTTP RESPONSE")
	fmt.Fprintln(w, " -------------------------")

	buf, err := dumpResponse(res, !isBinOctetBody(res.Header))
	if err != nil {
		return
	}

	bw := &bytes.Buffer{}
	err2 := writeIndented(bw, buf)
	if err2 != nil {
		message := fmt.Sprintf("Indentation failed with error: %v", err2)
		log.Info(message)
	}

	scanner := bufio.NewScanner(bw)
	for {
		if !scanner.Scan() {
			break
		}
		fmt.Fprintln(w, scanner.Text())
	}

	lf(log.Debug, w.String())
}

// WriteIndentedN indents all lines n spaces.
func WriteIndentedN(w io.Writer, b []byte, n int) error {
	s := bufio.NewScanner(bytes.NewReader(b))
	if !s.Scan() {
		return nil
	}
	l := s.Text()
	for {
		for x := 0; x < n; x++ {
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, l); err != nil {
			return err
		}
		if !s.Scan() {
			break
		}
		l = s.Text()
		if _, err := fmt.Fprint(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

// WriteIndented indents all lines four spaces.
func WriteIndented(w io.Writer, b []byte) error {
	return WriteIndentedN(w, b, 4)
}
