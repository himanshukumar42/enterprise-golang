package forms

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type ContactsForm struct{}

type CreateContactsForm struct {
	FirstName          string `form:"firstName" json:"firstName" binding:"required,min=3,max=100"`
	LastName           string `form:"lastName" json:"lastName" binding:"required,min=3,max=100"`
	Email              string `form:"email" json:"email" binding:"required,email"`
	CountryCode        string `form:"countryCode" json:"countryCode" binding:"required"`
	Mobile             string `form:"mobile" json:"mobile"`
	EventsNotification string `form:"eventsNotification" json:"eventsNotification"`
	Groups             string `form:"groups" json:"groups"`
	EventsType         string `form:"eventsType" json:"eventsType"`
	Status             string `form:"status" json:"status"`
}

type UpdateContactsForm struct {
	FirstName          string `form:"firstName" json:"firstName" binding:"min=3,max=100"`
	LastName           string `form:"lastName" json:"lastName" binding:"min=3,max=100"`
	Email              string `form:"email" json:"email" binding:"email"`
	CountryCode        string `form:"countryCode" json:"countryCode"`
	Mobile             string `form:"mobile" json:"mobile"`
	EventsNotification string `form:"eventsNotification" json:"eventsNotification"`
	Groups             string `form:"groups" json:"groups"`
	EventsType         string `form:"eventsType" json:"eventsType"`
	Status             string `form:"status" json:"status"`
}

func (f ContactsForm) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "please enter the Name"
		}
		return errMsg[0]
	case "min", "max":
		return "Name should be between 3 to 100 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f ContactsForm) Email(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "please enter your email"
		}
		return errMsg[0]
	case "min", "max", "email":
		return "please enter a valid email"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f ContactsForm) CountryCode(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "please enter your email"
		}
		return errMsg[0]
	case "min", "max":
		return "please enter a valid country code"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f ContactsForm) Mobile(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "please enter your mobile number"
		}
		return errMsg[0]
	case "min", "max":
		return "please enter a valid mobile number"
	default:
		return "Something went wrong, please try again later"
	}
}

func (f ContactsForm) Create(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "FirstName" {
				return f.Name(err.Tag())
			}
			if err.Field() == "LastName" {
				return f.Name(err.Tag())
			}
			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}
			if err.Field() == "CountryCode" {
				return f.CountryCode(err.Tag())
			}
			if err.Field() == "Mobile" {
				return f.Mobile(err.Tag())
			}
		}
	default:
		return "invalid request"
	}
	return "Something went wrong, please try again later"
}

func (f ContactsForm) Update(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "FirstName" {
				return f.Name(err.Tag())
			}
			if err.Field() == "LastName" {
				return f.Name(err.Tag())
			}
			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}
			if err.Field() == "CountryCode" {
				return f.CountryCode(err.Tag())
			}
			if err.Field() == "Mobile" {
				return f.Mobile(err.Tag())
			}
		}
	default:
		return "invalid request"
	}
	return "Something went wrong, please try again later"
}
