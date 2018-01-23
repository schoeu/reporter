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
	var date, mode string
	rsFilePath := "Result_Markdown.md"
	aType := 1

	flag.StringVar(&date, "date", "", "分析数据的日期, 默认当前日期（yyyy-MM-dd）")
	flag.IntVar(&aType, "type", 1, "分析数据类型")
	flag.StringVar(&mode, "mode", "", "是否为开发模式")
	flag.Parse()

	if date == "" {
		date = autils.GetCurrentDate(time.Now())
	}

	db := autils.OpenDb("postgres", config.PQFlowUrl)
	cwd := autils.GetCwd()
	rsFile := filepath.Join(cwd, rsFilePath)

	if mode != "" {
		// db = autils.OpenDb("postgres", config.PQTestUrl)
		// newSite 当月新增站点pv平均值之和
		// totle 站点当月pv总和的平均值
		// kvs TOP站点净增长总量

		// TODO
		// newSite, totle := mock.NewSite, mock.Total
		// kvs := mock.Kvs
		// handdlers.MarkdownMaker(rsFile)
		return
	}

	if aType == 0 {
		// 概览
		dataList := handdlers.GetOverview(db, date)
		fmt.Println(dataList, dataList.Diff)

	} else if aType == 1 {
		// 新增站点
		newSiteInfo := handdlers.SitePraise(db, date)
		fmt.Println(newSiteInfo, newSiteInfo.Newer)
	} else if aType == 2 {
		// TOP站点计算
		raiseNum := handdlers.TopSites(db, date)
		fmt.Println(raiseNum)
	} else if aType == 2 {
		handdlers.MarkdownMaker(rsFile)
	}

	defer db.Close()
}
