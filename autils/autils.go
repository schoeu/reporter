package autils

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const sqlReg = "(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\\b(select|update|and|or|delete|insert|trancate|char|into|substr|ascii|declare|exec|count|master|into|drop|execute)\\b)"

// 获取当前时间字符串
func GetCurrentDate(date time.Time) string {
	t := date.String()
	return strings.Split(t, " ")[0]
}

type anaChain struct {
	value   string
	content string
}

// 获取程序cwd
func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// 统一错误处理
func ErrHadle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// check sql string
func CheckSql(s string) string {
	match, _ := regexp.Match(sqlReg, []byte(s))
	if match {
		return ""
	}
	return s
}

// 创建数据库链接
func OpenDb(dbTyepe string, dbStr string) *sql.DB {
	if dbTyepe == "" {
		dbTyepe = "mysql"
	}
	db, err := sql.Open(dbTyepe, dbStr)
	ErrHadle(err)

	err = db.Ping()
	ErrHadle(err)
	return db
}

func GetMonthDate(now time.Time) (time.Time, time.Time) {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return firstOfMonth, lastOfMonth
}

func ParseTimeStr(date string) time.Time {
	shortTime := "2006-01-02"
	rsDate, err := time.Parse(shortTime, date)
	ErrHadle(err)
	return rsDate
}

func GetCircleDate(st, et string) (time.Time, time.Time) {
	startTime := ParseTimeStr(st)
	endTime := ParseTimeStr(et)
	sub := endTime.Sub(startTime.AddDate(0, 0, -1))
	lastStart := startTime.Add(-sub)
	lastEnd := startTime.AddDate(0, 0, -1)
	return lastStart, lastEnd
}
