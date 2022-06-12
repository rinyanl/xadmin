package traffic

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"xadmin/conf"

	"github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Address     = "207.148.103.112"
	Port        = 999
	TrafficList = []conf.UserTraffic{}
)

// 查询流量
func QueryUserTraffic() ([]conf.UserTraffic, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", Address, Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	c := command.NewStatsServiceClient(conn)

	filter := &command.QueryStatsRequest{
		Pattern: "",
		Reset_:  false,
	}

	res, err := c.QueryStats(context.Background(), filter)
	if err != nil {
		return []conf.UserTraffic{}, err
	}

	eReg := regexp.MustCompile(`[A-Za-z0-9.\-+_]+@[a-z0-9.\-+_]+\.[a-z]+`)
	tReg, _ := regexp.Compile(`uplink|downlink`)

	for _, val := range res.Stat {
		e := eReg.FindString(val.Name)
		types := tReg.FindString(val.Name)

		stat, index, update := CheckEmailIsExsit(e)

		if stat {
			if update == "uplink" {
				TrafficList[index].Uplink = val.Value
			}
			if update == "downlink" {
				TrafficList[index].Downlink = val.Value
			}
		}

		if !stat {
			cur := conf.UserTraffic{
				UserEmail: e,
				Uplink:    0,
				Downlink:  0,
			}
			if types == "uplink" {
				cur.Uplink = val.Value
			}
			if types == "downlink" {
				cur.Downlink = val.Value
			}
			TrafficList = append(TrafficList, cur)
		}
	}

	return TrafficList, nil

}

// 存在否、索引值、待更新项
func CheckEmailIsExsit(e string) (bool, int, string) {
	for index, check := range TrafficList {
		if e == check.UserEmail {
			if check.Uplink == 0 {
				return true, index, "uplink"
			}
			if check.Downlink == 0 {
				return true, index, "downlink"
			}
		}
	}
	return false, -1, ""
}
