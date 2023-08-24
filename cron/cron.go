package cron

import (
	"fmt"
	"strings"
)

// Value defines the value type
type Value string

// IsZero returns true if value is empry
func (v Value) IsZero() bool {
	return string(v) == ""
}

// Cron defines the cron structure
type Cron struct {
	Minute        Value `json:"minute"`
	Hour          Value `json:"hour"`
	DayOfTheMonth Value `json:"dayOfTheMonth"`
	Month         Value `json:"month"`
	DayOfTheWeek  Value `json:"dayOfTheWeek"`
}

// Format returns a formatted line that can be added to a crontab
func (c *Cron) Format(user, cmd string) string {
	if c.IsZero() {
		return ""
	}

	if user == "" {
		return fmt.Sprintf("%s %s", c.String(), cmd)
	}

	return fmt.Sprintf("%s %s %s", c.String(), user, cmd)
}

// IsValid returns true if the cron values are valid
func (c *Cron) IsValid() bool {
	if c.IsZero() {
		return false
	}

	_, err := Parse(c.Format("root", "cmd"))
	return err == nil
}

// IsZero returns true if cron is empry
func (c *Cron) IsZero() bool {
	return c.Minute.IsZero() ||
		c.Hour.IsZero() ||
		c.DayOfTheMonth.IsZero() ||
		c.Month.IsZero() ||
		c.DayOfTheWeek.IsZero()
}

func (c *Cron) String() string {
	if c.IsZero() {
		return ""
	}

	return fmt.Sprintf("%v %v %v %v %v", c.Minute, c.Hour, c.DayOfTheMonth,
		c.Month, c.DayOfTheWeek)
}

// Parse parses the given string and returns its Cron representation.
// Parsing only considers the first 5 elements, ignoring the rest.
func Parse(s string) (*Cron, error) {
	list := strings.Split(s, " ")
	if len(list) < 5 {
		return nil, fmt.Errorf("a cron expression must have at least 5 entries (and a command)")
	}

	if !minuteIsValid(list[0]) {
		return nil, fmt.Errorf("invalid expression for minutes")
	}
	if !hourIsValid(list[1]) {
		return nil, fmt.Errorf("invalid expression for hour")
	}
	if !dayOfTheMonthIsValid(list[2]) {
		return nil, fmt.Errorf("invalid expression for day-of-the-month")
	}
	if !monthIsValid(list[3]) {
		return nil, fmt.Errorf("invalid expression for month")
	}
	if !dayOfTheWeekIsValid(list[4]) {
		return nil, fmt.Errorf("invalid expression for day-of-the-week")
	}

	return &Cron{
			Minute:        Value(list[0]),
			Hour:          Value(list[1]),
			DayOfTheMonth: Value(list[2]),
			Month:         Value(list[3]),
			DayOfTheWeek:  Value(list[4]),
		},
		nil
}
