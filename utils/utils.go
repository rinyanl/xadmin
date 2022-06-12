package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
	"xadmin/cmd"
	"xadmin/conf"
)

var RunTimeReload int64
var ReadAndWriteDir = `./config.json`

func ReloadDebounce(run int64, f func() error) {
	// 每次触发前会生成一个最新时间戳、一个存全局、一个传入当前函数
	// 3s 后如果全局等于函数内的时间戳、执行目标函数
	time.Sleep(3 * time.Second)
	if run == RunTimeReload {
		f()
	}
}

func ReadXaryConfigFile() conf.XaryConfig {
	file, err := os.Open(ReadAndWriteDir)
	e := conf.XaryConfig{}

	if err != nil {
		log.Println("读取 xary 配置文件失败", err)
		return e
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("读取 xary 配置文件失败", err)
		return e
	}

	ReadXaryConfig := conf.XaryConfig{}
	err = json.Unmarshal(data, &ReadXaryConfig)
	if err != nil {
		log.Println("解析 xary 配置文件内容失败", err)
		return e
	}

	return ReadXaryConfig
}

func SaveXaryConfig(config conf.XaryConfig) {
	data, err := json.Marshal(&config)
	if err != nil {
		log.Println("转换为 json 配置文件失败", err)
	}

	err = ioutil.WriteFile(ReadAndWriteDir, data, os.ModePerm)

	if err != nil {
		log.Println("保存 xary 配置文件失败", err)
		return
	}

	time.Sleep(200 * time.Millisecond)
	t := time.Now().Unix()
	RunTimeReload = t
	go ReloadDebounce(t, cmd.RestartXary)
}

// 查看目录是否存在、存在则删除
func CheckAndDelCreateSub(email string) string {
	exPath, _ := os.Getwd()
	t := time.Now().Unix()
	s := strconv.FormatInt(t, 10)

	// 查看目录是否有其他该用户的文件夹 quantumultx
	qfiles, _ := ioutil.ReadDir(exPath + "/app/assets/quantumultx/")
	for _, f := range qfiles {
		s, _ := regexp.MatchString(email, f.Name())
		// 如果存在则删除该目录
		if s {
			fmt.Println("删除quantumultx目录", f.Name())
			os.RemoveAll(exPath + "/app/assets/quantumultx/" + f.Name())
		}
	}

	// 查看目录是否有其他该用户的文件夹 clash
	cfiles, _ := ioutil.ReadDir(exPath + "/app/assets/clash/")
	for _, f := range cfiles {
		s, _ := regexp.MatchString(email, f.Name())
		// 如果存在则删除该目录
		if s {
			fmt.Println("删除clash目录", f.Name())
			os.RemoveAll(exPath + "/app/assets/clash/" + f.Name())
		}
	}

	// 创建用户订阅目录
	dirname := email + "_" + s
	os.Mkdir(exPath+"/app/assets/quantumultx/"+dirname, 0777)
	os.Chmod(exPath+"/app/assets/quantumultx/"+dirname, 0777)

	os.Mkdir(exPath+"/app/assets/clash/"+dirname, 0777)
	os.Chmod(exPath+"/app/assets/clash/"+dirname, 0777)

	return dirname

}
