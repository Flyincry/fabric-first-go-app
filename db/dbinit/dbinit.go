package main

import (
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/db/model"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/db/utils"
)

func main() {
	// 清空数据库
	utils.CleanTable("webuser")
	// 创建初始用户
	user1 := model.User{
		Name:         "gyy",
		Password:     "gyy123",
		Role:         "Jeweler",
		Organization: "org1",
	}
	user1.AddUser()
	user2 := model.User{
		Name:         "ztz",
		Password:     "ztz123",
		Role:         "Bank",
		Organization: "org2",
	}
	user2.AddUser()
}
