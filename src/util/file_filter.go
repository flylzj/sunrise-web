package util

import "strings"

func FilterExcel(filename string) string{
	names := strings.Split(filename, ".")
	if len(names) != 2{
		return "文件名错误"
	}else if names[1] != "xlsx"{
		return "文件格式错误"
	}else {
		return "ok"
	}
}
