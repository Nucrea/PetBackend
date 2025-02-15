package models

type UserDTO struct {
	Id            string
	Email         string
	EmailVerified bool
	Secret        string
	Name          string
}

type UserUpdateDTO struct {
	Secret string
	Name   string
}
