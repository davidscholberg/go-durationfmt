// durationfmt provides a function to format durations according to a format
// string.
package durationfmt

import (
	"fmt"
	"time"
)

const Day = 24 * time.Hour
const Week = 7 * Day
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
// %w - # of weeks
// %d - # of days
// %h - # of hours
// %m - # of minutes
// %s - # of seconds
// %% - print a percent sign
// You can place a 0 before the h, m, and s modifiers to zeropad those values to
// two digits. Zeropadding is undefined for the other modifiers.
func Format(dur time.Duration, fmtStr string) (string, error) {
	var durUnits = map[string]*durationUnit{
		"y": &durationUnit{
			DurDivisor: Year,
		},
		"w": &durationUnit{
			DurDivisor: Week,
		},
		"d": &durationUnit{
			DurDivisor: Day,
		},
		"h": &durationUnit{
			DurDivisor: time.Hour,
		},
		"m": &durationUnit{
			DurDivisor: time.Minute,
		},
		"s": &durationUnit{
			DurDivisor: time.Second,
		},
	}

	sprintfFmt, durCount, err := parseFmtStr(fmtStr, durUnits)
	if err != nil {
		return "", err
	}

	durArray := make([]interface{}, durCount)
	calculateDurUnits(dur, durArray, durUnits)

	return fmt.Sprintf(sprintfFmt, durArray...), nil
}

// calculateDurUnits takes a duration and breaks it up into its constituent
// duration unit values.
func calculateDurUnits(dur time.Duration, durArray []interface{}, durUnits map[string]*durationUnit) {
	remainingDur := dur
	durCount := 0
	for _, c := range "ywdhms" {
		durChar := string(c)
		if durUnits[durChar].Present {
			durArray[durCount] = remainingDur / durUnits[durChar].DurDivisor
			remainingDur = remainingDur % durUnits[durChar].DurDivisor
			durCount++
		}
	}
}

// parseFmtStr parses the given duration format string into its constituent
// units.
// parseFmtStr returns a format string that can be passed to fmt.Sprintf and a
// count of how many duration units are in the format string.
func parseFmtStr(fmtStr string, durUnits map[string]*durationUnit) (string, int, error) {
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
		if _, ok := durUnits[fmtChar]; ok {
			durUnits[fmtChar].Present = true
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
