package models

import (
	"os"
)

//判断文件夹是否存在
func ExistsDir(dirname string) (bool, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(getwd + dirname)
	//os.Stat() 在文件夹存在的时候，error返回的是nil，这时候用os.IsExist(err)和os.IsNotExist(err) 都是false的，所以当文件夹存在即err == nil 时候，不用这两个方法再次判断。
	//当err != nil的时候，可以直接用os.IsNotExist(err) 或 os.IsExist(err) 方法判断是否不存在。
	if err == nil || !os.IsExist(err) && !os.IsNotExist(err) {
		return true, nil
	} else {
		return false, err
	}
}

//创建文件夹
func CreateDir(dirname string) (err error) {
	getwd, err := os.Getwd()
	if err != nil {
		return
	}
	err = os.Mkdir(getwd+dirname, os.ModePerm)
	if err != nil {
		return
	}
	return
}
