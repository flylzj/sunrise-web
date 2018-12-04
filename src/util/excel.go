package util

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"model"
	"os"
	"strconv"
)



func ReadXlsx(filename string) (abiids []string){
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows{
		abiid := row.Cells[0]
		abiids = append(abiids, abiid.Value)
	}
	return
}

func DomToExcel(goods []model.Good, filename string)string{
	os.Remove(filename)
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("库存变化表")
	if err != nil{
		fmt.Println(err)
		return ""
	}
	row := sheet.AddRow()
	row.AddCell().Value = "abiid"
	row.AddCell().Value = "mainname"
	row.AddCell().Value = "price"
	row.AddCell().Value = "stock"
	row.AddCell().Value = "stock_num"
	for _, good := range goods{
		row = sheet.AddRow()
		row.AddCell().Value = good.Abiid
		row.AddCell().Value = good.MainName
		row.AddCell().Value = strconv.Itoa(good.Price)
		row.AddCell().Value = good.Stock
		row.AddCell().Value = strconv.Itoa(good.IntStock)
	}
	err = file.Save(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filename
}
