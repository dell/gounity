package util

import (
	"errors"
	"github.com/dell/gounity/api"
	"regexp"
	"strconv"
	"strings"
)

var (
	NameEmptyError    = errors.New("name empty error")
	NameTooLongError  = errors.New("name too long error")
	InvalidCharacters = errors.New("name contains invalid characters or name doesn't start with alphabetic. Allowed characters are 'a-zA-Z0-9_-'")
)

//To validate the resource name
func ValidateResourceName(name string) (string, error) {
	name = strings.TrimSpace(name)
	re := regexp.MustCompile("^[A-Za-z][a-zA-Z0-9:_-]*$")

	if name == "" {
		return "", NameEmptyError
	} else if len(name) > api.MaxResourceNameLength {
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
