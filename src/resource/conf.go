package resource

import (
	"fmt"
	"model"
)

func SetTimeConf(hour, minute int){
	conf := model.Conf{}
	model.Db.First(&conf)
	conf.IntervalHour = hour
	conf.IntervalMinute = minute
	model.Db.Save(&conf)
}

func SetEmailConf(sender, sender_pwd, receiver string){
	conf := model.Conf{}
	model.Db.First(&conf)
	conf.Sender = sender
	conf.SenderPwd = sender_pwd
	conf.Receiver = receiver
	model.Db.Save(&conf)
}

func GetTimeConf()(int, int){
	conf := model.Conf{}
	model.Db.First(&conf)
	fmt.Println(conf)
	return conf.IntervalHour, conf.IntervalMinute
}

func GetEmailConf()(string, string, string){
	conf := model.Conf{}
	model.Db.First(&conf)
	return conf.Sender, conf.SenderPwd, conf.Receiver
}
