package service

type Service interface {
	// 一般資料
	SimpleData(username string) string
	// Admin
	AdminData(username string) string
	// HealthCheck
	HealthCheck() bool
}

type CommonService struct {
}

// 實例化
func NewCommonService() *CommonService {
	return &CommonService{}
}

func (s *CommonService) SimpleData(username string) string {
	return "hello " + username + " ,simple data, with simple authority"
}

func (s *CommonService) AdminData(username string) string {
	return "hello " + username + " ,admin data, with admin authority"

}

// HealthCheck
func (s *CommonService) HealthCheck() bool {
	return true
}
