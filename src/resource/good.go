package resource

import (
	"fmt"
	"model"
	"spider"
	"util"
)

func GetBeMonitoredGoods()[]model.Good{
	var goodsBeMonitor []model.GoodBeMonitored
	var goods []model.Good
	model.Db.Find(&goodsBeMonitor)
	for _, good := range goodsBeMonitor{
		g := SearchGood(good.Abiid)
		if g.IntStock < 0{
			g.IntStock = 0
		}
		goods = append(goods, g)
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

func DeleteAGood(abiid string)string{
	goodBeMonitored := SearchGoodBeMonitored(abiid)
	if goodBeMonitored.Abiid == ""{
		return "不存在"
	}
	model.Db.Delete(&goodBeMonitored)
	return "ok"
}

func AddGoodInBatches(filename string)(results [][2]string){
	//var wg sync.WaitGroup
	abiids := util.ReadXlsx(filename)
	for _, abiid := range abiids{
		result := AddAbiid(abiid)
		fmt.Println(result)
		results = append(results, [2]string{abiid, result})
	}
	return
}


func AddAbiid(abiid string)string{
	if SearchGoodBeMonitored(abiid).Abiid != ""{
		return "已存在"
	}
	var good model.Good
	goodBeMonitor := model.GoodBeMonitored{Abiid:abiid}
	model.Db.Create(&goodBeMonitor)
	if SearchGood(abiid).Abiid == ""{
		good.Abiid = abiid
		model.Db.Create(&good)
	}
	return "添加成功"
}