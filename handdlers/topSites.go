package handdlers

import (
	"database/sql"
	"fmt"
	"sync"
	"../autils"
)

func SitePraise(db *sql.DB, date string) {
	limit := "31"
	topSites := []string{}
	now := autils.ParseTimeStr(date)
	_, last := autils.GetMonthDate(now)
	sqlStr := "select domain, pv from site_detail where date = '"+last+"' order by pv desc offset 0 limit " + limit
	rows, err := db.Query(sqlStr)
	domain := ""
	pv := 0
	var totalNum int

	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)
		topSites = append(topSites, domain)
	}
	
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()

	totalNum := topSites[0]
	topSites = topSites[1:]

	firstDateStr := autils.GetCurrentDate(first)
	lastDateStr := autils.GetCurrentDate(last)

	newSiteFlow := siteFlow(db, newSites, firstDateStr, lastDateStr)
	total := getTotalFlow(db, firstDateStr, lastDateStr)
	fmt.Println(firstDateStr, lastDateStr)
	fmt.Println(newSiteFlow, total)
}
