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
		return
	}
	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows{
		abiid := row.Cells[0]
		_, err := strconv.Atoi(abiid.Value)
		if err != nil{
			continue
		}
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

func DomToExcelWithHightLight(goodsNeedBeNoticed []model.GoodBeNoticed, goods []model.Good, filename string)string{
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
	row.AddCell().Value = "before"
	row.AddCell().Value = "now"

	style := xlsx.NewStyle()

	fill := *xlsx.NewFill("solid", "00FF0000", "FF000000")
	font := *xlsx.NewFont(12, "Verdana")
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")

	style.Fill = fill
	style.Font = font
	style.Border = border

	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = true

	for _, good := range goodsNeedBeNoticed{
		row = sheet.AddRow()
		row.SetHeight(20)
		cell1 := row.AddCell()
		cell1.Value = good.Good.Abiid
		cell2 := row.AddCell()
		cell2.Value = good.Good.MainName
		cell3 := row.AddCell()
		cell3.Value = strconv.Itoa(good.Good.Price)
		cell4 := row.AddCell()
		cell4.Value = good.Good.Stock
		cell5 := row.AddCell()
		cell5.Value = strconv.Itoa(good.LastStock)
		cell6 := row.AddCell()
		cell6.Value = strconv.Itoa(good.Good.IntStock)
		cell1.SetStyle(style)
		cell2.SetStyle(style)
		cell3.SetStyle(style)
		cell4.SetStyle(style)
		cell5.SetStyle(style)
		cell6.SetStyle(style)
		//row.AddCell().Value = good.Abiid
		//row.AddCell().Value = good.MainName
		//row.AddCell().Value = strconv.Itoa(good.Price)
		//row.AddCell().Value = good.Stock
		//row.AddCell().Value = strconv.Itoa(good.IntStock)
	}
	for _, good := range goods{
		row = sheet.AddRow()
		cell1 := row.AddCell()
		cell1.Value = good.Abiid
		cell2 := row.AddCell()
		cell2.Value = good.MainName
		cell3 := row.AddCell()
		cell3.Value = strconv.Itoa(good.Price)
		cell4 := row.AddCell()
		cell4.Value = good.Stock
		cell5 := row.AddCell()
		cell5.Value = strconv.Itoa(good.IntStock)
		cell6 := row.AddCell()
		cell6.Value = strconv.Itoa(good.IntStock)
		//row.AddCell().Value = good.Abiid
		//row.AddCell().Value = good.MainName
		//row.AddCell().Value = strconv.Itoa(good.Price)
		//row.AddCell().Value = good.Stock
		//row.AddCell().Value = strconv.Itoa(good.IntStock)
	}
	err = file.Save(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filename
}
