package requests

import "user/internal/validation"

type CreateUserRequestWrapper struct {
	Firstname string `validate:"required,min=2,max=32"`
	Lastname  string `validate:"required,min=2,max=32"`
	Email     string `validate:"required,email,max=128"`
	Password  string `validate:"required,min=8,max=32"`
}

func (req *CreateUserRequestWrapper) Validate() error {
	return validation.ValidateStruct(req)
}
