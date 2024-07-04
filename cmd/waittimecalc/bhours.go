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

var errInvalidStartTime = errors.New("invalid START_TIME format")

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

// weekdays returns the weekdays that are considered business days.
// If DAYS is not set, it defaults to Monday through Friday.
// Returns an error if DAYS is invalid.
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

// clockTime returns the start time of the business day.
// If START_TIME is not set, it returns an error.
// Returns an error if START_TIME is invalid.
func clockTime() (ClockTime, error) {
	if os.Getenv("START_TIME") == "" {
		return ClockTime{}, errors.New("START_TIME environment variable not set")
	}

	parts := strings.Split(os.Getenv("START_TIME"), ":")
	if len(parts) != 2 {
		return ClockTime{}, errInvalidStartTime
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || isInvalidHour(hour) {
		return ClockTime{}, errInvalidStartTime
	}

	min, err := strconv.Atoi(parts[1])
	if err != nil || isInvalidHour(min) {
		return ClockTime{}, errInvalidStartTime
	}

	return ClockTime{hour, min}, nil
}

func isInvalidHour(hour int) bool {
	return hour < 0 || hour > 23
}

// prWaitMinutes returns the number of minutes to wait for PR approval.
// If PR_APPROVAL_WAIT_MINUTES is not set, it returns an error.
// Returns an error if PR_APPROVAL_WAIT_MINUTES is invalid.
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
