package cron

import (
	"slices"
	"strconv"
	"strings"
)

type cronExpr int

const (
	cronExprUnknown cronExpr = iota
	cronExprStar
	cronExprStarNth
	cronExprStarNthMinute
	cronExprStarNthHour
	cronExprStarNthMonth
	cronExprDigit
	cronExprMinute
	cronExprHour
	cronExprDay
	cronExprMonth
	cronExprWeek
	cronExprMonthName
	cronExprWeekName
	cronExprMinuteRange
	cronExprHourRange
	cronExprDayRange
	cronExprMonthRange
	cronExprWeekRange
	cronExprMontNamehRange
	cronExprWeekNameRange
	cronExprMinuteRangeNth
	cronExprHourRangeNth
	cronExprMonthRangeNth
	cronExprMonthNameRangeNth
)

type cronHint int

const (
	cronHintMinute cronHint = iota
	cronHintHour
	cronHintDay
	cronHintMonth
	cronHintWeek
)

var (
	namesOfDays = map[string]int{
		"sunday":    0,
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
		"sun":       0,
		"mon":       1,
		"tue":       2,
		"wed":       3,
		"thu":       4,
		"fri":       5,
		"sat":       6,
	}

	namesOfMonths = map[string]int{
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"may":       5,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
		"jan":       1,
		"feb":       2,
		"mar":       3,
		"apr":       4,
		"jun":       6,
		"jul":       7,
		"aug":       8,
		"sep":       9,
		"oct":       10,
		"nov":       11,
		"dec":       12,
	}
)

func classifyCronExpr(e string, h cronHint) cronExpr {
	if e == "*" {
		return cronExprStar
	}

	digitRangeReturn := func(lr []string, ret cronExpr) cronExpr {
		lv, err := strconv.ParseInt(lr[0], 10, 64)
		if err != nil {
			return cronExprUnknown
		}

		rv, err := strconv.ParseInt(lr[1], 10, 64)
		if err != nil {
			return cronExprUnknown
		}

		if lv >= rv {
			return cronExprUnknown
		}

		return ret
	}

	if strings.Contains(e, "/") {
		list := strings.Split(e, "/")

		if len(list) != 2 {
			return cronExprUnknown
		}

		above := classifyCronExpr(list[0], h)
		below := classifyCronExpr(list[1], h)

		if above == cronExprStar && below == cronExprDigit {
			return cronExprStarNth
		}

		if h == cronHintMinute && slices.Contains([]cronExpr{cronExprStar, cronExprMinuteRange}, above) &&
			slices.Contains([]cronExpr{cronExprMinute, cronExprDigit}, below) &&
			slices.Contains([]string{"2", "3", "4", "5", "6", "10", "12", "15", "20", "30"}, list[1]) {
			return cronExprMinuteRangeNth
		}

		if h == cronHintHour && slices.Contains([]cronExpr{cronExprStar, cronExprHourRange}, above) &&
			slices.Contains([]cronExpr{cronExprHour, cronExprDigit}, below) &&
			slices.Contains([]string{"2", "3", "4", "6", "8", "12"}, list[1]) {
			return cronExprHourRangeNth
		}

		if h == cronHintMonth && slices.Contains([]cronExpr{cronExprStar, cronExprMonthRange, cronExprMontNamehRange}, above) &&
			slices.Contains([]cronExpr{cronExprMonth, cronExprDigit}, below) &&
			slices.Contains([]string{"2", "3", "4", "6"}, list[1]) {
			return cronExprMonthRangeNth
		}

		return cronExprUnknown
	}

	if strings.Contains(e, "-") {
		list := strings.Split(e, "-")

		if len(list) != 2 {
			return cronExprUnknown
		}

		left := classifyCronExpr(list[0], h)
		right := classifyCronExpr(list[1], h)

		if left == cronExprMinute && right == cronExprMinute {
			return digitRangeReturn(list, cronExprMinuteRange)
		}

		if left == cronExprHour && right == cronExprHour {
			return digitRangeReturn(list, cronExprHourRange)
		}

		if left == cronExprDay && right == cronExprDay {
			return digitRangeReturn(list, cronExprDayRange)
		}

		if left == cronExprMonth && right == cronExprMonth {
			return digitRangeReturn(list, cronExprMonthRange)
		}

		if left == cronExprWeek && right == cronExprWeek {
			return digitRangeReturn(list, cronExprWeekRange)
		}

		if left == cronExprMonthName && right == cronExprMonthName {
			if namesOfMonths[list[0]] >= namesOfMonths[list[1]] {
				return cronExprUnknown
			}

			return cronExprMontNamehRange
		}

		if left == cronExprWeekName && right == cronExprWeekName {
			if namesOfDays[list[0]] >= namesOfDays[list[1]] {
				return cronExprUnknown
			}

			return cronExprWeekNameRange
		}

		return cronExprUnknown
	}

	if h == cronHintWeek {
		if _, ok := namesOfDays[strings.ToLower(e)]; ok {
			return cronExprWeekName
		}
	}

	if h == cronHintMonth {
		if _, ok := namesOfMonths[strings.ToLower(e)]; ok {
			return cronExprMonthName
		}
	}

	if v, err := strconv.ParseInt(e, 10, 64); err == nil {
		switch h {
		case cronHintMinute:
			if v >= 0 && v <= 59 {
				return cronExprMinute
			}
		case cronHintHour:
			if v >= 0 && v <= 59 {
				return cronExprHour
			}
		case cronHintDay:
			if v >= 1 && v <= 31 {
				return cronExprDay
			}
		case cronHintMonth:
			if v >= 1 && v <= 12 {
				return cronExprMonth
			}
		case cronHintWeek:
			if v >= 0 && v <= 7 {
				return cronExprWeek
			}
		default:
		}
		return cronExprDigit
	}

	return cronExprUnknown
}

func expressionIsValid(s string, h cronHint, validSet []cronExpr) bool {
	if strings.Contains(s, ",") {
		list := strings.Split(s, ",")
		for _, l := range list {
			if !slices.Contains(validSet, classifyCronExpr(l, h)) {
				return false
			}
		}
		return true
	}

	return slices.Contains(validSet, classifyCronExpr(s, h))
}

func minuteIsValid(s string) bool {
	validSet := []cronExpr{cronExprStar, cronExprStarNthMinute, cronExprMinute,
		cronExprMinuteRange, cronExprMinuteRangeNth}
	return expressionIsValid(s, cronHintMinute, validSet)
}

func hourIsValid(s string) bool {
	validSet := []cronExpr{cronExprStar, cronExprStarNthHour, cronExprHour,
		cronExprHourRange, cronExprHourRangeNth}
	return expressionIsValid(s, cronHintHour, validSet)
}

func dayOfTheMonthIsValid(s string) bool {
	validSet := []cronExpr{cronExprStar, cronExprDay, cronExprDayRange}
	return expressionIsValid(s, cronHintDay, validSet)
}

func monthIsValid(s string) bool {
	validSet := []cronExpr{cronExprStar, cronExprStarNthMonth, cronExprMonth,
		cronExprMonthName, cronExprMonthRange, cronExprMonthRangeNth, cronExprMonthNameRangeNth}
	return expressionIsValid(s, cronHintMonth, validSet)
}

func dayOfTheWeekIsValid(s string) bool {
	validSet := []cronExpr{cronExprStar, cronExprWeek,
		cronExprWeekRange, cronExprWeekNameRange}
	return expressionIsValid(s, cronHintWeek, validSet)
}
