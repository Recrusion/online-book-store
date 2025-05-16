package models

type Users struct {
	ID       int         `json:db:"id"`
	Username string      `json:db:"username"`
	Email    string      `json:db:"email"`
	Password string      `json:db:"password"`
	Address  interface{} `json"db:"address"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
