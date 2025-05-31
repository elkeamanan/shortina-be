package domain

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

func (r *RegisterUserRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Fullname == "" {
		return errors.New("fullname is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (r *RegisterUserRequest) ToNewUser() (*User, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New(),
		Email:    r.Email,
		Password: string(passwordBytes),
		Fullname: r.Fullname,
		Provider: LocalProvider,
	}, nil
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginUserRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type UpdateUserRequest struct {
	Fullname string `json:"fullname"`
}
