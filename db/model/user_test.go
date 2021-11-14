package model

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	//t.Run("测试添加用户：", testAddUser)
	t.Run("测试查询用户：", testGetUserByName)
	//t.Run("测试查询所有用户：", testGetUsers)
}

func testAddUser(t *testing.T) {
	fmt.Println("添加用户：")
	user := User{
		Name:         "ztz",
		Password:     "ztz123",
		Role:         "Jeweler",
		Organization: "org1",
	}
	user.AddUser()
}

func testGetUserByName(t *testing.T) {
	fmt.Println("查询用户：")
	user := User{
		Name: "gyy",
	}
	u, _ := user.GetUserByName()
	fmt.Println("查询结果：", u)
}

func testGetUsers(t *testing.T) {
	fmt.Println("查询所有用户：")
	user := &User{}
	us, _ := user.GetUsers()
	for k, v := range us {
		fmt.Printf("第%v个用户是：%v \n", k+1, v)
	}
}
