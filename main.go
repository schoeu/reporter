package main

import (
	"./autils"
	"./config"
	"./handdlers"
	"time"
	"flag"
)

func main() {
	date := ""
	flag.StringVar(&date, "date", "", "分析数据的日期, 默认当前日期（yyyy-MM-dd）")
	flag.Parse()

	if date == "" {
		date = autils.GetCurrentDate(time.Now())
	}

	//db := autils.OpenDb("postgres", config.PQFlowUrl)
	db := autils.OpenDb("postgres", config.PQTestUrl)
	defer db.Close()
	// 增长计算
	handdlers.SitePraise(db, date)
}