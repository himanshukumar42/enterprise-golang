package models

import "github.com/himanshukumar42/enterprise/forms"

type Contacts struct {
	ID                int64  `db:"id, primarykey, autoincrement" json:"id"`
	UserID            int64  `db:"userID" json:"-"`
	FirstName         string `db:"firstName" json:"firstName"`
	LastName          string `db:"lastName" json:"lastName"`
	Email             string `db:"email" json:"email"`
	CountryCode       string `db:"countryCode" json:"countryCode"`
	Mobile            string `db:"mobile" json:"mobile"`
	EventNotification string `db:"eventNotification" json:"eventNotification"`
	Groups            string `db:"groups" json:"groups"`
	EventsType        string `db:"eventsType" json:"eventsType"`
	Status            string `db:"status" json:"status"`
	UpdatedAt         int64  `db:"updated_at" json:"-"`
	CreatedAt         int64  `db:"created_at" json:"-"`
}

type ContactsModel struct{}

func (c ContactsModel) Create(userID int64, form forms.CreateContactsForm) (articleID int64, err error)
