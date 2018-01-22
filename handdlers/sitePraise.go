package handdlers

import (
	"database/sql"
	"fmt"
	"sync"
	"../autils"
)

func SitePraise(db *sql.DB, date string) {
	newSites := []string{}
	now := autils.ParseTimeStr(date)
	first, last := autils.GetMonthDate(now)
	_, lastMonthDate := autils.GetMonthDate(now.AddDate(0,-1,0))
	fmt.Println(lastMonthDate, last)
	sqlStr := "select domain from site_detail where date = '"+ autils.GetCurrentDate(last) +"' except  select domain from site_detail where date = '"+ autils.GetCurrentDate(lastMonthDate) +"'"
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

	newSiteFlow := siteFlow(db, newSites, firstDateStr, lastDateStr)
	total := getTotalFlow(db, firstDateStr, lastDateStr)

	fmt.Println(newSiteFlow, total)
}

func siteFlow(db *sql.DB, sites []string, first, last string) int {
	var count int
	
	var mutex sync.Mutex

	for _, v := range sites {
		mutex.Lock()
		sqlStr := "select ceil(avg(pv)) from site_detail where domain = '"+v+"' and date >= '"+ first +"' and date <= '"+ last +"'"
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
		fmt.Println(count)
		mutex.Unlock()
	}
	return count
}

func getTotalFlow(db *sql.DB, first, last string) int {
	sqlStr := "select ceil(avg(pv)) from site_detail where domain = '总和' and date >= '"+ first +"' and date <= '"+ last +"'"
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