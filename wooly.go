/* Package wooly wraps Gos default time.Time type to be more fluffy.

It sports a custom UnmarshalJSON implementation which trys to parse multiple time formats before failing.
Also the Time type is a pointer (in contrast to the default time.Time), so that zero values are not marshalled for omitempty tagged fields
by the json package.
*/
package wooly

import (
	"time"
)

// Layouts is a collection of time formats used in JSON unmarshaling of type Time.
var Layouts = []string{time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano}

// parseTime tries to parse a string as multiple time formats (like those of time.Parse).
// The first match wins. If no format can be parsed successfully, the error for the last
// tested formatstring is returned.
func parseTime(layouts []string, value string) (time.Time, error) {
	var err error
	for _, v := range layouts {
		var t time.Time
		t, err = time.Parse(v, value)
		if err != nil {
			continue
		}
		return t, nil
	}
	return time.Time{}, err
}

// Time is time.Time with custom json marshaling methods.
// Each Time object can have its own list of layouts. If it has own layouts these override the
// package-global layouts.
type Time struct {
	time.Time
	layouts []string
}

// NewTime returns a new Time value from a time.Time value.
func New(t time.Time) *Time {
	return &Time{Time: t}
}

// Parse parses a formatted string and returns the time value it represents. It tries to parse with all layouts, returning the error
// of the last tried layout if none succeeds. If layouts is nil, wooly.Layouts is used.
func Parse(layouts []string, value string) (*Time, error) {
	if layouts == nil {
		layouts = Layouts
	}
	x, err := parseTime(layouts, value)
	return New(x), err
}

func (t *Time) Layouts() []string {
	return t.layouts
}

func (t *Time) SetLayouts(layouts []string) {
	t.layouts = layouts
}

// MarshalJSON implements the json.Marshaler interface.
// To format the date, the first element of layouts is used (with the objects own layouts overriding the packages).
func (t *Time) MarshalJSON() ([]byte, error) {
	var layout string
	if t.layouts != nil {
		layout = t.layouts[0]
	} else {
		layout = Layouts[0]
	}

	return []byte(t.Format(`"` + layout + `"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The unmarshaling only fails when parsing has failed for every layout.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var layouts []string
	if t.layouts != nil {
		layouts = append(t.layouts, Layouts...)
	} else {
		layouts = Layouts
	}
	t.Time, err = parseTime(layouts, string(data[1:len(data)-1]))
	return
}