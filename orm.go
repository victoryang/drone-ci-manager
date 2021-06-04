package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
	"gorm.io/gorm"
)

const (
	MYSQL_CONN_STR = "xxxx:xxxxxxx@tcp(y.y.y.y:3306)/sce_rolling?parseTime=true&charset=utf8&loc=Local"
)

// MySQL orm
var ORM *gorm.DB

// TODO: add mysql table index
func init() {

	var err error
	dsn := MYSQL_CONN_STR
	ORM, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "sce_",	// table name prefix, table for `User` would be `t_users`
			SingularTable: false,	// use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err))
	}

	sqlDB, _ := ORM.DB()
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(30)
}
