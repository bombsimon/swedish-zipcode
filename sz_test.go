package sz

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Valid(t *testing.T) {
	cases := []struct {
		description string
		input       interface{}
		valid       bool
	}{
		{"invalid string", "foobar", false},
		{"invalid int", 99999, false},
		{"valid string", "12010", true},
		{"valid int", 12010, true},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			result := Valid(tc.input)

			if result != tc.valid {
				t.Fatal("result and expectation not equal")
			}
		})
	}
}

func TestZipCodes_Valid(t *testing.T) {
	var (
		zc        = NewZipCodes(false)
		byteVal   = []byte("12010")
		stringVal = "12010"
		intVal    = 12010
		int64Val  = int64(12010)
	)

	cases := []struct {
		description string
		input       interface{}
		valid       bool
	}{
		{"invalid string", "foobar", false},
		{"invalid int", 99999, false},
		{"invalid type", []string{"12010"}, false},
		{"valid byte slice", byteVal, true},
		{"valid string", stringVal, true},
		{"valid string ptr", &stringVal, true},
		{"valid int", intVal, true},
		{"valid int ptr", &intVal, true},
		{"valid int64", int64Val, true},
		{"valid int64 ptr", &int64Val, true},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			result := zc.Valid(tc.input)

			if result != tc.valid {
				t.Fatal("result and expectation not equal")
			}
		})
	}
}

func TestZipCodes_ValidWithHTTP(t *testing.T) {
	var (
		zc = NewZipCodes(true)

		existingZipCode = "47150"
		csvDir, _       = filepath.Split(zc.filepath)
		testCsv         = filepath.Join(csvDir, "test-csv.csv")
	)

	defer func() {
		if err := os.Remove(testCsv); err != nil {
			t.Fatal(err.Error())
		}
	}()

	// Set the client URL according to API specs.
	zc.ClientURL("https://github.com/bombsimon/swedish-zipcode")

	// Set the filepath to the temporary test file to store and re-parse.
	zc.filepath = testCsv

	// Unset all zip codes from the full CSV file to ensure we've always empty.
	zc.zipCodes = map[string]string{}

	cases := []struct {
		description string
		input       interface{}
		valid       bool
	}{
		{"not found in CSV, not found in API", "99999", false},
		{"not found in CSV, found in API", existingZipCode, true},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			result := zc.Valid(tc.input)

			if result != tc.valid {
				t.Fatal("result and expectation not equal")
			}
		})
	}

	if err := zc.Store(); err != nil {
		t.Fatal("could not store newly found zip codes")
	}

	// Ensure we've cleared all cached zip codes.
	zc.zipCodes = map[string]string{}
	zc.Read()

	if _, ok := zc.zipCodes[existingZipCode]; !ok {
		t.Fatal("zip code was not stored in CSV")
	}
}
