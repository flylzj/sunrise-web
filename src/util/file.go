package util

import (
	"os"
	"path"
)

func CreatePath(filetype string){
	_ = os.Mkdir("data", os.ModePerm)
	_ = os.Mkdir(path.Join("data", filetype), os.ModePerm)
}