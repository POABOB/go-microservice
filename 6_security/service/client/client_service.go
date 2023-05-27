package service

import (
	"context"
	"errors"

	"security/model"
)

// 錯誤訊息
var (
	ErrClientNotExist = errors.New("clientId is not exist")
	ErrClientSecret   = errors.New("invalid clientSecret")
)

// ClientDtails service interface
type ClientDetailsService interface {
	GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error)
}
