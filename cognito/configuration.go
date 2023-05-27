package cognito

import (
	"errors"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
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

func GetUsername() string {
	usernamePrompt := promptui.Prompt{
		Label: "Username",
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		color.Red("Error: %v", err.Error())
		os.Exit(1)
	}
	return username
}

func GetPassword() string {
	passwordPrompt := promptui.Prompt{
		Label:    "Password",
		Mask:     '*',
		Validate: CheckIfValidPassword,
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		color.Red("Error: %v", err.Error())
		os.Exit(1)
	}
	return password
}
