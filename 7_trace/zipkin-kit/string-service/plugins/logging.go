package plugins

import (
	"time"

	"zipkin-kit/string-service/service"

	"github.com/go-kit/log"
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

func (mw loggingMiddleware) Concat(a, b string) (ret string, err error) {
	// 結束後打印 log
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Concat",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Concat(a, b)
	return ret, err
}

func (mw loggingMiddleware) Diff(a, b string) (ret string, err error) {
	// 結束後打印 log
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Diff",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Diff(a, b)
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
