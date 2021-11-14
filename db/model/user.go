package model

import (
	"fmt"

	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/db/utils"
)

type User struct {
	Name         string
	Password     string
	Role         string
	Organization string
	OrgID        string
}

type Channel struct {
	Channel  string
	Username string
}

//AddUser 添加用户
func (user *User) AddUser() error {
	stmt, err := utils.Db.Prepare("INSERT INTO webuser(name, password, role, organization) VALUES($1,$2,$3,$4)")
	if err != nil {
		fmt.Println("预编译出现异常：", err)
		return err
	}
	_, err = stmt.Exec(user.Name, user.Password, user.Role, user.Organization)
	if err != nil {
		fmt.Println("执行出现异常：", err)
		return err
	}
	return nil
}

//"cbit", "cbitpassword", "Jeweler", "Org1"
//user.Name,user.Password,user.Role,user.Organization

//GetUserByName 通过用户名字查询用户
func (user *User) GetUserByName() (*User, error) {
	row := utils.Db.QueryRow("select name,password,role,organization,orgid from webuser where name = $1", user.Name)
	u := &User{}
	err := row.Scan(&u.Name, &u.Password, &u.Role, &u.Organization, &u.OrgID)
	return u, err
}

//GetUsers 获取所有用户
func (user *User) GetUsers() ([]*User, error) {
	rows, err := utils.Db.Query("select * from webuser")
	if err != nil {
		return nil, err
	}

	var users []*User
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.Name, &u.Password, &u.Role, &u.Organization, &u.OrgID)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
