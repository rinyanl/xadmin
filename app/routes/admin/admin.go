package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"xadmin/app/db/query"
	"xadmin/app/db/subconf"
	"xadmin/app/db/users"
	"xadmin/cmd"
	"xadmin/conf"
	"xadmin/utils"
	"xadmin/xary/xusers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func XaryRoute(r *gin.Engine) {
	api := r.Group("/api")

	{
		api.GET("/querytrafficandclear", func(c *gin.Context) {
			cur, err := query.QueryTrafficAndClear()

			if err != nil {
				m := fmt.Sprint(err)
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "查询流量失败：" + m,
					"data":   "a",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "查询流量成功 (查询后清空)",
				"data":   cur,
			})
		})

		api.GET("/querytraffic", func(c *gin.Context) {
			cur, err := query.QueryTraffic()

			if err != nil {
				m := fmt.Sprint(err)
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "查询流量失败：" + m,
					"data":   "a",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "查询流量成功",
				"data":   cur,
			})
		})

		api.GET("/runningstatus", func(c *gin.Context) {

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "获取运行状态成功",
				"data": conf.RunStatus{
					Xary:  cmd.RunstatusXary(),
					Nginx: cmd.RunstatusNginx(),
					Mongo: cmd.RunstatusMongo(),
				},
			})
		})

		api.GET("/restartxary", func(c *gin.Context) {

			err := cmd.RestartXary()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "xary 重启成功",
			})
		})

		api.GET("/restartnginx", func(c *gin.Context) {

			err := cmd.RestartNginx()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "nginx 重启成功",
			})
		})

		api.GET("/restartmongo", func(c *gin.Context) {

			err := cmd.RestartMongo()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "mongo 重启成功",
			})
		})

		api.GET("/userlist", func(c *gin.Context) {
			email := c.Query("email")
			p := c.DefaultQuery("page", "0")

			page, _ := strconv.ParseInt(p, 10, 64)

			ulist, total, err := users.UserList(page, email)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "查询出错、请稍后再试",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "获取用户列表成功",
				"data": conf.Userlist{
					List:  ulist,
					Total: total,
					Page:  page,
				},
			})
		})

		api.GET("/inboundlist", func(c *gin.Context) {

			ulist, err := xusers.InboundList()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "查询出错、请稍后再试",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "获取入站列表成功",
				"data": conf.Userlist{
					List:  ulist,
					Total: int64(len(ulist)),
				},
			})
		})

		api.GET("/xaryconf", func(c *gin.Context) {
			xaryconf := utils.ReadXaryConfigFile()
			conf, _ := json.Marshal(xaryconf)

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "获取配置成功",
				"data":   string(conf),
			})
		})

		api.GET("/savexaryconf", func(c *gin.Context) {
			config := c.Query("config")
			// xaryconf := utils.ReadXaryConfigFile()
			// conf, _ := json.Marshal(xaryconf)

			var newConf conf.XaryConfig
			json.Unmarshal([]byte(config), &newConf)
			utils.SaveXaryConfig(newConf)

			fmt.Println(newConf)

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "保存配置成功",
				// "data":   string(conf),
			})
		})

		api.GET("/clashrulelist", func(c *gin.Context) {
			clist, err := subconf.ClashRulelist()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "查询出错、请稍后再试",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "获取规则列表成功",
				"data":   clist,
			})
		})

		api.GET("/delclashrule", func(c *gin.Context) {
			pid := c.Query("_id")
			err := subconf.DelClashRule(pid)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "删除出错、请稍后再试",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "删除规则成功",
			})
		})

		api.GET("/downloadsubscribec", func(c *gin.Context) {
			fd := c.Query("subDir")
			//打开文件

			exPath, _ := os.Getwd()
			file, err := os.Open(exPath + "/app/assets/clash/" + fd + "/" + "config.yaml")

			if err != nil {
				fmt.Println("打开文件错误", err)
				return
			}

			// 读取文件内容
			data, err := ioutil.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusAlreadyReported, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "未找到订阅信息、文件读取出错",
				})
				fmt.Println("读取失败")
				return
			}

			re := regexp.MustCompile(`(\S*)@`)
			name := re.FindString(fd)

			c.Header("Content-Type", "application/octet-stream")
			c.Header("Content-Disposition", "attachment; filename="+name+".yaml")
			c.Header("Content-Transfer-Encoding", "binary")
			// c.File(fd + "/" + fn)
			c.Data(http.StatusOK, "application/octet-stream", data)
		})

		api.GET("/login", func(c *gin.Context) {
			un := c.Query("userName")
			up := c.Query("userPassword")

			if un != "dijia" || up != "qwer1234poiu" {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "账号或密码错误",
				})
				return
			}

			u := conf.UserCollection{
				Id: primitive.NewObjectID(),
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "登陆成功",
				"data":   u,
				"token":  u,
			})
		})

		api.POST("/editclashrule", func(c *gin.Context) {
			json := conf.ProxiesJson{}
			c.Bind(&json)

			err := subconf.EditClashRule(json)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "编辑规则成功",
			})

		})

		api.POST("/adduser", func(c *gin.Context) {
			json := conf.Clients{}
			c.Bind(&json)

			if len(strings.TrimSpace(json.Email)) <= 0 || len(strings.TrimSpace(json.Password)) <= 0 {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "邮箱、密码不可为空",
				})
				return
			}

			err := xusers.XaryAddUser(json.Password, json.Email)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "添加用户成功",
				"email":  json.Email,
			})

		})

		api.POST("/deluser", func(c *gin.Context) {
			json := conf.Clients{}
			c.Bind(&json)

			if len(strings.TrimSpace(json.Email)) <= 0 {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "邮箱不可为空",
				})
				return
			}

			err := xusers.XaryDelUser(json.Email)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "删除用户成功",
				"email":  json.Email,
			})

		})

		api.POST("/addinbound", func(c *gin.Context) {
			json := conf.Port{}
			c.Bind(&json)

			if json.Port <= 0 {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "端口不能为0",
				})
				return
			}

			err := xusers.XaryAddInbound(json.Port)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "添加入站成功",
				"port":   json.Port,
			})

		})

		api.POST("/delinbound", func(c *gin.Context) {
			json := conf.Port{}
			c.Bind(&json)

			if json.Port <= 0 {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    "端口不能为0",
				})
				return
			}

			err := xusers.XaryDelInbound(json.Port)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "删除入站成功",
				"port":   json.Port,
			})

		})

		api.POST("/createclashrule", func(c *gin.Context) {
			json := conf.Proxies{}
			c.Bind(&json)

			err := subconf.CreateClashRule(json)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusAlreadyReported,
					"msg":    fmt.Sprint(err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"msg":    "添加规则成功",
			})

		})

	}
}
