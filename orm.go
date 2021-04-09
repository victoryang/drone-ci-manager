package orm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
	"gorm.io/gorm"

	"git.snowballfinance.com/ops/sce-rolling/config"
)

const (
	MYSQL_CONN_STR = "root:123456@tcp(127.0.0.1:3306)/sce_rolling?parseTime=true&charset=utf8&loc=Local"
)

// MySQL orm
var ORM *gorm.DB

// TODO: add mysql table index
func init() {

	var err error
	dsn := config.GetString("MYSQL_CONN_STR")
	MySQL, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "sce_",	// table name prefix, table for `User` would be `t_users`
			SingularTable: false,	// use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err))
	}

	sqlDB, _ := MySQL.DB()
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(30)
}