package spider

import (
	"fmt"
	"model"
	"strconv"
	"time"
)

func GetAGood(abiid string) (model.Good, bool){
	info := model.Good{Abiid:abiid}
	token := GetToken()
	if token == ""{
		return info, false
	}
	GetGoodInfo(abiid, token, &info)
	GetGoodPrice(abiid, token, &info)
	if info.Abiid != "0"{
		return info, true
	}else {
		return info, false
	}
}

func GetGoodPrice(abiid string, token string, info *model.Good){
	url := fmt.Sprintf("http://srmemberapp.srgow.com/goods/prices/%s", abiid)
	headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
	data, _ := GetJsonData(url, "GET", headers, "")
	data = data.Get("data")
	info.Price, _ = data.Get("price").Int()
	info.RealPrice, _ = data.Get("realprice").Int()
	info.Stock, _ = data.Get("stock").String()
	info.IntStock, _ = data.Get("num").Int()
}

func GetGoodInfo(abiid string, token string, info *model.Good){
	url := fmt.Sprintf("http://b2carticleinfo.lib.cdn.srgow.com/api/v1/Article?languageid=1&abiid=%s", abiid)
	headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
	data, _ := GetJsonData(url, "GET", headers, "")
	abiidInt, _ := data.Get("abiid").Int()
    info.Abiid = strconv.Itoa(abiidInt)
	info.MainName, _ = data.Get("mainname").String()
	info.Subtitle, _ = data.Get("subtitle").String()
	info.BrandName, _ = data.Get("brandname").String()
	info.BrandId, _ = data.Get("brandcode").String()
	info.CategoryName, _ = data.Get("categoryname").String()
	info.CategoryId, _ = data.Get("categorycode").String()
}

func UpdateGoodInfo(){
	var goods []model.Good
	model.Db.Find(&goods)
	for _, good := range goods{
		g, result := GetAGood(good.Abiid)
		if result{
			good.Price = g.Price
			good.RealPrice = g.RealPrice
			good.Stock = g.Stock
			good.IntStock = g.IntStock
			good.LastUpdateTime = int(time.Now().Unix())
			model.Db.Save(&good)
			model.Db.Create(&model.GoodHistory{
				Abiid:good.Abiid,
				Stock:g.Stock,
				StockNum:g.IntStock,
				UpdateTime:int(time.Now().Unix()),
			})
			fmt.Println("更新", good.Abiid, "成功")
		}
	}
}

func UpdateGoodInfoWithInterval(){
	for{
		fmt.Println("更新任务开始")
		conf := &model.Conf{}
		model.Db.First(&conf)
		hour := conf.IntervalHour
		minute := conf.IntervalMinute
		if hour == 0 && minute == 0{
			hour = 0
			minute = 2
		}
		interval := time.Hour * time.Duration(hour) + time.Minute * time.Duration(minute)
		UpdateGoodInfo()
		fmt.Println("更新任务完成， 下次更新将在", time.Now().Add(interval))

		fmt.Println("提醒任务开始")
		Notice()
		fmt.Println("提醒任务完成")
		time.Sleep(interval)
	}
}

func GetNeedNotice()[]model.Good{
	var goods []model.Good
	model.Db.Find(&goods)
	var needNoticGoods []model.Good
	for _, good := range goods{
		var goodHistories []model.GoodHistory
		model.Db.Limit(2).Where("abiid = ?", good.Abiid).Order("update_time desc").Find(&goodHistories)
		if len(goodHistories) >= 2{
			if goodHistories[0].StockNum != goodHistories[1].StockNum{
				needNoticGoods = append(needNoticGoods, good)
			}
		}
	}
	return needNoticGoods
}

func Notice(){
	Goods := GetNeedNotice()
	if len(Goods) == 0{
		fmt.Println("商品库存没有变化，等待下一次更新")
	}else {
		for _, good := range Goods{

		}
	}
}
