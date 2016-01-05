package wooly

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleParse() {
	// Parse a time string using the layouts defined in wooly.Layouts
	t, _ := Parse(nil, "02 Jan 06 15:04 MST")

	fmt.Println(t.Format(time.UnixDate))
	// output:
	// Mon Jan  2 15:04:00 MST 2006
}

func ExampleParse_ownlayouts() {
	// Parse a time string using a own layout
	t, _ := Parse([]string{"2006-01-02 15:04"}, "2015-12-31 23:59")

	fmt.Println(t.Format(time.UnixDate))
	fmt.Println(t.Layouts())
	// output:
	// Thu Dec 31 23:59:00 UTC 2015
	// [2006-01-02 15:04]
}

func ExampleParse_nil() {
	// When no layouts match, Parse returns nil as Time and an error
	t, err := Parse([]string{""}, "02 Jan 06 15:04 MST")
	fmt.Println(t)
	fmt.Println(err)
	// output:
	// <nil>
	// parsing time "02 Jan 06 15:04 MST": extra text: 02 Jan 06 15:04 MST
}

func ExampleTime_MarshalJSON() {
	// Parse a time string using the layouts defined in wooly.Layouts
	t, _ := Parse(nil, "02 Jan 06 15:04 MST")

	// Marshal to JSON, using the first layout in wooly.Layouts
	buf, _ := json.Marshal(t)

	fmt.Println(string(buf))
	// Output:
	// "2006-01-02T15:04:00Z"
}

func ExampleTime_MarshalJSON_ownlayouts() {
	// Parse a time string using a own layout
	t, _ := Parse([]string{"2006-01-02 15:04"}, "2015-12-31 23:59")

	// Marshal to JSON, the first of the objects own layouts ("2006-01-02 15:04" in this case)
	buf, _ := json.Marshal(t)

	fmt.Println(string(buf))
	// Output:
	// "2015-12-31 23:59"
}

func ExampleTime_UnmarshalJSON() {
	// JSON-RubyDate
	buf := []byte(`"Mon Jan 02 15:04:05 -0700 2006"`)

	var x Time

	json.Unmarshal(buf, &x)

	fmt.Println(x.Time.Format(time.RFC850))
	// output:
	// Monday, 02-Jan-06 15:04:05 -0700
}

func ExampleTime_UnmarshalJSON_ownlayouts() {
	// Timestamp in a really nasty layout
	buf := []byte(`"06 January 2018 11 hours 44 minutes 55 seconds"`)

	var x Time
	// Set the layouts of the time struct
	x.SetLayouts([]string{"02 January 2006 15 hours 04 minutes 05 seconds"})

	json.Unmarshal(buf, &x)

	fmt.Println(x.Time.Format(time.UnixDate))
	// output:
	// Sat Jan  6 11:44:55 UTC 2018
}

func ExampleTime_UnmarshalJSON_fail() {
	// Timestamp in a really nasty layout
	buf := []byte(`"06 January 2018 11 hours 44 minutes 55 seconds"`)

	var x Time

	// This will fail as no matching layout is defined.
	err := json.Unmarshal(buf, &x)

	fmt.Println(err)
	fmt.Println(x.Time.Format(time.UnixDate))
	// output:
	// parsing time "06 January 2018 11 hours 44 minutes 55 seconds" as "Jan _2 15:04:05.000000000": cannot parse "06 January 2018 11 hours 44 minutes 55 seconds" as "Jan"
	// Mon Jan  1 00:00:00 UTC 0001
}
