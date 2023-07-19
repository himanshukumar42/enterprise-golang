package models

import (
	"encoding/json"
	"errors"

	"github.com/himanshukumar42/enterprise/db"
	"github.com/himanshukumar42/enterprise/forms"
	"github.com/lib/pq"
)

type Contacts struct {
	ID                 int64           `db:"id, primarykey, autoincrement" json:"id"`
	UserID             int64           `db:"user_id" json:"-"`
	FirstName          string          `db:"first_name" json:"firstName"`
	LastName           string          `db:"last_name" json:"lastName"`
	Email              string          `db:"email" json:"email"`
	CountryCode        string          `db:"country_code" json:"countryCode"`
	Mobile             string          `db:"mobile_number" json:"mobile"`
	EventsNotification string          `db:"events_notification" json:"eventsNotification"`
	Groups             json.RawMessage `db:"groups" json:"groups"`
	EventsType         json.RawMessage `db:"event_types" json:"eventsType"`
	Status             string          `db:"status" json:"status"`
	User               *JSONRaw        `db:"user" json:"user"`
}

type ContactsModel struct{}

func (c ContactsModel) Create(userID int64, form forms.CreateContactsForm) (contactsID int64, err error) {
	err = db.GetDB().QueryRow(`
	INSERT INTO 
		public.contacts(user_id, first_name, last_name, email, 
		country_code, mobile_number, events_notification, groups, event_types, status) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		userID, form.FirstName, form.LastName, form.Email, form.CountryCode, form.Mobile,
		form.EventsNotification, pq.Array(form.Groups), pq.Array(form.EventsType), form.Status).Scan(&contactsID)
	return contactsID, err
}

func (c ContactsModel) One(userID, id int64) (contacts Contacts, err error) {
	err = db.GetDB().SelectOne(&contacts, `
	SELECT 
		c.id, c.user_id, c.first_name, c.last_name, 
		c.email, c.country_code, c.mobile_number,
		c.events_notification, 
		COALESCE(array_to_json(c.groups), '[]'::json) AS groups, 
		COALESCE(array_to_json(c.event_types), '[]'::json) AS event_types,
		c.status, json_build_object('id', u.id, 'name', u.name, 'email', u.email)
		AS user from public.contacts c LEFT JOIN public.user u on c.user_id = u.id 
		WHERE c.user_id=$1 AND c.id=$2 LIMIT 1`, userID, id)
	return contacts, err
}

func (c ContactsModel) All(userID int64) (contacts []DataList, err error) {
	_, err = db.GetDB().Select(&contacts, `
	SELECT COALESCE(array_to_json(array_agg(row_to_json(d))), '[]') 
		AS data, (SELECT row_to_json(n) FROM ( SELECT count(c.id) AS total FROM public.contacts
	 	AS c WHERE c.user_id=$1 LIMIT 1 ) n ) AS meta FROM ( SELECT c.id, c.first_name AS "firstName", c.last_name as "lastName", 
		c.email, c.country_code as "countryCode", c.mobile_number as "mobile", c.events_notification as "eventsNotification", c.groups as "groups", c.event_types as "eventsType", 
		c.status as "status", json_build_object('id', u.id, 'name', u.name, 'email', u.email) AS user FROM public.contacts c 
		LEFT JOIN public.user u on c.user_id = u.id WHERE c.user_id=$1 ORDER BY c.id DESC) d`,
		userID)
	return contacts, err
}

func (c ContactsModel) Update(userID, id int64, form forms.CreateContactsForm) (err error) {
	operation, err := db.GetDB().Exec(`UPDATE public.contacts 
	SET 
		first_name=$2, last_name=$3, 
		email=$4, country_code=$5, mobile_number=$6, 
		events_notification=$7, groups=$8, events_type=$9,
		stat=$10 WHERE id=$1`,
		id, form.FirstName, form.LastName, form.Email, form.CountryCode,
		form.Mobile, form.EventsNotification, form.Groups, form.EventsType, form.Status)
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
	SET 
		first_name = CASE WHEN $2 != '' THEN $2 ELSE first_name END,
		last_name = CASE WHEN $3 != '' THEN $3 ELSE last_name END,
		email = CASE WHEN $4 != '' THEN $4 ELSE email END,
		country_code = CASE WHEN $5 != '' THEN $5 ELSE country_code END,
		mobile_number = CASE WHEN $6 != '' THEN $6 ELSE mobile_number END,
		events_notification = CASE WHEN $7 != '' THEN $7 ELSE events_notification END,
		groups = CASE WHEN $8::text[] != '{}' THEN $8 ELSE groups END,
		event_types = CASE WHEN $9::text[] != '{}' THEN $9 ELSE event_types END,
		status = CASE WHEN $10 != '' THEN $10 ELSE status END
	WHERE id = $1;
	`, id, form.FirstName, form.LastName, form.Email, form.CountryCode, form.Mobile,
		form.EventsNotification, pq.Array(form.Groups), pq.Array(form.EventsType), form.Status)

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
