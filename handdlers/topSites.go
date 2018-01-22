package handdlers

import (
	"../autils"
	"database/sql"
	"fmt"
)

type siteInfo struct {
	domain string
	pv     int
}

var (
	limit = "31"
)

func TopSites(db *sql.DB, date string) {
	topSites := []siteInfo{}
	si := siteInfo{}
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
		si.domain = domain
		si.pv = pv
		topSites = append(topSites, si)
	}

	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()

	firstDateStr := autils.GetCurrentDate(lastMonthDate)
	lastDateStr := autils.GetCurrentDate(last)

	lastTopList := getLastTop(db, firstDateStr, lastDateStr)

	diffList := []int{}
	for i, v := range topSites {
		if v.domain == lastTopList[i].domain {
			diffList = append(diffList, v.pv-lastTopList[i].pv)
		}
	}

	fmt.Println(diffList)
}

func getLastTop(db *sql.DB, lastMonth, monthTail string) []siteInfo {
	sqlStr := "select domain, pv from site_detail where date = '" + lastMonth + "' and domain in (select domain from site_detail where date = '" + monthTail + "' order by pv desc offset 0 limit " + limit + ")"
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)
	domain := ""
	pv := 0
	lastTopSites := []siteInfo{}
	lastSi := siteInfo{}

	for rows.Next() {
		err := rows.Scan(&domain, &pv)
		autils.ErrHadle(err)
		lastSi.domain = domain
		lastSi.pv = pv
		lastTopSites = append(lastTopSites, lastSi)
	}

	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()
	return lastTopSites
}
