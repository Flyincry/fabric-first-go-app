package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	Db  *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("postgres", "user=postgres password=gyy dbname=postgres sslmode=disable")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 清空数据，但不清除表格
func CleanTable(tableName string) error {
	stmt, err := Db.Prepare("TRUNCATE " + tableName + " RESTART IDENTITY CASCADE")
	if err != nil {
		fmt.Println("预编译出现异常：", err)
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("执行出现异常：", err)
		return err
	}
	return nil
}
