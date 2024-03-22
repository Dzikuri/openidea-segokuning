package model

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/itgelo/ozzo-validation/v4/is"
	uuid "github.com/satori/go.uuid"
)

type UserResponse struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	Name      string    `json:"name,omitempty"`
	ImageUrl  string    `json:"imageUrl,omitempty"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserAuthResponse struct {
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Name        string `json:"name,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

type UserCredentialType string

var UserCredentialTypes []interface{} = []interface{}{Email, Phone}

const (
	Email UserCredentialType = "email"
	Phone UserCredentialType = "phone"
)

type UserAuthRequest struct {
	CredentialType  UserCredentialType `json:"credentialType"`
	CredentialValue string             `json:"credentialValue"`
	Name            string             `json:"name"`
	Password        string             `json:"password"`
}

type UserLoginRequest struct {
	CredentialType  UserCredentialType `json:"credentialType"`
	CredentialValue string             `json:"credentialValue"`
	Password        string             `json:"password"`
}

type UserLinkEmailRequest struct {
	Email string    `json:"email"`
	Id    uuid.UUID `json:"id,omitempty"`
}

type UserLinkPhoneRequest struct {
	Phone string    `json:"phone"`
	Id    uuid.UUID `json:"id,omitempty"`
}

type UserUpdateAccount struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name"`
	ImageUrl string    `json:"imageUrl"`
}

func (p UserUpdateAccount) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.ImageUrl, validation.Required, is.URL),
	)
}

func (p UserLinkPhoneRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Phone, validation.Required.Error(ErrResRequiredField.Message), validation.Match(regexp.MustCompile(`^\+[0-9]{7,13}$`))),
	)
}

func (p UserLinkEmailRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required.Error(ErrResRequiredField.Message), is.Email),
	)
}

func (p UserAuthRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CredentialType, validation.Required.Error(ErrResRequiredField.Message), validation.In(UserCredentialTypes...)),

		// NOTE Validation Credential Value
		// CredentialValue validation based on CredentialType
		validation.Field(&p.CredentialValue, validation.Required, validation.By(func(value interface{}) error {
			switch p.CredentialType {
			case Email:
				if err := validation.Validate(value, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))); err != nil {
					return errors.New("Invalid email format")
				}
			case Phone:
				if err := validation.Validate(value, validation.Match(regexp.MustCompile(`^\+[0-9]{7,13}$`))); err != nil {
					return errors.New("Invalid phone number format")
				}
			}
			return nil
		})),
		// ===

		validation.Field(&p.Name, validation.Required.Error(ErrResRequiredField.Message), validation.Length(5, 50)),
		validation.Field(&p.Password, validation.Required.Error(ErrResRequiredField.Message), validation.Length(5, 15)),
	)
}

func (p UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CredentialType, validation.Required.Error(ErrResRequiredField.Message), validation.In(UserCredentialTypes...)),

		// NOTE Validation Credential Value
		// CredentialValue validation based on CredentialType
		validation.Field(&p.CredentialValue, validation.Required, validation.By(func(value interface{}) error {
			switch p.CredentialType {
			case Email:
				if err := validation.Validate(value, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))); err != nil {
					return errors.New("Invalid email format")
				}
			case Phone:
				if err := validation.Validate(value, validation.Match(regexp.MustCompile(`^\+[0-9]{7,13}$`))); err != nil {
					return errors.New("Invalid phone number format")
				}
			}
			return nil
		})),
		// ===

		validation.Field(&p.Password, validation.Required.Error(ErrResRequiredField.Message), validation.Length(5, 15)),
	)
}
