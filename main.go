package main

import (
	"./autils"
	"./config"
	"./handdlers"
	"flag"
	"time"
)

func main() {
	date := ""
	aType := 1

	flag.StringVar(&date, "date", "", "分析数据的日期, 默认当前日期（yyyy-MM-dd）")
	flag.IntVar(&aType, "type", 1, "分析数据类型")
	flag.Parse()

	if date == "" {
		date = autils.GetCurrentDate(time.Now())
	}

	db := autils.OpenDb("postgres", config.PQFlowUrl)
	//db := autils.OpenDb("postgres", config.PQTestUrl)

	if aType == 1 {
		// 增长计算
		handdlers.SitePraise(db, date)
	} else if aType == 2 {
		// TOP站点计算
		handdlers.TopSites(db, date)
	}
	defer db.Close()
}
