package models

type UserDTO struct {
	Id            string
	Email         string
	EmailVerified bool
	Secret        string
	FullName      string
}

type UserUpdateDTO struct {
	Secret   string
	FullName string
}
