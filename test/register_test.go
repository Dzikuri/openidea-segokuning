package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	validation "github.com/itgelo/ozzo-validation"
)

func TestRegisterEmailSuccessValidation(t *testing.T) {
	// Example usage
	req := Request{
		CredentialType:  "email",
		CredentialValue: "testadmin@mailinator.com",
		Name:            "fullname user",
		Password:        "admindev",
	}

	if err := req.Validate(); err != nil {
		fmt.Println("Validation error:", err)
		return
	}

	fmt.Println("Request is valid!")
}

func TestRegisterEmailFailedValidation(t *testing.T) {
	// Example usage
	req := Request{
		CredentialType:  "email",
		CredentialValue: "testadmin@mailinator",
		Name:            "fullname user",
		Password:        "admindev",
	}

	if err := req.Validate(); err != nil {
		fmt.Println("Validation error:", err)
		return
	}

	fmt.Println("Request is valid!")
}

func TestHashPassword(t *testing.T) {
	password := "admindev"

	hashPassword, err := helper.HashPassword(password)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("original password: ", password)
	fmt.Println("Hash salt password: ", string(hashPassword))

	err = helper.ComparePassword(hashPassword, password)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Hash salt not match")
	} else {

		fmt.Println("Hash salt password is match")
	}

}

func TestComparePassword(t *testing.T) {
	hashPassword := "$2a$08$P7.Bu2UT9LAgbp7.vS4RHOGZzI3awH9zAdfb7TiipJZtqNsNE2k7S"
	password := "admindev"
	fmt.Println("original password: ", password)
	fmt.Println("Hash salt password: ", string(hashPassword))

	err := helper.ComparePassword(hashPassword, password)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Hash salt not match")
	} else {
		fmt.Println("Hash salt password is match")
	}
}

// Request represents the JSON request structure
type Request struct {
	CredentialType  string `json:"credentialType"`
	CredentialValue string `json:"credentialValue"`
	Name            string `json:"name"`
	Password        string `json:"password"`
}

// Validate validates the request fields
func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		// CredentialType validation
		validation.Field(&r.CredentialType, validation.Required, validation.In("email", "phone")),

		// CredentialValue validation based on CredentialType
		validation.Field(&r.CredentialValue, validation.Required, validation.By(func(value interface{}) error {
			switch r.CredentialType {
			case "email":
				if err := validation.Validate(value, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))); err != nil {
					return fmt.Errorf("invalid email format")
				}
			case "phone":
				if err := validation.Validate(value, validation.Match(regexp.MustCompile(`^\+[0-9]{7,13}$`))); err != nil {
					return fmt.Errorf("invalid phone format")
				}
			}
			return nil
		})),

		// Name validation
		validation.Field(&r.Name, validation.Required),

		// Password validation
		validation.Field(&r.Password, validation.Required),
	)
}
