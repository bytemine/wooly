/*Package wooly wraps Gos default time.Time type to be more fluffy.

The Time type has a custom UnmarshalJSON implementation which trys to parse multiple time formats before failing.
Also the Time type is a pointer (in contrast to the default time.Time), so that zero values are not marshalled for omitempty tagged fields
by the json package.

The most common use should be to import wooly and set/add custom time formats to Layouts in a init function. This way, all methods of wooly.Time use the same layout strings. For special cases, a own set of layouts can be set to structs, but that should only be used in special cases.
*/
package wooly

import (
	"errors"
	"time"
)

// Layouts is a collection of time formats used in JSON unmarshaling of type Time.
// It is used for marshaling and unmarshaling if the Time has no own layouts set, with time.RFC3339Nano used as default format.
// The default set consists of all constants defined in time.
var Layouts = []string{time.RFC3339Nano, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano}

// ErrNoLayout is returned if no layouts are defined.
var ErrNoLayout = errors.New("No layouts defined")

// parseTime tries to parse a string as multiple time formats (like those of time.Parse).
// The first match wins. If no format can be parsed successfully, the error for the last
// tested formatstring is returned.
func parseTime(layouts []string, value string) (time.Time, error) {
	var err error
	for _, v := range layouts {
		var t time.Time
		t, err = time.Parse(v, value)
		
		// continue to the next layout
		if err != nil {
			continue
		}

		// parsed successfully
		return t, nil
	}

	if err != nil {
		return time.Time{}, err
	}

	return time.Time{}, ErrNoLayout
}

// Time is time.Time with custom json marshaling methods.
// Each Time object can have its own list of layouts. If it has own layouts these override the
// package-global layouts.
type Time struct {
	time.Time
	layouts []string
}

// New returns a new Time value from a time.Time value.
func New(t time.Time) *Time {
	return &Time{Time: t}
}

// Parse parses a formatted string and returns the time value it represents. It tries to parse with all layouts, returning the error
// of the last tried layout if none succeeds. If layouts is nil, wooly.Layouts is used. If it is not nil, the returned time has the supplied
// layouts set.
func Parse(layouts []string, value string) (*Time, error) {
	t := new(Time)
	if layouts == nil {
		layouts = Layouts
	} else {
		t.layouts = layouts
	}

	x, err := parseTime(layouts, value)
	if err != nil {
		return nil, err
	}

	t.Time = x

	return t, nil
}

// Layouts returns the layouts a Time object uses.
func (t *Time) Layouts() []string {
	return t.layouts
}

// SetLayouts sets the layouts a Time object uses.
func (t *Time) SetLayouts(layouts []string) {
	t.layouts = layouts
}

// selectLayouts chooses from the structs and the package layouts.
func (t *Time) selectLayouts() ([]string, error) {
	var l []string
	if t.layouts != nil {
		l = t.layouts
	} else if Layouts != nil {
		l = Layouts
	} else {
		return nil, ErrNoLayout
	}

	if len(l) < 1 {
		return nil, ErrNoLayout
	}

	return l, nil
}

// MarshalJSON implements the json.Marshaler interface.
// To format the date, the first element of layouts is used (with the objects own layouts overriding the packages).
func (t *Time) MarshalJSON() ([]byte, error) {
	layouts, err := t.selectLayouts()
	if err != nil {
		return nil, err
	}

	x := []byte(t.Format(`"` + layouts[0] + `"`))

	return x, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The unmarshaling only fails when parsing has failed for every layout (with the objects own layouts overriding the packages).
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	layouts, err := t.selectLayouts()
	if err != nil {
		return err
	}

	t.Time, err = parseTime(layouts, string(data[1:len(data)-1]))

	return err
}
