package xusers

import (
	"fmt"
	"time"
	"xadmin/app/db/users"
	"xadmin/conf"
	"xadmin/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ReadXaryConfig = conf.XaryConfig{}
)

func XaryAddUser(userPassword, userEmail string) error {
	if CheckXaryUserIsExsit(userEmail) != nil {
		return fmt.Errorf("xary 添加用户失败、已存在")
	}

	nu := conf.Clients{
		Password: userPassword,
		Flow:     "xtls-rprx-direct",
		Level:    0,
		Email:    userEmail,
	}

	for i, val := range ReadXaryConfig.Inbounds {
		if val.Port != 3008 {
			for j := 0; j < len(val.Settings.Clients); j++ {
				ReadXaryConfig.Inbounds[i].Settings.Clients = append(
					val.Settings.Clients,
					[]conf.Clients{
						nu,
					}...,
				)

			}
		}
	}

	utils.SaveXaryConfig(ReadXaryConfig)
	return nil
}

func XaryDelUser(userEmail string) error {
	if CheckXaryUserIsExsit(userEmail) != nil {
		for i, val := range ReadXaryConfig.Inbounds {
			if val.Tag != "api" {

				for j, user := range val.Settings.Clients {
					if user.Email == userEmail {
						ReadXaryConfig.Inbounds[i].Settings.Clients = append(val.Settings.Clients[:j], val.Settings.Clients[j+1:]...)
					}
				}

			}
		}
		users.DbDelUser(userEmail)
		utils.SaveXaryConfig(ReadXaryConfig)
		return nil
	}

	return fmt.Errorf("xary 删除用户失败、未找到用户")
}

func XaryAddInbound(port int) error {
	if CheckXaryInboundIsExsit(port) != nil {
		return fmt.Errorf("xary 添加入站失败、已存在")
	}

	tag := fmt.Sprintf("%d", time.Now().Unix())
	ni := []conf.Inbounds{
		{
			Tag:      tag,
			Port:     port,
			Protocol: "trojan",
			Settings: conf.Settings{
				Clients:    ReadXaryConfig.Inbounds[1].Settings.Clients,
				Decryption: "none",
				Fallbacks: []conf.Fallbacks{
					{
						Dest: 8080,
					},
				},
			},
			StreamSettings: &conf.StreamSettings{
				Network:  "tcp",
				Security: "xtls",
				XtlsSettings: &conf.XtlsSettings{
					AllowInsecure: false,
					MinVersion:    "1.2",
					Alpn: []string{
						"http/1.1",
					},
					Certificates: []conf.Certificates{
						{
							CertificateFile: "/etc/nginx/xray_cert/xray.crt",
							KeyFile:         "/etc/nginx/xray_cert/xray.key",
						},
					},
				},
			},
		},
	}

	ReadXaryConfig.Inbounds = append(ReadXaryConfig.Inbounds, ni...)
	utils.SaveXaryConfig(ReadXaryConfig)
	return nil
}

func XaryDelInbound(port int) error {
	if CheckXaryInboundIsExsit(port) != nil {
		for i, val := range ReadXaryConfig.Inbounds {
			if val.Tag != "api" {
				if val.Port == port {
					ReadXaryConfig.Inbounds = append(ReadXaryConfig.Inbounds[:i], ReadXaryConfig.Inbounds[i+1:]...)
					utils.SaveXaryConfig(ReadXaryConfig)
					return nil
				}

			}
		}
	}

	return fmt.Errorf("xary 删除入站失败、未找到")
}

func InboundList() ([]conf.InboundCollection, error) {
	config := utils.ReadXaryConfigFile()

	ilist := []conf.InboundCollection{}
	e := []conf.InboundCollection{}

	if config.Inbounds != nil {
		for _, val := range config.Inbounds {
			utotal := len(val.Settings.Clients)
			cur := conf.InboundCollection{
				Id:        primitive.NewObjectID(),
				Tag:       val.Tag,
				Port:      val.Port,
				Protocol:  val.Protocol,
				UserTotal: utotal,
			}
			ilist = append(ilist, cur)
		}
		return ilist, nil
	}
	return e, fmt.Errorf("入站列表为空")
}

func CheckXaryUserIsExsit(userEmail string) error {
	for _, val := range ReadXaryConfig.Inbounds {
		if val.Tag != "api" {
			for _, cli := range val.Settings.Clients {
				if cli.Email == userEmail {
					return fmt.Errorf("xary 添加用户失败、已存在")
				}
			}
		}
	}

	return nil
}

func CheckXaryInboundIsExsit(port int) error {
	for _, val := range ReadXaryConfig.Inbounds {
		if val.Port == port {
			return fmt.Errorf("xary 添加入站失败、已存在")
		}
	}
	return nil
}

func init() {
	ReadXaryConfig = utils.ReadXaryConfigFile()
}
