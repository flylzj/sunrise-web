package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Db *gorm.DB

func init_conf(db *gorm.DB){
	conf := Conf{IntervalHour:0, IntervalMinute:0,Sender:"",SenderPwd:"",Receiver:""}
	db.Create(&conf)
}

func init() {
	var err error
	Db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		panic("连接数据库失败")
	}
	if !Db.HasTable(&Good{}){
		Db.CreateTable(&Good{})
	}
	if !Db.HasTable(&GoodBeMonitored{}){
		Db.CreateTable(&GoodBeMonitored{})
	}
	if !Db.HasTable(&GoodHistory{}){
		Db.CreateTable(&GoodHistory{})
	}
	if !Db.HasTable(&Conf{}){
		Db.CreateTable(&Conf{})
		init_conf(Db)
	}
	return
}

type Good struct {
	ID				uint
	Abiid   		string
	MainName 		string
	Subtitle		string
	BrandId		    string
	BrandName		string
	CategoryId		string
	CategoryName	string
	Price			int
	RealPrice		int
	Stock           string
	IntStock		int
	LastUpdateTime	int
}

func (Good) TableName()string{
	return "good"
}

type GoodBeMonitored struct {
	ID				uint
	Abiid			string
}

func (GoodBeMonitored) TableName()string{
	return "good_be_monitored"
}

type GoodHistory struct {
	ID 		    uint
	Abiid	    string
	Stock       string
	StockNum	int
	UpdateTime	int
}

func (GoodHistory) TableName()string{
	return "good_history"
}

type Conf struct {
	ID				uint
	IntervalHour	int
	IntervalMinute	int
	Sender			string
	SenderPwd		string
	Receiver		string
}

func (Conf) TableName()string{
	return "conf"
}


