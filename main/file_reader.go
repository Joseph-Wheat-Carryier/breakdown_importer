package main

import (
	"CNXM_BRKD_READER/gorm"
	"CNXM_BRKD_READER/logger"
	"CNXM_BRKD_READER/minio_util"
	"bytes"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var bar = new(pb.ProgressBar)

type BreakdownTreeRel struct {
	StructTreeId int `gorm:"<-"`
	BreakdownId  int `gorm:"<-"`
	RootId       int `gorm:"<-"`
	Lft          int `gorm:"<-"`
	Rgt          int `gorm:"<-"`
}

func (b *BreakdownTreeRel) TableName() string {
	return "cnxm_breakdown_tree_rel"
}

func queryBreakdownTreeRel(id int, partId int) BreakdownTreeRel {
	db := gorm.GetDB()

	b := new(BreakdownTreeRel)
	sql := `
		select
    crvp.id as struct_tree_id,
    crvst.root_id,
    crvst.lft,
    crvst.rgt
from cnxm_resume_vehicle_parts crvp
    left join cnxm_resume_vehicle_struct_tree crvst on crvp.struct_id = crvst.id
where crvp.id = ? 
	`
	db.Raw(sql, partId).Scan(b)
	b.BreakdownId = id
	b.StructTreeId = partId
	return *b
}

func ReadFile(path string, filename string) {
	f, err := excelize.OpenFile(path)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 计算行数
	rowCounts := computeRows(f)
	bar.SetTotal(int64(rowCounts))
	bar.SetWriter(os.Stdout)
	bar.Start()

	sheetList := f.GetSheetList()

	for idx, name := range sheetList {
		log.Debug(fmt.Sprintf("去读第%d个sheet:%s", idx, name))
		ReadSheet(name, f)
	}

	bar.SetCurrent(bar.Total())
	bar.Finish()
	logger.GetLogger().Info("导入成功!!!")
}

func ReadSheet(sheetName string, f *excelize.File) {
	db := gorm.GetDB()
	rows, _ := f.Rows(sheetName)
	for idx := 1; rows.Next(); idx++ {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.GetLogger().Errorf("第%d行有误%v,", idx, err)
				}
			}()
			if idx == 1 {
				return
			}
			breakdown := NewBreakDown()

			var readColErr error
			vehicleName, err := ReadCell(sheetName, idx, "C", f)
			if err != nil {
				print2Log(fmt.Sprintf("%d行的数据导入失败, 原因:%s", idx, err.Error()))
				return
			}

			for _, col := range COLS {
				value, err := ReadCell(sheetName, idx, col, f)
				if err != nil {
					continue
				}
				err = breakdown.setValue(col, *value, *vehicleName)
				if err != nil {
					print2Log(fmt.Sprintf("%d行的数据导入失败, 原因:%s", idx, err.Error()))
					readColErr = err
					break
				}
			}

			if readColErr != nil {
				return
			}

			logger.GetLogger().Debug(breakdown.WorkTeam)
			// 设置dataId

			breakdown.DataId = uuid.New().String()
			isRepeat := checkRepeat(*breakdown)
			if isRepeat {
				fmt.Println("重复了: " + breakdown.BreakdownDescribe)
				bar.Increment()
				return
			}

			db.Create(breakdown)

			breakdownId := breakdown.Id
			carIdArr := strings.Split(breakdown.CarId, ",")
			for _, partIdStr := range carIdArr {
				partId, _ := strconv.Atoi(partIdStr)
				treeRel := queryBreakdownTreeRel(breakdownId, partId)
				db.Create(treeRel)
			}

			bar.Increment()
		}()
	}

	err := rows.Close()
	if err != nil {
		logger.GetLogger().Error(err)
	}
}

func checkRepeat(breakdown BreakDown) bool {
	db := gorm.GetDB()
	var happenDate *time.Time
	sql := `select happen_date from cnxm_breakdown where breakdown_describe = ?`

	db.Raw(sql, breakdown.BreakdownDescribe).Scan(&happenDate)
	if happenDate == nil {
		return false
	}
	return happenDate.Equal(*breakdown.HappenDate)
}

func ReadCell(sheet string, row int, col string, file *excelize.File) (*string, error) {
	db := gorm.GetDB()
	cellIdx := fmt.Sprintf("%s%d", col, row)
	if col != "P" && col != "Q" {
		value, err := file.GetCellValue(sheet, cellIdx)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return &value, nil
	} else {
		pics, err := file.GetPictures(sheet, fmt.Sprintf("%s%d", col, row))
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		ossFileIds := make([]string, 0, 10)
		for _, pic := range pics {
			id, err := Next()
			name := *id + pic.Extension
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}

			layout := "20060102"
			// 使用指定格式进行格式化
			currentTime := time.Now()
			formattedTime := currentTime.Format(layout)

			filePath := "/breakdown/file/" + formattedTime + "/"
			err = minio_util.UploadFile(name, filePath, bytes.NewReader(pic.File))
			if err != nil {
				logger.GetLogger().Error("上传图片失败,忽略...,name:" + name)
				return nil, err
			}

			ossFile := NewOssFile(*id, name, "/cnxm"+filePath+"/"+name)
			db.Create(ossFile)
			ossFileIds = append(ossFileIds, *id)
		}

		ossFileIdsStr := strings.Join(ossFileIds, ",")
		return &ossFileIdsStr, nil
	}

}

func computeRows(f *excelize.File) int {
	var count = 0
	for _, sheet := range f.GetSheetList() {
		dimension, _ := f.GetSheetDimension(sheet)
		if !strings.Contains(dimension, ":") {
			continue
		}
		end := strings.Split(dimension, ":")[1]

		pattern := `\d+$`
		regexpPattern := regexp.MustCompile(pattern)
		matches, _ := strconv.Atoi(regexpPattern.FindString(end))
		count += matches
	}

	return count
}

func print2Log(msg string) {
	// 文件路径
	filePath := "log.txt" // 替换为您要写入的文件路径
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	// 打开文件以进行写入，如果文件不存在则创建
	if err != nil {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Println("无法创建文件:", err)
			return
		}
	}

	defer file.Close() // 延迟关闭文件

	fmt.Fprintln(file, msg)
}
