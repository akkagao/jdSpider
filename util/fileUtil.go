package util

import (
	"fmt"
	"os"
)

/**
追加内容到指定为奸
 */
func AppendToFile(fileName, content string) {
	createFile(fileName)
	// 以追加模式，打开文件
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		f.WriteString(content)
	}
}

/**
创建文件
*/
func createFile(fileName string) {
	// 如果文件不存在则创建文件
	if _, err := os.Stat(fileName); err != nil && os.IsNotExist(err) {
		os.Create(fileName)
	}
}
