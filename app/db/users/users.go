package users

import (
	"context"
	"log"
	"time"
	"xadmin/app/db"
	"xadmin/conf"
	"xadmin/subscribe/clash"
	"xadmin/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbDelUser(userEmail string) {
	c := db.Client.Database("xary").Collection("users")
	_, err := c.DeleteOne(context.TODO(), bson.M{"userEmail": userEmail})
	if err != nil {
		return
	}
}

func CheckUserIsExsitAndAddUser(u conf.UserTraffic) (bool, error) {
	c := db.Client.Database("xary").Collection("users")

	result := conf.UserCollection{}
	err := c.FindOne(context.TODO(), bson.M{"userEmail": u.UserEmail}).Decode(&result)

	if err != nil {
		t := time.Now()
		f := t.Format("2006-01-02 15:04:05")
		p := t.Unix()

		config := utils.ReadXaryConfigFile()

		var pass string
		for _, val := range config.Inbounds[1].Settings.Clients {
			if val.Email == u.UserEmail {
				pass = val.Password
			}
		}

		subdir := utils.CheckAndDelCreateSub(u.UserEmail)
		clash.CreateVip(u.UserEmail, "vip1", subdir)

		nu := conf.UserCollection{
			Id:            primitive.NewObjectID(),
			UserEmail:     u.UserEmail,
			SubDir:        subdir,
			UserPassword:  pass,
			Uplink:        u.Uplink,
			Downlink:      u.Downlink,
			TotalUplink:   u.Uplink,
			TotalDownlink: u.Downlink,
			CreateTime:    f,
			TimeStamp:     p,
		}

		_, err = c.InsertOne(context.TODO(), nu)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func UserList(page int64, email string) ([]conf.UserCollection, int64, error) {
	c := db.Client.Database("xary").Collection("users")

	ulist := []conf.UserCollection{}
	e := []conf.UserCollection{}

	total, _ := c.CountDocuments(context.TODO(), bson.M{})
	filter := bson.M{}

	if (len(email)) > 0 {
		curtotal, _ := c.CountDocuments(context.TODO(), bson.M{"userEmail": email})
		total = curtotal

		filter = bson.M{"userEmail": email}
	}

	cur, err := c.Find(
		context.TODO(),
		filter,
		options.Find().SetSkip((page)*10).SetSort(bson.M{"timeStamp": -1}),
	)
	if err != nil {
		return e, total, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var result conf.UserCollection
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		ulist = append(ulist, result)
	}

	return ulist, total, nil
}

func init() {
	db.ConnectDB()
}
