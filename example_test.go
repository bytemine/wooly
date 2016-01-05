package wooly

import (
	"encoding/json"
	"fmt"
)

func ExampleTime_MarshalJSON() {
	// Parse a time string using the layouts defined in wooly.Layouts
	t, err := Parse(nil, "02 Jan 06 15:04 MST")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Marshal to JSON, using the first layout in wooly.Layouts
	buf, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(buf))
	// Output:
	// "2006-01-02T15:04:00Z"
}
