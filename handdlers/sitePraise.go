package handdlers

import (
	"strings"
	"time"
	"database/sql"
	"fmt"
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
	count := []int{}
	firstDateStr := autils.GetCurrentDate(first)
	lastDateStr := autils.GetCurrentDate(last)
	strArr := []string{}
	for _, v := range sites {
		sqlStr := "select avg(pv) from site_detail where domain = '"+v+"' and date >= '"+ firstDateStr +"' and date <= '"+ lastDateStr +"'"
		strArr = append(strArr, sqlStr)
	}
	rows, err := db.Query(strings.Join(strArr, " union all  "))
	var pv int
	for rows.Next() {
		err := rows.Scan(&pv)
		autils.ErrHadle(err)
		count = append(count, pv)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()
	fmt.Println(count)
}