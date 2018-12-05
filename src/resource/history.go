package resource

import (
	"fmt"
	"model"
)

func GetGoodHistory(abiid string)[]model.GoodHistory{
	var history []model.GoodHistory
	model.Db.Group("stock_num").Order("update_time desc").Limit(20).Find(&history, "abiid = ?", abiid)
	fmt.Println(history)
	return history
}
