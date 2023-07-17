package models

import (
	"errors"

	"github.com/himanshukumar42/enterprise/db"
	"github.com/himanshukumar42/enterprise/forms"
)

type Contacts struct {
	ID                int64    `db:"id, primarykey, autoincrement" json:"id"`
	UserID            int64    `db:"user_id" json:"-"`
	FirstName         string   `db:"first_name" json:"firstName"`
	LastName          string   `db:"last_name" json:"lastName"`
	Email             string   `db:"email" json:"email"`
	CountryCode       string   `db:"country_code" json:"countryCode"`
	Mobile            string   `db:"mobile_number" json:"mobile"`
	EventNotification string   `db:"events_notification" json:"eventNotification"`
	Groups            string   `db:"groups" json:"groups"`
	EventsType        string   `db:"events_type" json:"eventsType"`
	Status            string   `db:"stat" json:"status"`
	User              *JSONRaw `db:"user" json:"user"`
}

type ContactsModel struct{}

func (c ContactsModel) Create(userID int64, form forms.CreateContactsForm) (contactsID int64, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.contacts(user_id, first_name, last_name, email, country_code, mobile_number, events_notification, groups, events_type, stat) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", userID, form.FirstName, form.LastName, form.Email, form.CountryCode, form.Mobile, form.EventsNotification, form.Groups, form.EventsType, form.Status).Scan(&contactsID)
	return contactsID, err
}

func (c ContactsModel) One(userID, id int64) (contacts Contacts, err error) {
	err = db.GetDB().SelectOne(&contacts, "SELECT c.id, c.first_name, c.last_name, c.email, c.country_code, c.mobile_number, c.events_notification, c.groups, c.events_type, c.stat, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user from public.contacts c LEFT JOIN public.user u on c.user_id = u.id WHERE c.user_id=$1 AND c.id=$2 LIMIT 1", userID, id)
	// err = db.GetDB().SelectOne(&article, "SELECT a.id, a.title, a.content, a.updated_at, a.created_at, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.article a LEFT JOIN public.user u ON a.user_id = u.id WHERE a.user_id=$1 AND a.id=$2 LIMIT 1", userID, id)

	return contacts, err
}

func (c ContactsModel) All(userID int64) (contacts []DataList, err error) {
	_, err = db.GetDB().Select(&contacts, "SELECT COALESCE(array_to_json(array_agg(row_to_json(d))), '[]') AS data, (SELECT row_to_json(n) FROM ( SELECT count(c.id) AS total FROM public.contacts AS c WHERE c.user_id=$1 LIMIT 1 ) n ) AS meta FROM ( SELECT c.id, c.first_name, c.last_name, c.email, c.country_code, c.mobile_number, c.events_notification, c.groups, c.events_type, c.stat, json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.contacts c LEFT JOIN public.user u on c.user_id = u.id WHERE c.user_id=$1 ORDER BY c.id DESC) d", userID)
	return contacts, err
}

func (c ContactsModel) Update(userID, id int64, form forms.CreateContactsForm) (err error) {
	operation, err := db.GetDB().Exec("UPDATE public.contacts SET first_name=$2, last_name=$3, email=$4, country_code=$5, mobile_number=$6, events_notification=$7, groups=$8, events_type=$9, stat=$10 WHERE id=$1", id, form.FirstName, form.LastName, form.Email, form.CountryCode, form.Mobile, form.EventsNotification, form.Groups, form.EventsType, form.Status)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("updated 0 records")
	}
	return err
}

func (c ContactsModel) PartialUpdate(userID, id int64, form forms.UpdateContactsForm) (err error) {
	operation, err := db.GetDB().Exec(`UPDATE public.contacts
	SET first_name = COALESCE($2, first_name),
		last_name = COALESCE($3, last_name),
		email = COALESCE($4, email),
		country_code = COALESCE($5, country_code),
		mobile_number = COALESCE($6, mobile_number),
		events_notification = COALESCE($7, events_notification),
		groups = COALESCE($8, groups),
		events_type = COALESCE($9, events_type),
		stat = COALESCE($10, stat)
	WHERE id=$1;`, id, form.FirstName, form.LastName, form.Email, form.CountryCode, form.Mobile, form.EventsNotification, form.Groups, form.EventsType, form.Status)

	if err != nil {
		return err
	}
	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("updated 0 records")
	}
	return err
}

func (c ContactsModel) Delete(userID, id int64) (err error) {
	operation, err := db.GetDB().Exec("DELETE FROM public.contacts WHERE id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}
	return err
}
