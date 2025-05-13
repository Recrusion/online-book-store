package models

type Users struct {
	ID       int         `json:db:"id"`
	Username string      `json:db:"username"`
	Email    string      `json:db:"email"`
	Password string      `json:"-"`
	Address  interface{} `json"db:"address"`
}

type SignUp struct {
	Username string `json:db:"id"`
	Email    string `json:db:"email"`
	Password string `json:"-"`
	Address  string `json"db:"address"`
}

type SignIn struct {
	Email    string `json:db:"email"`
	Password string `json:"-"`
}
