package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type LoginInfo struct {
	UserID     int64     `json:"userId"`
	ClientIP   string    `json:"clientIP"`
	LoginState string    `json:"loginState"`
	LoginTime  time.Time `json:"loginTime"`
}

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	db := session.DB("user")
	c := db.C("login.info")
	err = c.Insert(NewLoginInfo(1234, "127.0.0.1", "success"))
	if err != nil {
		panic(err)
	}

	var loginHistory []LoginInfo
	err = c.Find(bson.M{"userid": 1234}).All(&loginHistory)
	for _, history := range loginHistory {
		fmt.Println(history)
	}

	fmt.Println("==========================")

	var lastLogin LoginInfo
	err = c.Find(bson.M{"userid": 1234}).Sort("-logintime").One(&lastLogin)
	if err != nil {
		panic(err)
	}
	fmt.Println(lastLogin)

}

func NewLoginInfo(id int64, clientIP string, loginState string) *LoginInfo{
	return &LoginInfo{
		UserID: id,
		ClientIP: clientIP,
		LoginState: loginState,
		LoginTime: time.Now(),
	}
}