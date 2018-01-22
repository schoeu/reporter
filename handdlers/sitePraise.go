package handdlers

import (
	"time"
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
	fmt.Println(len(newSites))
	siteFlow(db, newSites, first, last)
}

func siteFlow(db *sql.DB, sites []string, first time.Time, last time.Time) {
	var count int
	firstDateStr := autils.GetCurrentDate(first)
	lastDateStr := autils.GetCurrentDate(last)
	var mutex sync.Mutex

	for _, v := range sites {
		mutex.Lock()
		sqlStr := "select ceil(avg(pv)) from site_detail where domain = '"+v+"' and date >= '"+ firstDateStr +"' and date <= '"+ lastDateStr +"'"
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
}