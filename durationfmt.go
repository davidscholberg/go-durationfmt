// durationfmt provides a function to format durations according to a format
// string.
package durationfmt

import (
	"fmt"
	"time"
)

const Day = 24 * time.Hour
const Week = 7 * Day
const Month = 30 * Day
const Year = 365 * Day

// durationUnit represets a possible duration unit. A durationUnit object
// contains the divisor that the duration unit uses as well as if that duration
// unit is present in the duration format.
type durationUnit struct {
	Present    bool
	DurDivisor time.Duration
}

// Format formats the given duration according to the given format string.
// %y - # of years
// %o - # of months
// %w - # of weeks
// %d - # of days
// %h - # of hours
// %m - # of minutes
// %s - # of seconds
// %i - # of milliseconds
// %c - # of microseconds
// %n - # of nanoseconds
// %% - print a percent sign
// You can place a 0 before the h, m, and s modifiers to zeropad those values to
// two digits. Zeropadding is undefined for the other modifiers.
func Format(dur time.Duration, fmtStr string) (string, error) {
	var durUnitSlice = []durationUnit{
		{
			DurDivisor: Year,
		},
		{
			DurDivisor: Month,
		},
		{
			DurDivisor: Week,
		},
		{
			DurDivisor: Day,
		},
		{
			DurDivisor: time.Hour,
		},
		{
			DurDivisor: time.Minute,
		},
		{
			DurDivisor: time.Second,
		},
		{
			DurDivisor: time.Millisecond,
		},
		{
			DurDivisor: time.Microsecond,
		},
		{
			DurDivisor: time.Nanosecond,
		},
	}
	var durUnitMap = map[string]*durationUnit{
		"y": &durUnitSlice[0],
		"o": &durUnitSlice[1],
		"w": &durUnitSlice[2],
		"d": &durUnitSlice[3],
		"h": &durUnitSlice[4],
		"m": &durUnitSlice[5],
		"s": &durUnitSlice[6],
		"i": &durUnitSlice[7],
		"c": &durUnitSlice[8],
		"n": &durUnitSlice[9],
	}

	sprintfFmt, durCount, err := parseFmtStr(fmtStr, durUnitMap)
	if err != nil {
		return "", err
	}

	durArray := make([]interface{}, durCount)
	calculateDurUnits(dur, durArray, durUnitSlice)

	return fmt.Sprintf(sprintfFmt, durArray...), nil
}

// calculateDurUnits takes a duration and breaks it up into its constituent
// duration unit values.
func calculateDurUnits(dur time.Duration, durArray []interface{}, durUnitSlice []durationUnit) {
	remainingDur := dur
	durCount := 0
	for _, d := range durUnitSlice {
		if d.Present {
			durArray[durCount] = remainingDur / d.DurDivisor
			remainingDur = remainingDur % d.DurDivisor
			durCount++
		}
	}
}

// parseFmtStr parses the given duration format string into its constituent
// units.
// parseFmtStr returns a format string that can be passed to fmt.Sprintf and a
// count of how many duration units are in the format string.
func parseFmtStr(fmtStr string, durUnitMap map[string]*durationUnit) (string, int, error) {
	modifier, zeropad := false, false
	sprintfFmt := ""
	durCount := 0
	for _, c := range fmtStr {
		fmtChar := string(c)
		if modifier == false {
			if fmtChar == "%" {
				modifier = true
			} else {
				sprintfFmt += fmtChar
			}
			continue
		}
		if _, ok := durUnitMap[fmtChar]; ok {
			durUnitMap[fmtChar].Present = true
			durCount++
			if zeropad {
				sprintfFmt += "%02d"
				zeropad = false
			} else {
				sprintfFmt += "%d"
			}
		} else {
			switch fmtChar {
			case "0":
				zeropad = true
				continue
			case "%":
				sprintfFmt += "%%"
			default:
				return "", durCount, fmt.Errorf("incorrect duration modifier")
			}
		}
		modifier = false
	}
	return sprintfFmt, durCount, nil
}
