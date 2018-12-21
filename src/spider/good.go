package spider

import (
	"fmt"
	"math"
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
	data, err := GetJsonData(url, "GET", headers, "")
	if err != nil{
		return
	}
	data = data.Get("data")
	info.Price, _ = data.Get("price").Int()
	info.RealPrice, _ = data.Get("realprice").Int()
	info.Stock, _ = data.Get("stock").String()
	info.IntStock, _ = data.Get("num").Int()
}

func GetGoodInfo(abiid string, token string, info *model.Good){
	url := fmt.Sprintf("http://b2carticleinfo.lib.cdn.srgow.com/api/v1/Article?languageid=1&abiid=%s", abiid)
	headers := map[string]string{"Accept": "application/json", "Authorization": "Bearer " + token}
	data, err := GetJsonData(url, "GET", headers, "")
	if err != nil{
		return
	}
	abiidInt, _ := data.Get("abiid").Int()
    info.Abiid = strconv.Itoa(abiidInt)
	info.MainName, _ = data.Get("mainname").String()
	info.Subtitle, _ = data.Get("subtitle").String()
	info.BrandName, _ = data.Get("brandname").String()
	info.BrandId, _ = data.Get("brandcode").String()
	info.CategoryName, _ = data.Get("categoryname").String()
	info.CategoryId, _ = data.Get("categorycode").String()
}

func UpdateGoodInfo() bool{
	var goodsBeMonitored []model.GoodBeMonitored
	model.Db.Find(&goodsBeMonitored)
	errCount := 0
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
			model.Info.Println("更新", good.Abiid, "成功")
		}else {
			errCount += 1
		}
	}
	return errCount == len(goodsBeMonitored)
}

func UpdateGoodInfoWithInterval(){
	model.Info.Println("第一次爬虫将在", time.Now().Add(time.Second * 60), "开始")
	time.Sleep(time.Second * 10)
	for{
		model.Info.Println("更新任务开始")
		conf := model.Conf{}
		model.Db.First(&conf)
		hour := conf.IntervalHour
		minute := conf.IntervalMinute
		second := 0
		if hour == 0 && minute == 0{
			hour = 0
			minute = 0
			second = 5
		}
		interval := time.Hour * time.Duration(hour) + time.Minute * time.Duration(minute) + time.Second * time.Duration(second)
		hasHttpError := UpdateGoodInfo()
		if hasHttpError{
			model.Error.Println("网络或者服务器出现问题，程序即将暂停10分钟, 下次开始时间为", time.Now().Add(time.Minute * 10))
			time.Sleep(time.Minute * 10)
			continue
		}
		var goods []model.Good
		model.Db.Find(&goods)
		util.CreatePath("spider_data")
		util.DomToExcel(goods, path.Join("data", "spider_data","output.xlsx"))
		model.Info.Println("更新任务完成， 下次更新将在", time.Now().Add(interval))
		model.Info.Println("提醒任务开始")
		Notice(conf)
		model.Info.Println("提醒任务完成")
		time.Sleep(interval)
	}
}

func GetNeedNotice()([]model.GoodBeNoticed, []model.GoodBeNoticed){
	var goodsBeMonitor []model.GoodBeMonitored
	var goodsNeedBeNoticed []model.GoodBeNoticed
	var goods []model.GoodBeNoticed
	model.Db.Find(&goodsBeMonitor)
	//var needNoticeGoods []model.Good
	conf := model.Conf{}
	model.Db.First(&conf)
	var goodCountInterval float64
	if conf.GoodCountInterval == 0{
		goodCountInterval = 5
	}else {
		goodCountInterval = float64(conf.GoodCountInterval)
	}
	for _, good := range goodsBeMonitor{
		var g model.Good
		model.Db.Find(&g, "abiid = ?", good.Abiid)
		var goodHistories []model.GoodHistory
		model.Db.Limit(2).Where("abiid = ?", good.Abiid).Order("update_time desc").Find(&goodHistories)
		if len(goodHistories) >= 2  && goodHistories[0].StockNum != goodHistories[1].StockNum && (goodHistories[0].StockNum > 0 || goodHistories[1].StockNum > 0) && math.Abs(float64(goodHistories[1].StockNum - goodHistories[0].StockNum)) >= goodCountInterval{
			goodsNeedBeNoticed = append(goodsNeedBeNoticed, model.GoodBeNoticed{Good:g, LastStock:goodHistories[1].StockNum})
			//needNoticeGoods = append(needNoticeGoods, g)
		}else{
			if g.IntStock < 0{
				g.IntStock = 0
			}
			goods = append(goods, model.GoodBeNoticed{Good:g, LastStock:goodHistories[1].StockNum})
		}
	}
	return goodsNeedBeNoticed, goods
}

func Notice(conf model.Conf){
	GoodsNeedBeNotice, Goods := GetNeedNotice()
	var message string
	if len(GoodsNeedBeNotice) == 0{
		model.Info.Println("商品库存没有变化，等待下一次更新")
	}else {
		message = "以下商品发送了变化\n"
		for _, good := range GoodsNeedBeNotice{
			message += good.Good.Abiid + "\t" +good.Good.MainName + ":" + strconv.Itoa(good.LastStock) + " -> " + strconv.Itoa(good.Good.IntStock) + "\n"
		}
		util.CreatePath("email")
		filename := path.Join("data", "email", "output.xlsx")
		filename = util.DomToExcelWithHightLight(GoodsNeedBeNotice, Goods, filename)
		if conf.Sender != "" || conf.SenderPwd != "" || conf.Receiver != "" {
			sendEmail(conf.Sender, conf.SenderPwd, conf.Receiver, message, filename)
		}else {
			model.Info.Println("未配置邮箱或邮箱错误")
		}
	}
}