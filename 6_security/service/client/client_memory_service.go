package service

import (
	"context"

	"security/model"
)

// 記憶體儲存 Class
type InMemoryClientDetailsService struct {
	clientDetailsDict map[string]*model.ClientDetails // 使用字典當作 Session 使用
}

// 初始化
func NewInMemoryClientDetailService(clientDetailsList []*model.ClientDetails) *InMemoryClientDetailsService {
	// 建立 HashMap
	clientDetailsDict := make(map[string]*model.ClientDetails)

	// 如果傳遞進來的 List 不為空，那就依序傳入 HashMap
	if len(clientDetailsList) > 0 {
		for _, value := range clientDetailsList {
			clientDetailsDict[value.ClientId] = value
		}
	}

	// 返回實例
	return &InMemoryClientDetailsService{
		clientDetailsDict: clientDetailsDict,
	}
}

// 依據 Client ID 獲取 Client Detail
func (service *InMemoryClientDetailsService) GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error) {
	// HashMap 中查找
	if clientDetails, ok := service.clientDetailsDict[clientId]; ok {
		// 比對私鑰是否正確
		if clientDetails.ClientSecret == clientSecret {
			return clientDetails, nil
		}

		return nil, ErrClientSecret
	}

	return nil, ErrClientNotExist
}
