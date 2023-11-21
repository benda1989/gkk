package db

import (
	"fmt"
	"gkk/config"
	"gkk/tool"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbs []*gorm.DB

func Close() {
	for _, v := range dbs {
		db, _ := v.DB()
		db.Close()
	}
}

// 参数：不传默认调用方配置文件，或者 配置文件（DbConf）或者 配置文件中的app的name
func Open(args ...any) *gorm.DB {
	var c *config.DbConf
	if len(args) > 0 {
		switch arg := args[0].(type) {
		case string:
			c = config.Get(arg).Db
		case *config.DbConf:
			c = arg
		}
	} else {
		c = config.Get().Db
	}
	var gia gorm.Dialector
	switch c.Type {
	case "mysql":
		gia = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name))
	default:
		gia = postgres.New(postgres.Config{
			DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				c.Host, c.User, c.Password, c.Name, c.Port),
			PreferSimpleProtocol: true,
		})
	}
	db, err := gorm.Open(gia, &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: false,
	}})
	tool.Exit(err, "failed to open database")
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdle)
	sqlDB.SetMaxOpenConns(c.MaxOpen)
	dbs = append(dbs, db)
	return db
}

func W(m map[string]any, tx *gorm.DB) *gorm.DB {
	for k, v := range m {
		if !tool.IsEmpty(v) {
			tx = tx.Where(k, v)
		}
	}
	return tx
}

func Rollback(tx *gorm.DB) {
	if err := recover(); err != nil {
		tx.Rollback()
		panic(err)
	}
}
