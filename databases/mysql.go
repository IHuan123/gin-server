package databases

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	dsn := "root:123456789@tcp(127.0.0.1:3306)/test_database?charset=utf8mb4&parseTime=True"
	DB,err = sqlx.Connect("mysql",dsn)
	if err != nil{
		panic(err)
	}
	// 设置最大连接数
	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(time.Minute*5)
	//err = DB.Ping()
	//if err != nil {
	//	panic(err)
	//}
}

//关闭数据库连接
func CloseMysql(){
	err := DB.Close()
	if err!=nil{
		panic(err)
	}
}
