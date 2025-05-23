package main

import (
	"CNXM_BRKD_READER/gorm"
	"CNXM_BRKD_READER/logger"
	"fmt"
	gorm2 "gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type BreakDown struct {
	Id                        int    `gorm:"primaryKey"`
	RecordSheetCode           string `gorm:"<-"`
	TrackId                   int
	Sort                      int
	VehicleId                 int
	CarId                     string
	Kilometres                int
	ApplyDate                 *time.Time
	ReceiveDate               *time.Time
	HappenDate                *time.Time
	BreakdownDescribe         string
	FileId                    string
	RespondTime               *time.Time
	BreakdownAffect           string
	BreakdownLevel            int
	RepairTime                *time.Time
	BreakdownCode             string
	HandleResult              int
	HandleMan                 string
	HandleMethod              int
	BreakdownSource           string
	BreakdownResult           int
	ResultSort                int
	ChangePartId              *int
	ChangeNum                 int
	HandleCondition           string
	FollowCondition           string
	Notes                     string
	ImportantBreakdown        int
	ImportantBreakdownContent string
	BreakdownAffectSort       int
	AffectTrainNum            int
	CreateBy                  string
	CreateTime                *time.Time `gorm:"<-"`
	UpdateBy                  string
	UpdateTime                *time.Time
	DelFlag                   bool
	ResumeVehiclePartId       int
	TypeId                    int
	Classify                  string
	DataId                    string
	BreakdownStatus           int
	BreakdownFormId           string
	BreakdownReportId         string
	BreakdownFormWhere        int
	BreakdownReportWhere      int
	ChangePartOrderCode       string
	AfterProcessing           string
	WorkTeam                  string
}

func (b *BreakDown) TableName() string {
	return "cnxm_breakdown"
}

func NewBreakDown() *BreakDown {
	breakdown := new(BreakDown)
	curTime := time.Now()
	breakdown.UpdateTime = &curTime
	breakdown.CreateTime = &curTime
	breakdown.RecordSheetCode = GenerateSheetCode()
	return breakdown
}

var layout = "2006-01-02"

var COLS = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q"}

func (b *BreakDown) setValue(col string, value string, vehicleName string) error {
	db := gorm.GetDB()

	var err error = nil

	switch col {
	case "A":
		return nil
	case "B":
		b.ApplyDate, err = parseTime(value)
		b.ReceiveDate = b.ApplyDate
		b.HappenDate = b.ApplyDate
		printLog(err)
		return nil
	case "C":
		id := *selectVehicleId(db, value)
		b.VehicleId = id
		return nil
	case "D":
		if value == "" {
			return nil
		}
		b.Kilometres, err = strconv.Atoi(value)
		printLog(err)
		return err
	case "E":
		idArr := selectCarId(db, value, vehicleName)
		if len(idArr) == 0 {
			logger.GetLogger().Errorf("没找到车厢：%s, 请先创建车厢", value)
			return fmt.Errorf("没找到车厢：%s, 请先创建车厢", value)
		}
		b.CarId = strings.Join(idArr, ",")
		return nil
	case "F":
		b.BreakdownDescribe = value
		return nil
	case "G":
		b.BreakdownSource = value
		return nil
	case "H":
		switch value {
		case "无":
			b.BreakdownAffect = "0"
			break
		case "故障扣修":
			b.BreakdownAffect = "1"
			break
		case "退出服务":
			b.BreakdownAffect = "2"
			break
		case "晚点2分钟以下":
			b.BreakdownAffect = "3"
			break
		case "晚点2-5分钟":
			b.BreakdownAffect = "4"
			break
		case "晚点5分钟以上":
			b.BreakdownAffect = "5"
			break
		case "清客":
			b.BreakdownAffect = "6"
			break
		case "正线救援":
			b.BreakdownAffect = "7"
			break
		case "下线":
			b.BreakdownAffect = "8"
			break
		default:
			//log.info("------" + data.get(i).get("breakdownAffect") + "------");
			logger.GetLogger().Errorf("错误的影响范围: %s", value)
			b.BreakdownAffect = "-1000"
		}
		return nil
	case "I":
		var id = new(int)
		db.Raw("select id from cnxm_resume_vehicle_struct_type tp where tp.type_name = ?", value).Scan(id)
		if id == nil {
			logger.GetLogger().Errorf("未找到功能构型%s", value)
			return fmt.Errorf("未找到功能构型%s", value)
		}
		b.TypeId = *id
		b.Classify = value
		return nil
	case "J":
		b.HandleCondition = value
		return nil
	case "K":
		if value == "未处理" || value == "待处理" {
			b.BreakdownStatus = 0
		}

		if value == "已处理" || value == "已关闭" {
			b.BreakdownStatus = 1
		}

		if strings.Contains(value, "跟踪") {
			b.BreakdownStatus = 2
		}
		return nil
	case "L":
		b.RepairTime, err = parseTime(value)
		return err
	case "M":
		b.HandleMan = value
		return nil
	case "N":
		b.CreateBy = value
		return nil
	case "O":
		id := new(string)
		db.Raw("select id from sys_depart where depart_name= ?", value).Scan(id)
		b.WorkTeam = *id
		if id == nil {
			return fmt.Errorf("未找到车厢%s", value)
		}
		return nil
	case "P":
		b.FileId = value
		return nil
	case "Q":
		b.AfterProcessing = value
		return nil
	}

	return nil
}

func printLog(err error) {
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
}

func parseTime(str string) (*time.Time, error) {
	str = strings.TrimSpace(str)
	timeStr := new(string)

	parseFormat := func(splitter string) {
		arr := strings.Split(str, splitter)
		arr2 := make([]string, 0, 3)
		for _, s := range arr {
			s = fmt.Sprintf("%02s", s)
			arr2 = append(arr2, s)
		}
		str := strings.Join(arr2, splitter)
		timeStr = &str
		logger.GetLogger().Debug(*timeStr)
	}

	parseFormat("-")
	timeParsed, err := time.Parse(layout, *timeStr)
	if err != nil {
		logger.GetLogger().Debug("尝试使用令一个种时间格式")
		parseFormat("/")
		timeParsed, err = time.Parse("2006/01/02", *timeStr)
		if err != nil {
			logger.GetLogger().Error(err)
		}
	}
	return &timeParsed, err
}

func selectVehicleId(db *gorm2.DB, vehicleCode string) *int {
	id := new(int)
	sql := `select id from cnxm_resume_vehicle where vehicle_code = ?`
	db.Raw(sql, vehicleCode).Scan(id)
	return id
}

func selectCarId(db *gorm2.DB, carName string, vehicleName string) (id []string) {
	carName = strings.TrimSpace(carName)
	if carName == "全车" {
		return selectQuanChe(db, vehicleName)
	}
	id = make([]string, 0, 2)

	var carArr []string

	if strings.Contains(carName, "/") {
		carArr = strings.Split(carName, "/")
	} else if strings.Contains(carName, ",") {
		carArr = strings.Split(carName, ",")
	} else if strings.Contains(carName, "，") {
		carArr = strings.Split(carName, "，")
	} else if strings.Contains(carName, "、") {
		carArr = strings.Split(carName, "、")
	} else if len(carName) == 6 && (!strings.Contains(carName, "A") || !strings.Contains(carName, "B")) {
		carArr = []string{
			carName[0:3],
			carName[3:6],
		}
	} else if len(carName) == 8 && (strings.Contains(carName, "A") || strings.Contains(carName, "B")) {
		carArr = []string{
			carName[0:4],
			carName[4:8],
		}
	} else {
		carArr = strings.Split(carName, ",")
	}

	formatCar := func(carStr string) string {
		return fmt.Sprintf("%04s", carStr)
	}

	transformed := make([]string, 0, 10)
	for _, car := range carArr {
		if !endsWith(car, "A") && !endsWith(car, "B") {
			transformed = append(transformed, formatCar(car+"A"), formatCar(car+"B"))
		} else {
			transformed = append(transformed, formatCar(car))
		}
	}

	sql := `select crvp.id
		from cnxm_resume_vehicle_parts crvp
		left join cnxm_resume_vehicle_struct_tree crvst on crvp.struct_id = crvst.id
		left join cnxm_resume_car_type car on car.id = crvst.car_type_id
		left join cnxm_resume_vehicle vehicle on vehicle.id = crvp.vehicle_id
  		where crvst.depth = 1
		and crvp.car_serial_number in ?
		and vehicle.vehicle_name = ?`
	db.Raw(sql, transformed, vehicleName).Scan(&id)
	return id
}

func selectQuanChe(db *gorm2.DB, vehicleName string) []string {
	id := make([]string, 0, 2)
	sql := `select crvp.id
		from cnxm_resume_vehicle_parts crvp
		left join cnxm_resume_vehicle_struct_tree crvst on crvp.struct_id = crvst.id
		left join cnxm_resume_car_type car on car.id = crvst.car_type_id
		left join cnxm_resume_vehicle vehicle on vehicle.id = crvp.vehicle_id
  		where crvst.depth = 1
		and vehicle.vehicle_name = ?`
	db.Raw(sql, vehicleName).Scan(&id)
	return id
}

func endsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}
