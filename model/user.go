package model

type User struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"is_admin"`
	Uuid      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
