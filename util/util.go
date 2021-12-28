package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	NameEmptyError    = errors.New("name empty error")
	NameTooLongError  = errors.New("name too long error")
	InvalidCharacters = errors.New("name contains invalid characters or name doesn't start with alphabetic. Allowed characters are 'a-zA-Z0-9_-'")
)

const (
	UnityLog = "unitylog"
)

func GetRunIdLogger(ctx context.Context) *logrus.Entry {
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

var singletonLog *logrus.Logger
var once sync.Once

//This is a singleton method which returns log object.
//Type singletonLog initialized only once.
func GetLogger() *logrus.Logger {
	once.Do(func() {
		singletonLog = logrus.New()
		fmt.Println("gounity logger initiated. This should be called only once.")

		//Gounity users can make use of this environment variable to initialize log level. Default level will be Info
		logLevel := os.Getenv("X_CSI_LOG_LEVEL")

		ChangeLogLevel(logLevel)

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
	})

	return singletonLog
}

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
		//Default level will be Info
		fallthrough

	default:
		singletonLog.Level = logrus.InfoLevel
	}
}

//To validate the resource name
func ValidateResourceName(name string, maxLength int) (string, error) {
	name = strings.TrimSpace(name)
	re := regexp.MustCompile("^[A-Za-z][a-zA-Z0-9:_-]*$")

	if name == "" {
		return "", NameEmptyError
	} else if len(name) > maxLength {
		return "", NameTooLongError
	} else if !re.MatchString(name) {
		return "", InvalidCharacters
	}

	return name, nil
}

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
