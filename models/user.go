package models

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

func (m UserModel) Login(from forms.LoginForm) (user User, token Token, err error) {

}
