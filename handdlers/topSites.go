package handdlers

import (
	"../autils"
	"database/sql"
	"sort"
	"time"
)

var (
	limit = "30"
)

type kv struct {
	Key   string
	Value int
}

type topSitesRs struct {
	Total      int
	Raise      int
	Rate       float32
	TotalRate  float32
	CTotalRate float32
	SigleSite  []sigleSite
}

type sigleSite struct {
	Domain string
	Pv     int
	Rate   float32
}

func TopSites(db *sql.DB, date string) topSitesRs {
	last, lastMonthDate := getLastTime(date)
	firstDateStr := autils.GetCurrentDate(lastMonthDate)
	lastDateStr := autils.GetCurrentDate(last)

	ss := sigleSite{}
	ssCtt := []sigleSite{}

	// 当前TOP流量
	currentTop := getCurrentTop(db, date)
	// TOP流量环比
	//lCurrentTop := getCurrentTop(db, firstDateStr)

	// 新增总流量
	diff, _ := getRaiseNum(db, lastDateStr, firstDateStr)
	// 当月总流量
	allFlow := getAllFlow(db, lastDateStr)
	// 上月总流量
	//lAllFlow := getAllFlow(db, firstDateStr)

	sortedTop := sortMap(currentTop)

	// 获取TOP单站数据
	for _, v := range sortedTop {
		ss.Domain = v.Key
		ss.Pv = v.Value
		ss.Rate = float32(v.Value) / float32(allFlow) * 100
		ssCtt = append(ssCtt, ss)
	}

	// 获取当前TOP占总流量&各站点环比差
	topTotal, diffList := getTopTotal(db, date, currentTop)
	// 获取TOP占总流量&各站点环比差 环比
	//lTopTotal, _ := getTopTotal(lCurrentTop)

	// Map按value排序
	sortedMap := sortMap(diffList)

	// TOP新增流量
	topSum := 0
	for _, v := range sortedMap {
		topSum += v.Value
	}
	tsr := topSitesRs{}
	// TOP总量
	tsr.Total = topTotal
	// TOP新增量
	tsr.Raise = topSum
	// 新增流量占总增长量比例
	tsr.Rate = float32(topSum) / float32(diff) * 100
	// TOP总量占总流量比例
	tsr.TotalRate = float32(topTotal) / float32(allFlow) * 100
	//tsr.CTotalRate = float32(lTopTotal) / float32(lAllFlow) * 100

	// TOP单站数据
	tsr.SigleSite = ssCtt

	return tsr
}

func sortMap(data map[string]int) []kv {
	var tmpKV []kv
	for k, v := range data {
		tmpKV = append(tmpKV, kv{k, v})
	}

	sort.Slice(tmpKV, func(i, j int) bool {
		return tmpKV[i].Value >= tmpKV[j].Value
	})

	return tmpKV
}

func getTopTotal(db *sql.DB, date string, data map[string]int) (int, map[string]int) {
	last, lastMonthDate := getLastTime(date)
	firstDateStr := autils.GetCurrentDate(lastMonthDate)
	lastDateStr := autils.GetCurrentDate(last)

	t := 0
	diffList := map[string]int{}

	// TOP流量环比
	lastTopList := getLastTop(db, firstDateStr, lastDateStr)

	for i, v := range data {
		t += v
		diffList[i] = v - lastTopList[i]
	}
	return t, diffList
}

func getCurrentTop(db *sql.DB, date string) map[string]int {
	topSites := map[string]int{}
	last, _ := getLastTime(date)
	sqlStr := "select domain, pv from site_detail where date = '" + autils.GetCurrentDate(last) + "' order by pv desc offset 1 limit " + limit
	rows, err := db.Query(sqlStr)
	domain := ""
	pv := 0

	for rows.Next() {
		err := rows.Scan(&domain, &pv)
		autils.ErrHadle(err)
		topSites[domain] = pv
	}

	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()
	return topSites
}

func getLastTime(date string) (time.Time, time.Time) {
	now := autils.ParseTimeStr(date)
	_, last := autils.GetMonthDate(now)
	_, lastMonthDate := autils.GetMonthDate(now.AddDate(0, -1, 0))
	return last, lastMonthDate
}

func getLastTop(db *sql.DB, lastMonth, monthTail string) map[string]int {
	sqlStr := "select domain, pv from site_detail where date = '" + lastMonth + "' and domain in (select domain from site_detail where date = '" + monthTail + "' order by pv desc offset 1 limit " + limit + ")"
	rows, err := db.Query(sqlStr)
	domain := ""
	pv := 0
	lastTopSites := map[string]int{}
	for rows.Next() {
		err := rows.Scan(&domain, &pv)
		autils.ErrHadle(err)
		lastTopSites[domain] = pv
	}

	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()
	return lastTopSites
}
