package cognito

import (
	"errors"
	"regexp"
)

func CheckIfValidPassword(password string) error {
	// Password must be 8 characters or greater
	if len(password) < 8 {
		return errors.New("password must be 8 characters or longer")
	}

	// Password must contain at least one number
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	// Password must contain at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]{1,}`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Password must contain at least one lowercase letter
	if !regexp.MustCompile(`[a-z]{1,}`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Password must contain at least one special character
	if !regexp.MustCompile(`[$*.[\]{}()?\-\"!@#%&/\\,><':;|_~+=]{1,}`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
