package handdlers

import (
	"../autils"
	"database/sql"
)

func GetOverview(db *sql.DB, date string) (int, int) {
	allFlowCh := make(chan int)
	dCountCh := make(chan int)
	go getAllFlow(db, allFlowCh, date)
	go getDCount(db, dCountCh, date)
	allFlow := <-allFlowCh
	dCount := <-dCountCh
	return allFlow, dCount
}

// 当前流量
func getAllFlow(db *sql.DB, ch chan int, day string) {
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

	ch <- total
}

// 域名总数
func getDCount(db *sql.DB, ch chan int, day string) {
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

	ch <- total
}
