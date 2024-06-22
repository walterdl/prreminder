package main

import (
	"slices"
	"time"

	"github.com/walterdl/prremind/notifiertypes"
)

func calcWaitingTime(input notifiertypes.NotifierPayload) (time.Duration, error) {
	tz, err := timezone()
	if err != nil {
		return 0, err
	}
	bHours, err := newBusinessHours()
	if err != nil {
		return 0, err
	}

	now := time.Now().In(tz)
	target := now.Add(time.Duration(bHours.prWaitMinutes) * time.Minute)

	if isInBusinessDay(target, bHours) {
		bDayStart := businessDayStart(target, bHours.start)
		if target.Before(bDayStart) {
			target = addToReachBusinessDayStart(target, bHours.start)
			return target.Sub(now), nil
		}
		return target.Sub(now), nil
	}

	return nextBusinessDayStart(target, bHours).Sub(now), nil
}

func isInBusinessDay(t time.Time, bHours BusinessHours) bool {
	return slices.Contains(bHours.weekdays, t.Weekday())
}

func businessDayStart(t time.Time, startTime ClockTime) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), startTime.hour, startTime.min, 0, 0, t.Location())
}

func addToReachBusinessDayStart(t time.Time, startTime ClockTime) time.Time {
	bDayStart := businessDayStart(t, startTime)
	return t.Add(time.Duration(bDayStart.Sub(t).Minutes()) * time.Minute)
}

func nextBusinessDayStart(t time.Time, bHours BusinessHours) time.Time {
	for {
		t = t.Add(24 * time.Hour)
		if isInBusinessDay(t, bHours) {
			return businessDayStart(t, bHours.start)
		}
	}
}
