package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initMySQL() (err error) {
	// DSN:Data Source Name
	dsn := "root:0123456@tcp(127.0.0.1:3306)/sql_demo"
	// 去初始化全局的db对象而不是新声明的一个db变量
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)
		return
	}
	db.SetConnMaxLifetime(time.Second * 10) //连接存活的最长时间
	db.SetMaxOpenConns(200)                 //最大连接数
	db.SetMaxIdleConns(10)                  //最大空闲连接数
	return
}

type user struct {
	id   int
	age  int
	name string
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	row := db.QueryRow(sqlStr, 1)
	// 如果不调用Scan这个方法，则连接不会释放，具体查看源码
	err := row.Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接，这只是为了避免下面for循环未执行完就奔溃，导致Next方法中没有关闭连接，但是是否会发生重复关闭的问题
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

func main() {
	if err := initMySQL(); err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)
	}
	// 做完错误检查之后，确保d不为nil
	defer db.Close() // 注意这行代码要写在上面err判断的下面
	fmt.Println("connect to db success")
	queryRowDemo()
	queryMultiRowDemo()
}
