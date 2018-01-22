package handdlers

import (
	"time"
	"database/sql"
	"fmt"
	"../autils"
)

func SitePraise(db *sql.DB, date string) {
	newSites := []string{}
	first, last := autils.GetMonthDate(date)
	sqlStr := "select domain from site_detail where date = '"+ autils.GetCurrentDate(last) +"' except  select domain from site_detail where date = '"+ autils.GetCurrentDate(first) +"'"
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
	siteFlow(db, newSites, first, last)
}

func siteFlow(db *sql.DB, sites []string, first time.Time, last time.Time) {
	count := 0
	for _, v := range sites {
		sqlStr := "select pv from site_detail where date >= '"+ autils.GetCurrentDate(first) +"' date <= '"+ autils.GetCurrentDate(last) +"' and domain = '"+ v + "'"
		rows, err := db.Query(sqlStr)
		var pv int
		for rows.Next() {
			err := rows.Scan(&pv)
			autils.ErrHadle(err)
			count = count + pv
		}
		err = rows.Err()
		autils.ErrHadle(err)
		defer rows.Close()
	}
}