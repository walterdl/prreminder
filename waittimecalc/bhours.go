package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type BusinessHours struct {
	tz            *time.Location
	weekdays      []time.Weekday
	start         ClockTime
	prWaitMinutes int
}
type ClockTime struct {
	hour int
	min  int
}

func newBusinessHours() (BusinessHours, error) {
	tz, err := timezone()
	if err != nil {
		return BusinessHours{}, err
	}
	wdays, err := weekdays()
	if err != nil {
		return BusinessHours{}, err
	}
	start, err := clockTime()
	if err != nil {
		return BusinessHours{}, err
	}
	prWaitMins, err := prWaitMinutes()
	if err != nil {
		return BusinessHours{}, err
	}

	return BusinessHours{tz, wdays, start, prWaitMins}, nil
}

func timezone() (*time.Location, error) {
	if os.Getenv("TIMEZONE") == "" {
		return nil, errors.New("TIMEZONE environment variable not set")
	}

	return time.LoadLocation(os.Getenv("TIMEZONE"))
}

func weekdays() ([]time.Weekday, error) {
	if os.Getenv("DAYS") == "" {
		// Default to Monday through Friday
		return []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}, nil
	}

	rawDays := strings.Split(os.Getenv("DAYS"), ",")
	result := make([]time.Weekday, len(rawDays))
	for i, raw := range rawDays {
		val, err := strconv.Atoi(raw)
		if err != nil || isInvalidWeekday(val) {
			return nil, fmt.Errorf(("invalid day value: %s"), raw)
		}
		result[i] = time.Weekday(val)
	}

	return result, nil
}

func isInvalidWeekday(val int) bool {
	return val < 0 || val > 6
}

func clockTime() (ClockTime, error) {
	if os.Getenv("START_TIME") == "" {
		return ClockTime{}, errors.New("START_TIME environment variable not set")
	}

	parts := strings.Split(os.Getenv("START_TIME"), ":")
	if len(parts) != 2 {
		return ClockTime{}, errors.New("invalid START_TIME format")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || isInvalidHour(hour) {
		return ClockTime{}, errors.New("invalid START_TIME format")
	}

	min, err := strconv.Atoi(parts[1])
	if err != nil || isInvalidHour(min) {
		return ClockTime{}, errors.New("invalid START_TIME format")
	}

	return ClockTime{hour, min}, nil
}

func isInvalidHour(hour int) bool {
	return hour < 0 || hour > 23
}

func prWaitMinutes() (int, error) {
	val, err := strconv.Atoi(os.Getenv("PR_APPROVAL_WAIT_MINUTES"))

	if err != nil {
		return 0, errors.New("invalid PR_APPROVAL_WAIT_MINUTES value")
	}
	if val < 0 {
		return 0, errors.New("PR_APPROVAL_WAIT_MINUTES must be greater than 0")
	}

	return val, nil
}
