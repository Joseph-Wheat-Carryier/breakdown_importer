package main

import (
	"CNXM_BRKD_READER/logger"
	"CNXM_BRKD_READER/minio_util"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"path/filepath"
)

func main1() {
	f, err := excelize.OpenFile("/Users/liumingju/code_world/demo/CNXM_BRKD_READER/resource/3333.xlsx")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	pics, err := f.GetPictures("Sheet1", "P2")

	pic := pics[0]
	dumpPic(pic.File)

	rows := computeRows(f)
	fmt.Println(rows)

}

func main() {
	var path string
	prompt := &survey.Input{
		Message: "请输入需要导入的excel路径",
	}
	survey.AskOne(prompt, &path)

	// 获取完整文件名
	filename := filepath.Base(path)
	// 获取扩展名
	ext := filepath.Ext(filename)
	// 去掉扩展名
	nameWithoutExt := filename[:len(filename)-len(ext)]

	// 初始化日志
	logger.InitLogger(filename)

	minio_util.InitMinio()
	ReadFile(path, nameWithoutExt)
}

func dumpPic(byteArr []byte) {
	// 文件路径
	filePath := "image1.png" // 替换为要写入的文件路径

	// 将字节数组写入文件
	err := ioutil.WriteFile(filePath, byteArr, 0644)
	if err != nil {
		fmt.Println("写入文件错误:", err)
		return
	}

	fmt.Println("文件写入完成.")
}
