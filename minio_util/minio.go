package minio_util

import (
	"CNXM_BRKD_READER/logger"
	"github.com/minio/minio-go/v6"
	"io"
)

//var endpoint = "172.16.1.62:9000"

var endpoint = "192.168.8.3:9000"

//var accessKeyID = "admin"

var accessKeyID = "crrcdt"
var secretAccessKey = "KF@32rjb"
var useSSL = false
var bucketName = "cnxm"

var minioClient *minio.Client

func InitMinio() {
	var err error
	minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		logger.GetLogger().Error("创建 MinIO 客户端失败", err)
		return
	}
	logger.GetLogger().Debug("创建 MinIO 客户端成功")
}

func GetMinioClient() *minio.Client {
	if minioClient != nil {
		return minioClient
	}

	return minioClient
}

func UploadFile(objectName string, filePath string, reader io.Reader) error {
	// 指定上传文件类型
	//contentType := "application/zip"

	// 调用 FPutObject 接口上传文件。
	fullPath := filePath + "/" + objectName
	_, err := minioClient.PutObject(bucketName, fullPath, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		logger.GetLogger().Error("上传文件失败", err)
		return err
	}

	logger.GetLogger().Debugf("上传文件 %s 成功\n", objectName)
	logger.GetLogger().Debugf("路径: %s", fullPath)
	return nil
}
