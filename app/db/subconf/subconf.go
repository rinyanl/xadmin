package subconf

import (
	"context"
	"fmt"
	"xadmin/app/db"
	"xadmin/conf"
	"xadmin/subscribe/clash"
	"xadmin/subscribe/quantumultx"
	"xadmin/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 重置所有用户的订阅
func ResetAllSub() error {
	c := db.Client.Database("xary").Collection("users")
	// 获取所有会员
	cur, err := c.Find(
		context.TODO(),
		bson.M{},
	)

	if err != nil {
		return fmt.Errorf("获取所有会员出错")
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		result := conf.UserCollection{}
		err := cur.Decode(&result)
		if err != nil {
			return fmt.Errorf("获取所有会员出错")
		}

		u := db.Client.Database("nothing_user").Collection("users")
		user := conf.UserCollection{}
		u.FindOne(context.TODO(), bson.M{"_id": result.Id}).Decode(&user)

		// 生成新目录
		dirname := utils.CheckAndDelCreateSub(user.UserEmail)

		// 生成匹配的 yaml 文件
		err = clash.CreateVip(user.UserEmail, "vip1", dirname)
		if err != nil {
			return fmt.Errorf("重置失败、yaml 用户：%v", user.UserEmail)
		}

		// 生成匹配的 txt 文件
		err = quantumultx.CreateVip1Quantumult(user.UserEmail, dirname)
		if err != nil {
			return fmt.Errorf("重置失败、txt 用户：%v", user.UserEmail)
		}

		// 更新vip订阅信息
		update := bson.M{"subDir": dirname}
		c.FindOneAndUpdate(context.TODO(), bson.M{
			"_id": result.Id,
		}, bson.M{"$set": update})

		if err != nil {
			// fmt.Println("失败、重置订阅")
			return fmt.Errorf("更新订阅失败、txt 用户：%v", fmt.Sprint(err))
		}

	}

	return nil
}

// 新增 clash 规则
func CreateClashRule(p conf.Proxies) error {
	c := db.Client.Database("xary").Collection("subconf")

	nc := conf.Proxies{
		Id:     primitive.NewObjectID(),
		Name:   p.Name,
		Server: p.Server,
		Port:   p.Port,
	}

	_, err := c.InsertOne(context.TODO(), nc)
	if err != nil {
		return fmt.Errorf("添加 clash 规则出错")
	}
	return nil
}

// 编辑 clash 规则
func EditClashRule(p conf.ProxiesJson) error {
	c := db.Client.Database("xary").Collection("subconf")
	pid, _ := primitive.ObjectIDFromHex(p.Id)
	update := bson.M{"name": p.Name, "server": p.Server, "port": p.Port}
	c.FindOneAndUpdate(context.TODO(), bson.M{"_id": pid}, bson.M{"$set": update})

	return nil
}

// 所有 clash 规则列表
func ClashRulelist() ([]conf.Proxies, error) {
	c := db.Client.Database("xary").Collection("subconf")
	// total, _ := c.CountDocuments(context.TODO(), bson.M{})

	carr := []conf.Proxies{}
	e := []conf.Proxies{}

	cur, err := c.Find(
		context.TODO(),
		bson.M{},
	)
	if err != nil {
		return e, fmt.Errorf("获取clash 规则出错")
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		result := conf.Proxies{}
		err := cur.Decode(&result)
		if err != nil {
			return e, fmt.Errorf("获取clash 规则出错")
		}
		carr = append(carr, result)
	}
	return carr, nil
}

// 删除规则
func DelClashRule(p string) error {
	c := db.Client.Database("xary").Collection("subconf")

	pid, _ := primitive.ObjectIDFromHex(p)
	_, err := c.DeleteOne(context.TODO(), bson.M{"_id": pid})
	if err != nil {
		return fmt.Errorf("删除规则失败")
	}

	return nil
}

func init() {
	db.ConnectDB()
}
