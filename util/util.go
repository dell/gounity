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
	RUNIDLOG = "runidlog"
	RUNID    = "runid"
)

func GetRunIdLogger(ctx context.Context) *logrus.Entry {
	temp := ctx.Value(RUNIDLOG)
	entry := &logrus.Entry{}
	if reflect.TypeOf(temp) == reflect.TypeOf(entry) {
		return ctx.Value(RUNIDLOG).(*logrus.Entry)
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
