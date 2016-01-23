package log_parser

import (
	"testing"
)

func TestExtractingIPs(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here
`

	expected := `8.8.8.8
8.8.4.4
208.122.23.23
`

	test(t, expected, logContents, 1)
}

func TestExtractingTimes(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here
`

	expected := `2015-08-23 12:37:03
2015-08-23 12:37:04
2015-08-23 12:37:05
`

	test(t, expected, logContents, 0)
}

func TestExtractingTexts(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here
`

	expected := `As far as we can tell this is a DNS
Yet another DNS, how quaint!
There is definitely some trend here
`

	test(t, expected, logContents, 2)
}

func TestMoreThanOneEmptyLineAtEndOfLog(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here
`

	expected := `8.8.8.8
8.8.4.4
208.122.23.23
`

	test(t, expected, logContents, 1)
}

func TestLogDoesNotEndInNewLine(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here`

	expected := `8.8.8.8
8.8.4.4
208.122.23.23
`

	test(t, expected, logContents, 1)
}

func TestLogStartsWithNewLine(t *testing.T) {
	logContents := `
	2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here`

	expected := `8.8.8.8
8.8.4.4
208.122.23.23
`

	test(t, expected, logContents, 1)
}

func TestOneLineLog(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS`

	expected := `8.8.8.8
`

	test(t, expected, logContents, 1)
}

func TestEmptyLog(t *testing.T) {
	test(t, ``, ``, 1)
}

func TestEmptyLineLog(t *testing.T) {
	logContents := `
	`
	test(t, ``, logContents, 1)
}

func TestMoreLinesThanExample(t *testing.T) {
	logContents := `2015-08-23 12:37:03 8.8.8.8 As far as we can tell this is a DNS
2015-08-23 12:37:04 8.8.4.4 Yet another DNS, how quaint!
2015-08-23 12:37:05 208.122.23.23 There is definitely some trend here
2015-10-22 08:22:05 127.0.0.1 A campus crashes in the rainbow!
2015-10-22 08:22:06 127.0.0.1 Localhost?? Something is wrong here
2015-10-22 08:22:07 127.0.0.1 The amber libel flies a pope.
2015-10-22 08:22:07 42.42.42.42 Time is an illusion. Lunchtime doubly so.
`

	expected := `8.8.8.8
8.8.4.4
208.122.23.23
127.0.0.1
127.0.0.1
127.0.0.1
42.42.42.42
`

	test(t, expected, logContents, 1)
}

func test(t *testing.T, expected, logContents string, column uint8) {
	found := ExtractColumn(logContents, column)

	if found != expected {
		t.Errorf("Expected\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found)
	}
}
