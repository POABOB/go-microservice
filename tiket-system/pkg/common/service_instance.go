package common

// 服務實例相關資訊結構體
type ServiceInstance struct {
	Host          string // Host
	Port          int    // Port
	Weight        int    // 權重
	CurrentWeight int    // 當前權重
	GrpcPort      int    // Grpc Port
}
