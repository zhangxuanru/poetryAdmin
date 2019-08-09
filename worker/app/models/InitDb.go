package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"poetryAdmin/worker/app/config"
	"time"
)

//数据库连接
func InitDb() (err error) {
	if err = orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		return err
	}
	if err = orm.RegisterDataBase("default", "mysql", config.G_Conf.DataSource); err != nil {
		return err
	}
	//SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制。
	//SetMaxIdleConns用于设置闲置的连接数
	orm.SetMaxIdleConns("default", 1000)
	orm.SetMaxOpenConns("default", 2000)
	orm.DefaultTimeLoc = time.UTC
	orm.Debug = true
	return nil
}
