package spider

import "fmt"

func MainSpider(){
	go UpdateGoodInfoWithInterval()
	fmt.Println("爬虫模块加载完成")
}
