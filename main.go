package main

import (
	"fmt"
	"flag"
	"time"
	"path/filepath"
	"./autils"
	"./config"
	"./handdlers"
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

	if aType == 1 {
		// 增长计算
		newSite, totle := handdlers.SitePraise(db, date)
		fmt.Println(newSite, totle)
	} else if aType == 2 {
		// TOP站点计算
		raiseNum := handdlers.TopSites(db, date)
		fmt.Println(raiseNum)
	}
	
	handdlers.MarkdownMaker(rsFile)

	defer db.Close()
}
