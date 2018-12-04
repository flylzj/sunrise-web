package spider

import (
	"fmt"
	"model"
	"path"
	"strconv"
	"time"
	"util"
	)

func GetAGood(abiid string) (model.Good, bool){
	info := model.Good{Abiid:abiid}
	token := GetToken()
	if token == "" || abiid == ""{
		return info, false
	}
	GetGoodInfo(abiid, token, &info)
	GetGoodPrice(abiid, token, &info)
	if info.Abiid != "0"{
		info.LastUpdateTime = int(time.Now().Unix())
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
	var goodsBeMonitored []model.GoodBeMonitored
	model.Db.Find(&goodsBeMonitored)
	for _, goodBeMonitored := range goodsBeMonitored{
		g, result := GetAGood(goodBeMonitored.Abiid)
		if result{
			var good model.Good
			model.Db.Find(&good, "abiid = ?", goodBeMonitored.Abiid)
			good.MainName = g.MainName
			good.Subtitle = g.Subtitle
			good.BrandId = g.BrandId
			good.BrandName = g.BrandName
			good.CategoryId = g.CategoryId
			good.CategoryName = g.CategoryName
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
	fmt.Println("第一次爬虫将在", time.Now().Add(time.Second * 60), "开始")
	time.Sleep(time.Second * 60)
	for{
		fmt.Println("更新任务开始")
		conf := model.Conf{}
		model.Db.First(&conf)
		hour := conf.IntervalHour
		minute := conf.IntervalMinute
		if hour == 0 && minute == 0{
			hour = 0
			minute = 2
		}
		interval := time.Hour * time.Duration(hour) + time.Minute * time.Duration(minute)
		UpdateGoodInfo()
		var goods []model.Good
		model.Db.Find(&goods)
		util.CreatePath("spider_data")
		util.DomToExcel(goods, path.Join("data", "spider_data","output.xlsx"))
		fmt.Println("更新任务完成， 下次更新将在", time.Now().Add(interval))
		fmt.Println("提醒任务开始")
		Notice(conf)
		fmt.Println("提醒任务完成")
		time.Sleep(interval)
	}
}

func GetNeedNotice()([]model.Good, []model.Good){
	var goodsBeMonitor []model.GoodBeMonitored
	var goods []model.Good
	model.Db.Find(&goodsBeMonitor)
	var needNoticeGoods []model.Good
	for _, good := range goodsBeMonitor{
		var g model.Good
		model.Db.Find(&g, "abiid = ?", good.Abiid)
		var goodHistories []model.GoodHistory
		model.Db.Limit(2).Where("abiid = ?", good.Abiid).Order("update_time desc").Find(&goodHistories)

		if len(goodHistories) >= 2 && goodHistories[0].StockNum != goodHistories[1].StockNum{
				needNoticeGoods = append(needNoticeGoods, g)
		}else{
			goods = append(goods, g)
		}
	}
	fmt.Println(needNoticeGoods)
	fmt.Println(goods)
	return needNoticeGoods, goods
}

func Notice(conf model.Conf){
	GoodNeedBeNotice, Goods := GetNeedNotice()
	var message string
	if len(GoodNeedBeNotice) == 0{
		message = "没有商品变化"
	}else {
		message += "以下商品发送了变化\n"
	}
	if len(Goods) == 0{
		fmt.Println("商品库存没有变化，等待下一次更新")
	}else {
		for _, good := range GoodNeedBeNotice{
			message += good.Abiid + "\t" +good.MainName + "\n"
		}
		util.CreatePath("email")
		filename := path.Join("data", "email", "output.xlsx")
		filename = util.DomToExcelWithHightLight(GoodNeedBeNotice, Goods, filename)
		if conf.Sender != "" || conf.SenderPwd != "" || conf.Receiver != "" {
			sendEmail(conf.Sender, conf.SenderPwd, conf.Receiver, message, filename)
		}else {
			fmt.Println("未配置邮箱或邮箱错误")
		}
	}
}