package wooly

import (
	"testing"
	"time"
)

var testLayouts = []string{time.ANSIC, time.RubyDate}

func TestParseTime(t *testing.T) {
	x, err := parseTime(testLayouts, "Mon Jan 02 15:04:05 -0700 2006")
	if err != nil {
		t.Error(err)
	}

	if h, m, s := x.Clock(); h != 15 || m != 04 || s != 05 {
		t.Fail()
	}

	if y, m, d := x.Date(); y != 2006 || m != time.January || d != 2 {
		t.Fail()
	}

	_, err = parseTime(testLayouts, "FailDate")
	if err == nil {
		t.Fail()
	}
}

func TestMarshalJSON(t *testing.T) {
	y, err := time.Parse(time.RubyDate, "Mon Jan 02 15:04:05 -0700 2006")

	x := New(y)
	jx, err := x.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	if string(jx) != `"Mon Jan  2 15:04:05 2006"` {
		t.Log(string(jx))
		t.Fail()
	}
}

func TestUnmarshalJSON(t *testing.T) {
	x := new(Time)
	x.SetLayouts(testLayouts)

	// ANSIC
	err := x.UnmarshalJSON([]byte(`"Mon Jan  2 15:04:05 2006"`))
	if err != nil {
		t.Error(err)
	}

	if h, m, s := x.Clock(); h != 15 || m != 04 || s != 05 {
		t.Fail()
	}

	if y, m, d := x.Date(); y != 2006 || m != time.January || d != 2 {
		t.Fail()
	}

	// RubyDate
	err = x.UnmarshalJSON([]byte(`"Mon Jan 02 15:04:05 -0700 2006"`))
	if err != nil {
		t.Error(err)
	}

	if h, m, s := x.Clock(); h != 15 || m != 04 || s != 05 {
		t.Fail()
	}

	if y, m, d := x.Date(); y != 2006 || m != time.January || d != 2 {
		t.Fail()
	}

	// Kitchen fails as it is not in x.layouts
	err = x.UnmarshalJSON([]byte(`"3:04PM"`))
	if err == nil {
		t.Fail()
	}
}
