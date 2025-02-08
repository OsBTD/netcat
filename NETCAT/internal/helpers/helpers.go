package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// HandleArguments processes command-line arguments and returns the port to use.
func HandleArgs() (port string) {
	if len(os.Args) < 2 {
		port = "8989"
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		port = ""
	}
	return
}

// ValidName checks if the username is valid:
func ValidUserName(name string) error {
	if len(name) < 3 {
		return errors.New("name should contain at least have 3 characters")
	}
	if len(name) > 13 {
		return errors.New("name should contain less than 13 characters")
	}
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-') {
			return errors.New("name should only contain alphanumeric characters, underscores, or hyphens")
		}
	}
	return nil
}

// ValidMessage validates if a message is not empty and contains only printable ASCII characters.
func ValidMessage(message string) (string, error) {
	if message == "" {
		return "", nil
	}
	for _, char := range message {
		if char < 32 || char > 126 {
			return "", errors.New("message should only contain printable ascii characters")
		}
	}
	return message, nil
}

// SetPrefix creates a formatted string with a timestamp and username,
func SetPrefix(name string) string {
	timestamp := time.Now().Format(time.DateTime)
	return fmt.Sprintf("[%s][%s]:", timestamp, name)
}
