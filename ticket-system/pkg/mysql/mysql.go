package mysql

import (
	"fmt"

	conf "github.com/POABOB/go-microservice/ticket-system/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

var (
	engin *gorose.Engin
	err   error
)

func InitMysql(hostMysql, portMysql, userMysql, pwdMysql, dbMysql string) {
	DbConfig := gorose.Config{
		// Default database configuration
		Driver: "mysql",                                                                                                              // Db 驅動 (mysql,sqlite,postgres,oracle,mssql)
		Dsn:    userMysql + ":" + pwdMysql + "@tcp(" + hostMysql + ":" + portMysql + ")/" + dbMysql + "?charset=utf8&parseTime=true", // 資料庫連線資訊
		Prefix: "",                                                                                                                   // Table prefix
		// 最多開啟連線池 Max open connections, default value 0 means unlimit.
		SetMaxOpenConns: 300,
		// 對多閒置連線 Max idle connections, default value is 1.
		SetMaxIdleConns: 10,
	}

	if engin, err = gorose.Open(&DbConfig); err != nil {
		conf.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
