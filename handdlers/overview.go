package handdlers

import (
	"../autils"
	"database/sql"
)

type orResult struct {
	AllFlow     int
	DomainCount int
	Diff        int
	Rate        float32
	CircleRate  float32
}

func GetOverview(db *sql.DB, date string) orResult {
	now := autils.ParseTimeStr(date)
	_, last := autils.GetMonthDate(now)
	_, lastMonthDate := autils.GetMonthDate(now.AddDate(0, -1, 0))
	_, tMonthTime := autils.GetMonthDate(now.AddDate(0, -2, 0))

	lastDateStr := autils.GetCurrentDate(last)
	lastMDateStr := autils.GetCurrentDate(lastMonthDate)
	tMonthStr := autils.GetCurrentDate(tMonthTime)

	allFlow := getAllFlow(db, lastDateStr)
	dCount := getDCount(db, lastDateStr)
	diff, rate := getRaiseNum(db, lastDateStr, lastMDateStr)

	// 环比
	_, cRate := getRaiseNum(db, lastMDateStr, tMonthStr)

	rs := orResult{}
	rs.AllFlow = allFlow
	rs.DomainCount = dCount
	rs.Diff = diff
	rs.Rate = rate
	rs.CircleRate = cRate

	return rs
}

// 当前流量
func getAllFlow(db *sql.DB, day string) int {
	rows, err := db.Query("select click from all_flow where date = '" + day + "'")
	autils.ErrHadle(err)

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

// 域名总数
func getDCount(db *sql.DB, day string) int {
	rows, err := db.Query("select count(domain) from site_detail where date = '" + day + "'")
	autils.ErrHadle(err)

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

// 增长流量
func getRaiseNum(db *sql.DB, lastDate, newDadte string) (int, float32) {
	rows, err := db.Query("select click from all_flow where date = '" + lastDate + "' or date = '" + newDadte + "' order by ana_date desc")
	autils.ErrHadle(err)

	var nums []int
	var pv int
	for rows.Next() {
		err := rows.Scan(&pv)
		autils.ErrHadle(err)
		nums = append(nums, pv)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	if len(nums) > 1 {
		diff := nums[0] - nums[1]
		rate := float32(diff) / float32(nums[1]) * 100
		return diff, rate
	}

	defer rows.Close()
	return 0, 0.0
}
