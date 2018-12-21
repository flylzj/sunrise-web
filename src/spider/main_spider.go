package spider

import "model"

func MainSpider(){
	UpdateGoodInfoWithInterval()
	model.Info.Println("爬虫模块加载完成")
}
