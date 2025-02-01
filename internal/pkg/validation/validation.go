package validation

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/google/uuid"
)

func EmailValidation(email string) (string, error) {
	//get email
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	emailErr := validation.Validate(email, validation.Required)
	if emailErr != nil {
		log.Println(emailErr)
		return "", emailErr
	}
	return email, nil
}
func PasswordValidation(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasLowerCase bool
		hasUpperCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			hasSpecial = true
		}

		if hasLowerCase && hasUpperCase && hasDigit && hasSpecial {
			break
		}
	}

	return hasLowerCase && hasUpperCase && hasDigit && hasSpecial
}
func NameValiddation(name string) bool {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ToUpper(string(name[0]))
	if len(name) < 2 || len(name) > 50 {
		return false
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

func ValidateUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	if !unicode.IsLetter(rune(username[0])) {
		return false
	}

	pattern := "^[a-zA-Z0-9_.-]+$"
	match, _ := regexp.MatchString(pattern, username)
	return match
}
