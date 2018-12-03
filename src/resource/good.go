package resource

import (
	"model"
	"spider"
)

func GetBeMonitoredGoods()[]model.Good{
	var goodsBeMonitor []model.GoodBeMonitored
	var goods []model.Good
	model.Db.Find(&goodsBeMonitor)
	for _, good := range goodsBeMonitor{
		goods = append(goods, SearchGood(good.Abiid))
	}
	return goods
}

func SearchGood(abiid string) model.Good{
	var good model.Good
	model.Db.Find(&good, "abiid = ?", abiid)
	return good
}

func SearchGoodBeMonitored(abiid string)model.GoodBeMonitored{
	var goodBeMonitor model.GoodBeMonitored
	model.Db.Find(&goodBeMonitor, "abiid = ?", abiid)
	return goodBeMonitor
}

func AddGood(abiid string)string{
	if SearchGoodBeMonitored(abiid).Abiid != ""{
		return "已存在"
	}
	var good model.Good
	var result bool
	good, result = spider.GetAGood(abiid)
	if result{
	    goodBeMonitor := model.GoodBeMonitored{Abiid:abiid}
	    model.Db.Create(&goodBeMonitor)
	    if SearchGood(abiid).Abiid == ""{
			model.Db.Create(&good)
		}
		return "添加成功"
	}else {
		return "添加失败"
	}
}