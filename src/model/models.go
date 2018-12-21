package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"log"
	"os"
)



var(
	Db *gorm.DB
	Info *log.Logger
	Error *log.Logger
)

func init_conf(db *gorm.DB){
	conf := Conf{IntervalHour:0, IntervalMinute:2,Sender:"",SenderPwd:"",Receiver:"", GoodCountInterval:5}
	db.Create(&conf)
}

func init() {
	errFile, err := os.OpenFile("err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil{
		log.Fatalln("打开文件失败", err.Error())
	}
	Info = log.New(os.Stdout, "Info:", log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stdout, errFile), "Error:", log.Ldate | log.Ltime | log.Llongfile)
	Db, err = gorm.Open("sqlite3", "test3.db")
	if err != nil {
		Error.Fatalln("数据库链接错误")
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
	Db.AutoMigrate(&Good{}, &GoodBeMonitored{}, &GoodHistory{}, &Conf{})
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
	GoodCountInterval	int
}

func (Conf) TableName()string{
	return "conf"
}


type GoodBeNoticed struct {
	Good	Good
	LastStock	int
}

