package handdlers

import (
	"../autils"
	"database/sql"
<<<<<<< HEAD
=======
	"fmt"
	"sort"
>>>>>>> 36daf9260d0a72a7220ff43a5ba2529656a16f96
)

var (
	limit = "31"
)

type kv struct {
	Key   string
	Value int
}

func TopSites(db *sql.DB, date string) int {
	topSites := map[string]int{}
	now := autils.ParseTimeStr(date)
	_, last := autils.GetMonthDate(now)
	_, lastMonthDate := autils.GetMonthDate(now.AddDate(0, -1, 0))
	sqlStr := "select domain, pv from site_detail where date = '" + autils.GetCurrentDate(last) + "' order by pv desc offset 0 limit " + limit
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

	firstDateStr := autils.GetCurrentDate(lastMonthDate)
	lastDateStr := autils.GetCurrentDate(last)

	lastTopList := getLastTop(db, firstDateStr, lastDateStr)

	diffList := map[string]int{}

	for i, v := range topSites {
		diffList[i] = v - lastTopList[i]
	}

<<<<<<< HEAD
	topSum := 0
	for _, v := range diffList {
		topSum += v
=======
	var tmpKV []kv
	for k, v := range diffList {
		tmpKV = append(tmpKV, kv{k, v})
>>>>>>> 36daf9260d0a72a7220ff43a5ba2529656a16f96
	}

	sort.Slice(tmpKV, func(i, j int) bool {
		return tmpKV[i].Value >= tmpKV[j].Value
	})
	tmpKV = tmpKV[1:]

	topSum := 0
	for _, v := range tmpKV {
		topSum += v.Value
	}

<<<<<<< HEAD
	// sort.Slice(tmpKV, func(i, j int) bool {
	// 	return tmpKV[i].Value >= tmpKV[j].Value
	// })
	// return tmpKV
	return topSum
=======
	fmt.Println(topSum)
>>>>>>> 36daf9260d0a72a7220ff43a5ba2529656a16f96
}

func getLastTop(db *sql.DB, lastMonth, monthTail string) map[string]int {
	sqlStr := "select domain, pv from site_detail where date = '" + lastMonth + "' and domain in (select domain from site_detail where date = '" + monthTail + "' order by pv desc offset 0 limit " + limit + ")"
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
