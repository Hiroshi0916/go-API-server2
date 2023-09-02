package entities

type User struct {
	ID       string `json:"login_id"`
	Password string `json:"password"`
}
