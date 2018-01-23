package handdlers

import (
	"../autils"
	"database/sql"
	"fmt"
	"sync"
)

type nSiteInfo struct {
	Newer   int
	NewFlow float32
	CnFlow  float32
}

func SitePraise(db *sql.DB, date string) nSiteInfo {
	newSites := []string{}
	now := autils.ParseTimeStr(date)
	first, last := autils.GetMonthDate(now)
	_, lastMonthDate := autils.GetMonthDate(now.AddDate(0, -1, 0))
	sqlStr := "select domain from site_detail where date = '" + autils.GetCurrentDate(last) + "' except  select domain from site_detail where date = '" + autils.GetCurrentDate(lastMonthDate) + "'"
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)
	domain := ""
	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)
		newSites = append(newSites, domain)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()

	firstDateStr := autils.GetCurrentDate(first)
	lastDateStr := autils.GetCurrentDate(last)
	tMDateStr := autils.GetCurrentDate(lastMonthDate)

	newSiteFlow := siteFlow(db, newSites, firstDateStr, lastDateStr)
	total := getTotalFlow(db, firstDateStr, lastDateStr)

	// 环比
	cNewSiteFlow := siteFlow(db, newSites, tMDateStr, firstDateStr)
	cTotal := getTotalFlow(db, tMDateStr, firstDateStr)

	nsi := nSiteInfo{}
	nsi.Newer = len(newSites)
	nsi.NewFlow = float32(newSiteFlow) / float32(total)
	nsi.CnFlow = float32(cNewSiteFlow) / float32(cTotal)
	return nsi
}

func siteFlow(db *sql.DB, sites []string, first, last string) int {
	var count int
	var mutex sync.Mutex
	for _, v := range sites {
		mutex.Lock()
		sqlStr := "select ceil(avg(pv)) from site_detail where domain = '" + v + "' and date >= '" + first + "' and date <= '" + last + "'"
		fmt.Println(sqlStr)
		rows, err := db.Query(sqlStr)
		var pv int
		for rows.Next() {
			err := rows.Scan(&pv)
			autils.ErrHadle(err)
			if pv != 0 {
				count = count + pv
			}
		}
		err = rows.Err()
		autils.ErrHadle(err)
		defer rows.Close()
		mutex.Unlock()
	}
	return count
}

func getTotalFlow(db *sql.DB, first, last string) int {
	sqlStr := "select ceil(avg(pv)) from site_detail where domain = '总和' and date >= '" + first + "' and date <= '" + last + "'"
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)
	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		autils.ErrHadle(err)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()
	return total
}
