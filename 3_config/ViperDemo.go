package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

var Resume ResumeInformation

// init() 建構子，初始化函數
func init() {
	// 獲取環境變數
	viper.AutomaticEnv()
	// 設定配置名稱、路徑
	initDefault()

	// 讀取 config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}
	if err := sub("ResumeInformation", &Resume); err != nil {
		log.Fatal("Fail to parse config", err)
	}
}
func initDefault() {
	// 文件名稱
	viper.SetConfigName("resume")
	// 文件路就
	viper.AddConfigPath("./config/")
	// windows環境為%GOPATH，linux環境為$GOPATH
	viper.AddConfigPath("$GOPATH/src/")
	// 設定文件類型
	viper.SetConfigType("yaml")
}
func main() {
	fmt.Printf(
		"姓名: %s\n興趣: %s\n性别: %s \n年齡: %d \n",
		Resume.Name,
		Resume.Habits,
		Resume.Sex,
		Resume.Age,
	)
	// 讀取YAML
	parseYaml(viper.GetViper())
	// 該變數是否為某值
	fmt.Println(Contains("Basketball", Resume.Habits))
}

func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	// 只要是Slice、Array多筆，就for迴圈判斷屬性
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	// 如果是Map，就判斷值有沒有效
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

type ResumeInformation struct {
	Name   string
	Sex    string
	Age    int
	Habits []interface{}
}

type ResumeSetting struct {
	RegisterTime      string
	Address           string
	ResumeInformation ResumeInformation
}

// 使用 Viper 記憶體中的環境變數，反序列化為Struct
func parseYaml(v *viper.Viper) {
	var resumeConfig ResumeSetting
	if err := v.Unmarshal(&resumeConfig); err != nil {
		fmt.Printf("err:%s", err)
	}
	fmt.Println("resume config:\n ", resumeConfig)
}

// 將 key 解析出 sub-tree，然後反序列化為 Struct
func sub(key string, value interface{}) error {
	log.Printf("配置文件的前綴為：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
