package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}
