/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package api

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func isBinOctetBody(h http.Header) bool {
	return h.Get(HeaderKeyContentType) == headerValContentTypeBinaryOctetStream
}

var singletonLog *logrus.Logger
var once sync.Once

//This is a singleton method which returns log object.
//Type singletonLog initialized only once.
func GetLogger() *logrus.Logger {
	once.Do(func() {
		singletonLog = logrus.New()
		fmt.Println("gounity logger initiated. This should be called only once.")
		var debug bool
		debugStr := os.Getenv("GOUNITY_DEBUG")
		debug, _ = strconv.ParseBool(debugStr)
		if debug {
			fmt.Println("Enabling debug for gounity")
			singletonLog.Level = logrus.DebugLevel
			singletonLog.SetReportCaller(true)
			singletonLog.Formatter = &logrus.TextFormatter{
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					filename := strings.Split(f.File, "dell/gounity")
					if len(filename) > 1 {
						return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("dell/gounity%s:%d", filename[1], f.Line)
					} else {
						return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
					}
				},
			}
		}
	})

	return singletonLog
}

func logRequest(
	ctx context.Context,
	req *http.Request,
	lf func(func(args ...interface{}), string)) {

	w := &bytes.Buffer{}

	fmt.Fprintln(w)
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOUNITY HTTP REQUEST")
	fmt.Fprintln(w, " -------------------------")

	buf, err := httputil.DumpRequest(req, !isBinOctetBody(req.Header))
	if err != nil {
		return
	}

	WriteIndented(w, buf)
	fmt.Fprintln(w)

	lf(log.Debug, w.String())
}

func logResponse(
	ctx context.Context,
	res *http.Response,
	lf func(func(args ...interface{}), string)) {

	w := &bytes.Buffer{}

	fmt.Fprintln(w)
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOUNITY HTTP RESPONSE")
	fmt.Fprintln(w, " -------------------------")

	buf, err := httputil.DumpResponse(res, !isBinOctetBody(res.Header))
	if err != nil {
		return
	}

	bw := &bytes.Buffer{}
	WriteIndented(bw, buf)

	scanner := bufio.NewScanner(bw)
	for {
		if !scanner.Scan() {
			break
		}
		fmt.Fprintln(w, scanner.Text())
	}

	log.Debug(w.String())
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
