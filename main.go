package main

import (
	"./autils"
	"./config"
	"./handdlers"
	"flag"
	"fmt"
	"path/filepath"
	"time"
)

func main() {
	var date, mode, starttime, endtime string
	rsFilePath := "Result_Markdown.md"
	aType := -1

	flag.StringVar(&date, "date", "", "分析数据的日期, 默认当前日期（yyyy-MM-dd）")
	flag.StringVar(&starttime, "s", "", "分析数据的起始日期, 默认一个月前（yyyy-MM-dd）")
	flag.StringVar(&endtime, "e", "", "分析数据的结束日期, 默认当前日期（yyyy-MM-dd）")
	flag.IntVar(&aType, "type", -1, "分析数据类型")
	flag.StringVar(&mode, "debug", "", "是否为开发模式")
	flag.Parse()

	now := time.Now()

	if date == "" {
		date = autils.GetCurrentDate(now)
	}

	if starttime == "" {
		starttime = autils.GetCurrentDate(now.AddDate(0, -1, 0))
	}

	if endtime == "" {
		endtime = autils.GetCurrentDate(now)
	}

	dbUrl := config.PQFlowUrl

	cwd := autils.GetCwd()
	rsFile := filepath.Join(cwd, rsFilePath)

	if mode != "" {
		dbUrl = config.PQTestUrl
		// db = autils.OpenDb("postgres", config.PQTestUrl)
		// newSite 当月新增站点pv平均值之和
		// totle 站点当月pv总和的平均值
		// kvs TOP站点净增长总量

		// TODO
		// newSite, totle := mock.NewSite, mock.Total
		// kvs := mock.Kvs
		// handdlers.MarkdownMaker(rsFile)
	}

	db := autils.OpenDb("postgres", dbUrl)

	if aType == 0 || aType < 0 {
		// 概览
		dataList := handdlers.GetOverview(db, starttime, endtime)
		fmt.Println(dataList, dataList.Diff)

	} else if aType == 1 || aType < 0 {
		// 新增站点
		newSiteInfo := handdlers.SitePraise(db, date)
		fmt.Println(newSiteInfo, newSiteInfo.Newer)
	} else if aType == 2 || aType < 0 {
		// TOP站点计算
		raiseNum := handdlers.TopSites(db, date)
		fmt.Println(raiseNum)
	} else if aType == 3 || aType < 0 {
		handdlers.MarkdownMaker(rsFile)
	} else if aType == 4 || aType < 0 {
		a, b := autils.GetCircleDate("2017-10-10", "2017-12-10")
		fmt.Println(autils.GetCurrentDate(a), autils.GetCurrentDate(b))
	}

	defer db.Close()
}
