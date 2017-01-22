package durationfmt

import (
	"fmt"
	"time"
)

const Day = 24 * time.Hour
const Week = 7 * Day
const Year = 365 * Day

type duration struct {
	Present    bool
	DurDivisor time.Duration
}

//type durationUnits struct {
//	Years   duration
//	Weeks   duration
//	Days    duration
//	Hours   duration
//	Minutes duration
//	Seconds duration
//}

// Format formats the given duration according to the given format string.
// %y - # of years
// %w - # of weeks
// %d - # of days
// %h - # of hours
// %m - # of minutes
// %s - # of seconds
// %% - print a percent sign
func Format(dur time.Duration, fmtStr string) (string, error) {
	var durationUnits = map[string]*duration{
		"y": &duration{
			DurDivisor: Year,
		},
		"w": &duration{
			DurDivisor: Week,
		},
		"d": &duration{
			DurDivisor: Day,
		},
		"h": &duration{
			DurDivisor: time.Hour,
		},
		"m": &duration{
			DurDivisor: time.Minute,
		},
		"s": &duration{
			DurDivisor: time.Second,
		},
	}

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
		if _, ok := durationUnits[fmtChar]; ok {
			durationUnits[fmtChar].Present = true
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
				return "", fmt.Errorf("incorrect duration modifier")
			}
		}
		modifier = false
	}

	remainingDur := dur
	durationArray := make([]interface{}, durCount)
	durCount = 0
	for _, c := range "ywdhms" {
		durChar := string(c)
		if durationUnits[durChar].Present {
			durationArray[durCount] = remainingDur / durationUnits[durChar].DurDivisor
			remainingDur = remainingDur % durationUnits[durChar].DurDivisor
			durCount++
		}
	}

	return fmt.Sprintf(sprintfFmt, durationArray...), nil
}
