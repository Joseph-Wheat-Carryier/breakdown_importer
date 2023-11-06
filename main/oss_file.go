package main

import "time"

type OssFile struct {
	Id         string `gorm:"primaryKey"`
	FileName   string
	Url        string
	CreateBy   string
	CreateTime time.Time
	UpdateBy   string
	UpdateTime time.Time
}

func (o *OssFile) TableName() string {
	return `oss_file`
}

func NewOssFile(id string, fileName string, url string) (file *OssFile) {
	file = new(OssFile)
	file.Id = id
	file.FileName = fileName
	file.Url = url
	file.CreateBy = "Import"
	file.CreateTime = time.Now()
	file.UpdateBy = "Import"
	file.UpdateTime = time.Now()
	return file
}
