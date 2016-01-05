# wooly
## about
Package wooly wraps Gos default time.Time type to be more fluffy.

The Time type has a custom UnmarshalJSON implementation which trys to parse multiple time formats before failing. Also the Time type is a pointer (in contrast to the default time.Time), so that zero values are not marshalled for omitempty tagged fields by the json package.

The most common use should be to import wooly and set/add custom time formats to Layouts in a init function. This way, all methods of wooly.Time use the same layout strings. For special cases, a own set of layouts can be set to structs, but that should only be used in special cases.

## documentation
[Documentation on godoc.org](https://godoc.org/github.com/bytemine/wooly)
