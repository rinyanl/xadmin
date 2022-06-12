package query

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
	"xadmin/app/db"
	"xadmin/app/db/users"
	"xadmin/conf"
	"xadmin/xary/traffic"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateTraffic() {
	for {
		// isRun := cmd.RunstatusXary()
		// if !isRun {
		// 	log.Println("检测到 xary 未启动、将在3秒后重试", isRun)
		// 	time.Sleep(3 * time.Second)
		// 	go UpdateTraffic()
		// 	runtime.Goexit()
		// }

		utlist, err := traffic.QueryUserTraffic()
		if err != nil {
			log.Println("查询流量超时、将在3秒后重试", err)
			time.Sleep(3 * time.Second)
			go UpdateTraffic()
			runtime.Goexit()
		}

		c := db.Client.Database("xary").Collection("users")
		for _, val := range utlist {
			stats, err := users.CheckUserIsExsitAndAddUser(val)
			if err != nil {
				log.Println("更新流量到数据库、重写用户出错、将在3秒后重试", err)
				time.Sleep(3 * time.Second)
				go UpdateTraffic()
				runtime.Goexit()
			}

			if stats {
				update := bson.M{
					"uplink":        val.Uplink,
					"totalUplink":   val.Uplink,
					"downlink":      val.Downlink,
					"totalDownlink": val.Downlink,
				}
				c.FindOneAndUpdate(context.TODO(), bson.M{
					"userEmail": val.UserEmail,
				}, bson.M{"$set": update})
				// "$inc" // 测试先用set
			}

		}

		// log.Println(time.Now().Format("15:04:05"), "运行中同步流量")
		traffic.TrafficList = []conf.UserTraffic{}
		time.Sleep(3 * time.Second)
	}
}

func QueryTraffic() ([]conf.UserCollection, error) {
	c := db.Client.Database("xary").Collection("users")

	cur := []conf.UserCollection{}
	e := []conf.UserCollection{}

	r, _ := c.Find(context.TODO(), bson.M{})

	for r.Next(context.TODO()) {
		result := conf.UserCollection{}
		err := r.Decode(&result)
		if err != nil {
			return e, fmt.Errorf("解析用户信息失败")
		}

		cur = append(cur, result)
	}
	return cur, nil
}

func QueryTrafficAndClear() ([]conf.UserCollection, error) {
	c := db.Client.Database("xary").Collection("users")

	cur := []conf.UserCollection{}
	e := []conf.UserCollection{}

	r, _ := c.Find(context.TODO(), bson.M{})

	for r.Next(context.TODO()) {
		result := conf.UserCollection{}
		err := r.Decode(&result)
		if err != nil {
			return e, fmt.Errorf("解析用户信息失败")
		}
		cur = append(cur, result)

		result.Uplink = 0
		result.Downlink = 0
		_, err = c.UpdateOne(context.TODO(), bson.M{
			"userEmail": result.UserEmail,
		}, bson.M{"$set": result})

		if err != nil {
			return e, fmt.Errorf("查询后清空流量失败")
		}
	}
	return cur, nil
}

func init() {
	db.ConnectDB()
	go UpdateTraffic()
}
