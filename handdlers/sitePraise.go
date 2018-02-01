package handdlers

import (
	"../autils"
	"database/sql"
	"sync"
)

type nSiteInfo struct {
	Newer   int
	NewFlow float32
	CnFlow  float32
}

func SitePraise(db *sql.DB, st, et string) nSiteInfo {
	newSites := []string{}
	sqlStr := "select domain from site_detail where date = '" + et + "' except  select domain from site_detail where date = '" + st + "'"
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

	cs, ce := autils.GetCircleDate(st, et)
	csStr := autils.GetCurrentDate(cs)
	ceStr := autils.GetCurrentDate(ce)

	newSiteFlow := siteFlow(db, newSites, st, et)
	total := getTotalFlow(db, st, et)

	// 环比
	cNewSiteFlow := siteFlow(db, newSites, csStr, ceStr)
	cTotal := getTotalFlow(db, csStr, ceStr)

	nsi := nSiteInfo{}
	nsi.Newer = len(newSites)
	nsi.NewFlow = float32(newSiteFlow) / float32(total) * 100
	nsi.CnFlow = float32(cNewSiteFlow) / float32(cTotal) * 100
	return nsi
}

// 新增站点的平均pv
func siteFlow(db *sql.DB, sites []string, first, last string) int {
	var count int
	var mutex sync.Mutex
	var pv sql.NullInt64
	for _, v := range sites {
		mutex.Lock()
		sqlStr := "select ceil(avg(pv)) from site_detail where domain = '" + v + "' and date >= '" + first + "' and date <= '" + last + "'"
		rows, err := db.Query(sqlStr)
		for rows.Next() {
			err := rows.Scan(&pv)
			autils.ErrHadle(err)
			realPv := int(pv.Int64)
			if realPv != 0 {
				count = count + realPv
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
	rows, err := db.Query(sqlStr)
	var total sql.NullInt64
	for rows.Next() {
		err := rows.Scan(&total)
		autils.ErrHadle(err)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()
	return int(total.Int64)
}
