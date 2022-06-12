package quantumultx

import (
	"fmt"
	"io/ioutil"
	"xadmin/subscribe/clash"

	"os"
)

var quantumultConfig = []byte("")

// 创建 vip1 quantumult 的订阅
func CreateVip1Quantumult(email, dirname string) error {
	vip1 := []string{}

	// 获取并生成数组
	sub, _, err := clash.ClashRulelist()
	if err != nil {
		return err
	}
	for _, v := range sub {
		str := "trojan=" + v.Server + ":" + fmt.Sprint(v.Port) + ", password=" + email + ", over-tls=true, tls-verification=true, fast-open=false, udp-relay=false, tag=" + v.Name + ""
		vip1 = append(vip1, str)
	}

	// 处理订阅为 txt
	str := ""
	for _, v := range vip1 {
		str += v + "\n"
	}

	quantumultConfig = []byte(str)
	err = SaveQuantumultSub(dirname)
	if err != nil {
		return err
	}
	return nil
}

// 保存文件
func SaveQuantumultSub(dirname string) error {
	exPath, _ := os.Getwd()
	err := ioutil.WriteFile(exPath+"/app/assets/quantumultx/"+dirname+"/config.txt", quantumultConfig, os.ModePerm)
	if err != nil {
		fmt.Printf("%v  错误", err)
		return err
	}
	return nil
}
