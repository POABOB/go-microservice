package main

import (
	"config/conf"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	http.HandleFunc("/resumes", func(w http.ResponseWriter, req *http.Request) {
		//q := events.goreq.URL.Query().Get("q")
		_, _ = fmt.Fprintf(w, "個人資訊：\n")
		_, _ = fmt.Fprintf(w, "姓名：%s，\n性别：%s，\n年齡 %d!", conf.Resume.Name, conf.Resume.Sex, conf.Resume.Age)                                           //这个写入到w的是输出到客户端的
		_, _ = fmt.Fprintf(w, "姓名：%s，\n性别：%s，\n年齡 %s!", viper.GetString("resume.name"), viper.GetString("resume.sex"), viper.GetString("resume.age")) //这个写入到w的是输出到客户端的
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
