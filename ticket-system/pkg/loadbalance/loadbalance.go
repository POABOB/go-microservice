package loadbalance

import (
	"errors"
	"math/rand"

	"github.com/POABOB/go-microservice/ticket-system/pkg/common"
)

// 附載均衡
type LoadBalance interface {
	SelectService(service []*common.ServiceInstance) (*common.ServiceInstance, error)
}

// 隨機策略
type RandomLoadBalance struct {
}

func (loadBalance *RandomLoadBalance) SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error) {
	if len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	return services[rand.Intn(len(services))], nil
}

// 權重平滑策略
type WeightRoundRobinLoadBalance struct {
}

func (loadBalance *WeightRoundRobinLoadBalance) SelectService(services []*common.ServiceInstance) (best *common.ServiceInstance, err error) {
	if len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}

	total := 0
	for i := 0; i < len(services); i++ {
		w := services[i]
		if w == nil {
			continue
		}

		// CurrentWeight 會被動態修改
		w.CurrentWeight += w.Weight

		total += w.Weight
		// CurrentWeight 最大即是
		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}

	if best == nil {
		return nil, nil
	}

	// 被選中的權重要減去全部
	best.CurrentWeight -= total
	return best, nil
}
