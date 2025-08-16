package calendar

import (
	"strconv"
	"strings"
	"time"
)

// GetMonthName returns the full name of the month
func GetMonthName(month time.Time) string {
	return month.Format("January")
}

// GetYear returns the year as a string
func GetYear(month time.Time) string {
	return month.Format("2006")
}

// GetDaysInMonth returns the number of days in the given month
func GetDaysInMonth(month time.Time) int {
	// Get the first day of the next month, then subtract one day
	firstOfNextMonth := time.Date(month.Year(), month.Month()+1, 1, 0, 0, 0, 0, month.Location())
	lastOfThisMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfThisMonth.Day()
}

// GetFirstDayOfMonth returns the first day of the given month
func GetFirstDayOfMonth(month time.Time) time.Time {
	return time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
}

// GetLastDayOfMonth returns the last day of the given month
func GetLastDayOfMonth(month time.Time) time.Time {
	firstDayNextMonth := time.Date(month.Year(), month.Month()+1, 1, 0, 0, 0, 0, month.Location())
	return firstDayNextMonth.AddDate(0, 0, -1)
}

// GetWeekday returns the weekday (0=Sunday, 1=Monday, etc.) for the first day of the month
func GetWeekday(month time.Time) int {
	firstDay := GetFirstDayOfMonth(month)
	return int(firstDay.Weekday())
}

// IsLeapYear checks if the given year is a leap year
func IsLeapYear(year int) bool {
	// A leap year is divisible by 4, but not by 100, unless also divisible by 400
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// GetWeekOfYear returns the ISO week number for a given date
func GetWeekOfYear(date time.Time) int {
	_, week := date.ISOWeek()
	return week
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ParseTime parses a time string in HH:MM format
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// FormatDate formats a date as YYYY-MM-DD
func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

// FormatTime formats a time as HH:MM
func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

// ValidateTimeString validates that a time string is in HH:MM format and valid
func ValidateTimeString(timeStr string) bool {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}

	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return false
	}

	return hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59
}

// GetDayOfWeekHeaders returns the day-of-week headers (Sunday first)
func GetDayOfWeekHeaders() []string {
	return []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
}

// IsToday checks if the given date is today
func IsToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day()
}

// IsSameDate checks if two dates represent the same day (ignoring time)
func IsSameDate(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.Month() == date2.Month() && date1.Day() == date2.Day()
}

// GetCalendarWeeks returns the weeks needed to display a month's calendar
// Each week is represented as an array of day numbers (0 for empty cells)
func GetCalendarWeeks(month time.Time) [][]int {
	firstDay := GetFirstDayOfMonth(month)
	daysInMonth := GetDaysInMonth(month)
	startWeekday := int(firstDay.Weekday()) // 0=Sunday, 1=Monday, etc.

	weeks := [][]int{}
	currentWeek := make([]int, 7)

	// Fill in the empty days before the first day of the month
	for i := 0; i < startWeekday; i++ {
		currentWeek[i] = 0
	}

	// Fill in the days of the month
	dayNum := 1
	for dayNum <= daysInMonth {
		for weekDay := startWeekday; weekDay < 7 && dayNum <= daysInMonth; weekDay++ {
			currentWeek[weekDay] = dayNum
			dayNum++
		}
		weeks = append(weeks, currentWeek)
		currentWeek = make([]int, 7)
		startWeekday = 0 // Reset for subsequent weeks
	}

	return weeks
}
