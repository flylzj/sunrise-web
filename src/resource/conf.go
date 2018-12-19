package resource

import (
	"model"
)

func SetTimeConf(hour, minute, goodCountInterval int)bool{
	conf := model.Conf{}
	model.Db.First(&conf)
	conf.IntervalHour = hour
	conf.IntervalMinute = minute
	conf.GoodCountInterval = goodCountInterval
	model.Db.Save(&conf)
	return true
}

func SetEmailConf(sender, sender_pwd, receiver string){
	conf := model.Conf{}
	model.Db.First(&conf)
	conf.Sender = sender
	conf.SenderPwd = sender_pwd
	conf.Receiver = receiver
	model.Db.Save(&conf)
}

func GetTimeConf()(int, int, int){
	conf := model.Conf{}
	model.Db.First(&conf)
	return conf.IntervalHour, conf.IntervalMinute, conf.GoodCountInterval
}

func GetEmailConf()(string, string, string){
	conf := model.Conf{}
	model.Db.First(&conf)
	return conf.Sender, conf.SenderPwd, conf.Receiver
}
