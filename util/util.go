/*
 Copyright © 2019-2024 Dell Inc. or its subsidiaries. All Rights Reserved.

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

package util

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Error messages
var (
	ErrorNameEmpty         = errors.New("name empty error")
	ErrorNameTooLong       = errors.New("name too long error")
	ErrorInvalidCharacters = errors.New("name contains invalid characters or name doesn't start with alphabetic. Allowed characters are 'a-zA-Z0-9_-'")
)

// UnityLog constant
const (
	UnityLog = "unitylog"
)

// UnityLogStruct is structure of UnityLog
type UnityLogStruct struct {
	unityLog string
}

// GetRunIDLogger function returns entry if exists
func GetRunIDLogger(ctx context.Context) *logrus.Entry {
	rlog := ctx.Value(UnityLog)
	entry := &logrus.Entry{}
	if rlog != nil && reflect.TypeOf(rlog) == reflect.TypeOf(entry) {
		entry = rlog.(*logrus.Entry)
	}
	if len(entry.Data) > 0 {
		return entry
	}

	log := GetLogger()
	return log.WithContext(ctx)
}

var (
	singletonLog *logrus.Logger
	once         sync.Once
)

// GetLogger is a singleton method which returns log object.
// Type singletonLog initialized only once.
func GetLogger() *logrus.Logger {
	once.Do(func() {
		singletonLog = logrus.New()
		fmt.Println("gounity logger initiated. This should be called only once.")

		// Gounity users can make use of this environment variable to initialize log level. Default level will be Info
		logLevel := os.Getenv("X_CSI_LOG_LEVEL")

		ChangeLogLevel(logLevel)

		singletonLog.SetReportCaller(true)
		singletonLog.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := strings.Split(f.File, "dell/gounity")
				if len(filename) > 1 {
					return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("dell/gounity%s:%d", filename[1], f.Line)
				}
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
			},
		}
	})

	return singletonLog
}

// ChangeLogLevel method returns log level
func ChangeLogLevel(logLevel string) {
	if singletonLog == nil {
		GetLogger()
	}

	switch strings.ToLower(logLevel) {

	case "debug":
		singletonLog.Level = logrus.DebugLevel
		break

	case "warn", "warning":
		singletonLog.Level = logrus.WarnLevel
		break

	case "error":
		singletonLog.Level = logrus.ErrorLevel
		break

	case "info":
		// Default level will be Info
		fallthrough

	default:
		singletonLog.Level = logrus.InfoLevel
	}
}

// ValidateResourceName function validate the resource name
func ValidateResourceName(name string, maxLength int) (string, error) {
	name = strings.TrimSpace(name)
	re := regexp.MustCompile("^[A-Za-z][a-zA-Z0-9:_-]*$")

	if name == "" {
		return "", ErrorNameEmpty
	} else if len(name) > maxLength {
		return "", ErrorNameTooLong
	} else if !re.MatchString(name) {
		return "", ErrorInvalidCharacters
	}

	return name, nil
}

// ValidateDuration function validates duration
func ValidateDuration(durationStr string) (uint64, error) {
	if durationStr != "" {
		durationArr := strings.Split(durationStr, ":")
		if len(durationArr) != 4 {
			return 0, errors.New("enter duration in Days:Hours:Mins:Secs, Ex: 1:23:52:50")
		}
		days, err := strconv.Atoi(durationArr[0])
		if err != nil {
			return 0, errors.New("invalid days in duration in Days:Hours:Mins:Secs, Ex: 1:23:52:50")
		}
		hours, err := strconv.Atoi(durationArr[1])
		if err != nil {
			return 0, errors.New("invalid hours in duration in Days:Hours:Mins:Secs, Ex: 1:23:52:50")
		}
		minutes, err := strconv.Atoi(durationArr[2])
		if err != nil {
			return 0, errors.New("invalid minutes in duration in Days:Hours:Mins:Secs, Ex: 1:23:52:50")
		}
		secs, err := strconv.Atoi(durationArr[3])
		if err != nil {
			return 0, errors.New("invalid seconds in duration in Days:Hours:Mins:Secs, Ex: 1:23:52:50")
		}
		if days < 0 {
			return 0, errors.New("days should be >= 0")
		}
		if hours >= 24 || hours < 0 {
			return 0, errors.New("hours should be - to 23")
		}
		if minutes >= 60 || minutes < 0 || secs >= 60 || secs < 0 {
			return 0, errors.New("hours, minutes and seconds should be in between 0-60")
		}
		seconds := uint64(days*24*60*60 + hours*60*60 + minutes*60 + secs)
		return seconds, nil
	}

	return 0, nil
}

// GetSecuredCipherSuites returns a slice of secured cipher suites.
// It iterates over the tls.CipherSuites() and appends the ID of each cipher suite to the suites slice.
// The function returns the suites slice.
func GetSecuredCipherSuites() (suites []uint16) {
	securedSuite := tls.CipherSuites()
	for _, v := range securedSuite {
		suites = append(suites, v.ID)
	}
	return suites
}
