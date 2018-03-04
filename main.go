package main

import (
	"./autils"
	"./config"
	"./handdlers"
	"flag"
	"fmt"
	"time"
)

func main() {
	var mode, starttime, endtime string
	// rsFilePath := "Result_Markdown.md"
	aType := -1

	flag.StringVar(&starttime, "start", "", "分析数据的起始日期, 默认一个月前（yyyy-MM-dd）")
	flag.StringVar(&endtime, "end", "", "分析数据的结束日期, 默认当前日期（yyyy-MM-dd）")
	flag.IntVar(&aType, "type", -1, "分析数据类型")
	flag.StringVar(&mode, "debug", "", "是否为开发模式")
	flag.Parse()

	now := time.Now()

	if starttime == "" {
		starttime = autils.GetCurrentDate(now.AddDate(0, -1, 0))
	}

	if endtime == "" {
		endtime = autils.GetCurrentDate(now)
	}

	dbUrl := config.PQFlowUrl

	// cwd := autils.GetCwd()
	// rsFile := filepath.Join(cwd, rsFilePath)

	if mode != "" {
		dbUrl = config.PQTestUrl
	}

	db := autils.OpenDb("postgres", dbUrl)

	if aType == 0 || aType < 0 {
		// 概览
		dataList := handdlers.GetOverview(db, starttime, endtime)
		fmt.Println(dataList, dataList.Diff)

	} else if aType == 1 || aType < 0 {
		// 新增站点
		newSiteInfo := handdlers.SitePraise(db, starttime, endtime)
		fmt.Println(newSiteInfo)
	} else if aType == 2 || aType < 0 {
		// TOP站点计算
		raiseNum := handdlers.TopSites(db, starttime, endtime)
		fmt.Println(raiseNum)
	}

	defer db.Close()
}
