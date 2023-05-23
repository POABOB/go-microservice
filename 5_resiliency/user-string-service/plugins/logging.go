package plugins

import (
	"time"

	"use-string-service/service"

	"github.com/go-kit/kit/log"
)

// loggingMiddleware 主要是整合 Service 和 Logger 的模組
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// Middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) UseStringService(operationType, a, b string) (ret string, err error) {
	// 結束後打印 log
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "UseStringService",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.UseStringService(operationType, a, b)
	return ret, err
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
	// 結束後打印 log
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = mw.Service.HealthCheck()
	return
}
