package resource

import "model"

func GetGoodHistory(abiid string)[]model.GoodHistory{
	var history []model.GoodHistory
	model.Db.Find(&history, "abiid = ?", abiid).Order("update_time")
	return history
}
