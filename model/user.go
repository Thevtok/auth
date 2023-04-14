package model

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	C_Username string `json:"c_username"`
}
