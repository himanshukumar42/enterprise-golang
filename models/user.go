package models

import (
	"github.com/himanshukumar42/enterprise/db"
	"github.com/himanshukumar42/enterprise/forms"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64  `db:"id, primarykey, autoincrement" json:"id"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	Name      string `db:"name" json:"name"`
	UpdatedAt int64  `db:"updatedAt" json:"-"`
	CreatedAt int64  `db:"createdAt" json:"-"`
}

// UserModel..
type UserModel struct{}

var authModel = new(AuthModel)

func (m UserModel) Login(form forms.LoginForm) (user User, token Token, err error) {
	err = db.GetDB().SelectOne(&user, "SELECT id, email, password, name, updated_at, created_at FROM public.user WHERE email=LOWER($1) LIMIT 1", form.Email)
	if err != nil {
		return user, token, err
	}

	// Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	if err != nil {
		return user, token, err
	}

	// Generate the JWT auth token
	tokenDetails, err := authModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	saveErr := authModel.CreateAuth(user.ID, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}
	return user, token, nil
}
