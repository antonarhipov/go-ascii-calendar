package calendar

import (
	"testing"
	"time"
)

func TestGetMonthName(t *testing.T) {
	tests := []struct {
		name     string
		month    time.Time
		expected string
	}{
		{"January", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), "January"},
		{"February", time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC), "February"},
		{"December", time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC), "December"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMonthName(tt.month)
			if result != tt.expected {
				t.Errorf("GetMonthName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetYear(t *testing.T) {
	tests := []struct {
		name     string
		month    time.Time
		expected string
	}{
		{"2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), "2025"},
		{"2024", time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC), "2024"},
		{"1999", time.Date(1999, time.June, 15, 0, 0, 0, 0, time.UTC), "1999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetYear(tt.month)
			if result != tt.expected {
				t.Errorf("GetYear() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDaysInMonth(t *testing.T) {
	tests := []struct {
		name     string
		month    time.Time
		expected int
	}{
		{"January", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), 31},
		{"February non-leap", time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC), 28},
		{"February leap", time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), 29},
		{"April", time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC), 30},
		{"December", time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC), 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDaysInMonth(tt.month)
			if result != tt.expected {
				t.Errorf("GetDaysInMonth() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		expected bool
	}{
		{"2024 leap", 2024, true},
		{"2025 not leap", 2025, false},
		{"2000 leap (divisible by 400)", 2000, true},
		{"1900 not leap (divisible by 100 but not 400)", 1900, false},
		{"2004 leap", 2004, true},
		{"1999 not leap", 1999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLeapYear(tt.year)
			if result != tt.expected {
				t.Errorf("IsLeapYear(%d) = %v, want %v", tt.year, result, tt.expected)
			}
		})
	}
}

func TestGetWeekday(t *testing.T) {
	tests := []struct {
		name     string
		month    time.Time
		expected int
	}{
		{"January 2025 starts on Wednesday", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), 3},
		{"February 2025 starts on Saturday", time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC), 6},
		{"August 2025 starts on Friday", time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC), 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWeekday(tt.month)
			if result != tt.expected {
				t.Errorf("GetWeekday() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name      string
		dateStr   string
		expectErr bool
		expected  time.Time
	}{
		{"Valid date", "2025-08-16", false, time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC)},
		{"Valid leap day", "2024-02-29", false, time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)},
		{"Invalid format", "25-08-16", true, time.Time{}},
		{"Invalid date", "2025-13-01", true, time.Time{}},
		{"Empty string", "", true, time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.dateStr)

			if tt.expectErr && err == nil {
				t.Errorf("ParseDate(%s) expected error but got none", tt.dateStr)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("ParseDate(%s) unexpected error: %v", tt.dateStr, err)
			}

			if !tt.expectErr && !result.Equal(tt.expected) {
				t.Errorf("ParseDate(%s) = %v, want %v", tt.dateStr, result, tt.expected)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name      string
		timeStr   string
		expectErr bool
		expected  string
	}{
		{"Valid morning time", "09:30", false, "09:30"},
		{"Valid afternoon time", "14:45", false, "14:45"},
		{"Valid midnight", "00:00", false, "00:00"},
		{"Valid end of day", "23:59", false, "23:59"},
		{"Invalid format", "9:30", false, "09:30"},
		{"Invalid hour", "25:30", true, ""},
		{"Invalid minute", "12:60", true, ""},
		{"Empty string", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.timeStr)

			if tt.expectErr && err == nil {
				t.Errorf("ParseTime(%s) expected error but got none", tt.timeStr)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("ParseTime(%s) unexpected error: %v", tt.timeStr, err)
			}

			if !tt.expectErr && result.Format("15:04") != tt.expected {
				t.Errorf("ParseTime(%s) = %v, want %v", tt.timeStr, result.Format("15:04"), tt.expected)
			}
		})
	}
}

func TestValidateTimeString(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		expected bool
	}{
		{"Valid morning", "09:30", true},
		{"Valid afternoon", "14:45", true},
		{"Valid midnight", "00:00", true},
		{"Valid end of day", "23:59", true},
		{"Single digit hour", "9:30", true},
		{"Single digit minute", "09:5", true},
		{"Invalid hour", "25:30", false},
		{"Invalid minute", "12:60", false},
		{"No colon", "0930", false},
		{"Empty string", "", false},
		{"Non-numeric", "ab:cd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateTimeString(tt.timeStr)
			if result != tt.expected {
				t.Errorf("ValidateTimeString(%s) = %v, want %v", tt.timeStr, result, tt.expected)
			}
		})
	}
}

func TestIsToday(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"Today", today, true},
		{"Yesterday", yesterday, false},
		{"Tomorrow", tomorrow, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsToday(tt.date)
			if result != tt.expected {
				t.Errorf("IsToday() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsSameDate(t *testing.T) {
	date1 := time.Date(2025, time.August, 16, 10, 30, 0, 0, time.UTC)
	date2 := time.Date(2025, time.August, 16, 15, 45, 0, 0, time.UTC)
	date3 := time.Date(2025, time.August, 17, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		date1    time.Time
		date2    time.Time
		expected bool
	}{
		{"Same date different time", date1, date2, true},
		{"Different dates", date1, date3, false},
		{"Same date and time", date1, date1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSameDate(tt.date1, tt.date2)
			if result != tt.expected {
				t.Errorf("IsSameDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCalendarWeeks(t *testing.T) {
	// Test August 2025 (starts on Friday, 5th weekday, has 31 days)
	aug2025 := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	weeks := GetCalendarWeeks(aug2025)

	// August 2025 should have 6 weeks (starts on Friday, so first week has 2 days)
	if len(weeks) != 6 {
		t.Errorf("GetCalendarWeeks(Aug 2025) returned %d weeks, want 6", len(weeks))
	}

	// First week should have 0,0,0,0,0,1,2
	firstWeek := weeks[0]
	expected := []int{0, 0, 0, 0, 0, 1, 2}
	for i, day := range firstWeek {
		if day != expected[i] {
			t.Errorf("First week day %d = %d, want %d", i, day, expected[i])
		}
	}

	// Last week should end with 31
	lastWeek := weeks[len(weeks)-1]
	found31 := false
	for _, day := range lastWeek {
		if day == 31 {
			found31 = true
			break
		}
	}
	if !found31 {
		t.Error("Last week should contain day 31")
	}
}

func TestGetDayOfWeekHeaders(t *testing.T) {
	headers := GetDayOfWeekHeaders()
	expected := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}

	if len(headers) != 7 {
		t.Errorf("GetDayOfWeekHeaders() returned %d headers, want 7", len(headers))
	}

	for i, header := range headers {
		if header != expected[i] {
			t.Errorf("Header %d = %s, want %s", i, header, expected[i])
		}
	}
}
