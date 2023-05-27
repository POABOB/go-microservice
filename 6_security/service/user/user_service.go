package service

import (
	"context"
	"errors"

	"security/model"
)

// 錯誤訊息
var (
	ErrUserNotExist  = errors.New("username is not exist")
	ErrEmailNotExist = errors.New("email is not exist")
	ErrPassword      = errors.New("invalid password")
)

// UserDetails service interface
type UserDetailsService interface {
	// Get UserDetails By username or email
	GetUserDetailByUsernameOrEmail(ctx context.Context, username, email, password string) (*model.UserDetails, error)
}
