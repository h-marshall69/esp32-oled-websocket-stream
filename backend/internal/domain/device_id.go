package domain

import (
	"errors"
	"regexp"
)

var deviceIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)

var ErrInvalidDeviceID = errors.New("invalid device id")

func ValidateDeviceID(id string) error {
	if !deviceIDPattern.MatchString(id) {
		return ErrInvalidDeviceID
	}

	return nil
}
