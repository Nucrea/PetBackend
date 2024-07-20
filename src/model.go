package src

type UserDTO struct {
	Id     string
	Login  string
	Secret string
	Name   string
}

type UserDAO struct {
	Id     string `json:"id"`
	Login  string `json:"login"`
	Secret string `json:"secret"`
	Name   string `json:"name"`
}
