package main

import (
	"CNXM_BRKD_READER/gorm"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateSheetCode() string {
	db := gorm.GetDB()

	breakdown := new(BreakDown)
	db.Raw("select create_time,record_sheet_code from cnxm_breakdown order by create_time desc,record_sheet_code desc").Scan(breakdown)

	var sc int
	if breakdown.RecordSheetCode == "" {
		sc = 0
	} else {
		sc, _ = strconv.Atoi(breakdown.RecordSheetCode)
		sc = sc % 10000
	}

	if breakdown.CreateTime == nil {
		breakdown.CreateTime = new(time.Time)
		*breakdown.CreateTime = time.Now()
	}
	timeStr := strings.Replace(breakdown.CreateTime.Format("2006-01-02"), "-", "", -1)

	newSc := sc + 1
	return fmt.Sprintf("%s%04d", timeStr, newSc)
}
