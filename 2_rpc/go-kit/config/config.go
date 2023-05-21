package config

import (
	"log"
	"os"

	kitlog "github.com/go-kit/log"
)

var Logger *log.Logger
var KitLogger kitlog.Logger

func init() {
	// 定義Log紀錄錯誤，並定義時間格式
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	// 定義Kit Logger 設定相關錯誤
	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.DefaultTimestampUTC)
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)
}
