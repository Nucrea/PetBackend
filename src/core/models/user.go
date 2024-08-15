package models

type UserDTO struct {
	Id     string
	Email  string
	Secret string
	Name   string
}

type UserUpdateDTO struct {
	Secret string
	Name   string
}
