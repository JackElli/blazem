package user

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Password string `json:"password"`
}
