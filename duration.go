package rules

import (
	"bytes"
	"strconv"

	"github.com/pkg/errors"
)

// Duration represents an interval duration
type Duration uint64

// ParseDuration parses a duration
func ParseDuration(str string) (dur Duration, err error) {
	var (
		parsingNum = true
		numStart   int
		numEnd     int
		i          int
		r          rune
	)

	finishPart := func() error {
		c, err := strconv.ParseFloat(str[numStart:numEnd], 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse coefficient")
		}
		if c < 0 {
			return errors.New("negative durations are invalid")
		}

		u, err := parseDurationUnit(str[numEnd:i])
		if err != nil {
			return errors.Wrapf(err, "failed to parse unit %q", str[numEnd:i])
		}

		dur += Duration(c * float64(u.Multiplier))

		numStart, numEnd = i, i
		return nil
	}

	for i, r = range str {
		// fmt.Printf("parsing %c, i %v\n", r, i)
		// check if numeric
		if (r >= '0' && r <= '9') || r == '.' {
			if !parsingNum {
				// fmt.Printf("1\n")
				if err := finishPart(); err != nil {
					return 0, err
				}
			}
			parsingNum = true
			continue
		}
		if parsingNum {
			numEnd = i
			parsingNum = false
		}
	}

	if numEnd > numStart {
		i = len(str)
		if err := finishPart(); err != nil {
			return 0, err
		}
	}

	return dur, nil
}

// MustParseDuration parses the dur or panics
func MustParseDuration(dur string) Duration {
	d, err := ParseDuration(dur)
	if err != nil {
		panic(err)
	}
	return d
}

// String returns a friendly representation of duration
func (d Duration) String() string {
	out := &bytes.Buffer{}

	var (
		co Duration
	)

	if d == Duration(0) {
		return "0s"
	}

	for i := len(durationUnits) - 1; i >= 0 && d > 0; i-- {
		// fmt.Printf("d:%v %v:%v \n", uint64(d), durationUnits[i].Unit, uint64(durationUnits[i].Multiplier))
		if d >= durationUnits[i].Multiplier {
			co = d / durationUnits[i].Multiplier
			out.WriteString(strconv.FormatUint(uint64(co), 10))
			out.WriteString(durationUnits[i].Unit)
			d -= co * durationUnits[i].Multiplier
		}
	}

	return out.String()
}

type durationUnit struct {
	Unit       string
	Multiplier Duration
}

// Duration units
var (
	durationUnits = [...]durationUnit{
		durationUnit{Unit: "s", Multiplier: 1},
		durationUnit{Unit: "m", Multiplier: 60},
		durationUnit{Unit: "h", Multiplier: 60 * 60},
		durationUnit{Unit: "d", Multiplier: 60 * 60 * 24},
		durationUnit{Unit: "w", Multiplier: 60 * 60 * 24 * 7},
		durationUnit{Unit: "mo", Multiplier: 60 * 60 * 24 * 7 * 30},
		durationUnit{Unit: "y", Multiplier: 60 * 60 * 24 * 7 * 30 * 12},
	}
)

func parseDurationUnit(unit string) (durationUnit, error) {
	for _, du := range durationUnits {
		if du.Unit == unit {
			return du, nil
		}
	}
	return durationUnit{}, errors.New("no unit found")
}
