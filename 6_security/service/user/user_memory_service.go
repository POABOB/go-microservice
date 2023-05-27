package service

import (
	"context"

	"security/model"
)

// 記憶體儲存 Class
type InMemoryUserDetailsService struct {
	userDetailsDict map[string]*model.UserDetails
}

// 根據 Username or Email 獲取 User Detail
func (service *InMemoryUserDetailsService) GetUserDetailByUsernameOrEmail(ctx context.Context, username, email, password string) (*model.UserDetails, error) {
	var (
		userDetails *model.UserDetails
		ok          bool
	)

	// 先找 Email
	if email != "" {
		userDetails, ok = service.userDetailsDict[email]
	}
	// Email 找不到，從 Username 找
	if !ok && username != "" {
		userDetails, ok = service.userDetailsDict[username]
	}

	// 如果其中一個有找到
	if ok {
		// 比較密碼是否相同
		if userDetails.Password == password {
			return userDetails, nil
		}

		return nil, ErrPassword
	}

	return nil, ErrUserNotExist
}

// 初始化
func NewInMemoryUserDetailsService(userDetailsList []*model.UserDetails) *InMemoryUserDetailsService {
	// 建立 HashMap
	userDetailsDict := make(map[string]*model.UserDetails)

	// 如果有值，那就依序傳入
	if len(userDetailsList) > 0 {
		for _, value := range userDetailsList {
			userDetailsDict[value.Username] = value
		}
	}

	// 返回實例
	return &InMemoryUserDetailsService{
		userDetailsDict: userDetailsDict,
	}
}
