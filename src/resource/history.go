package resource

import "model"

func GetGoodHistory(abiid string)[]model.GoodHistory{
	var history []model.GoodHistory
	model.Db.Order("update_time").Limit(20).Find(&history, "abiid = ?", abiid)
	return history
}
