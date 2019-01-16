package sz

import "testing"

func Test_Valid(t *testing.T) {
	cases := []struct {
		description string
		input       interface{}
		valid       bool
	}{
		{
			"invalid string",
			"foobar",
			false,
		},
		{
			"invalid int",
			99999,
			false,
		},
		{
			"valid string",
			"12010",
			true,
		},
		{
			"valid int",
			12010,
			true,
		},
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

func TestCache_Valid(t *testing.T) {
	var (
		cache     = NewCache()
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
		{
			"invalid string",
			"foobar",
			false,
		},
		{
			"invalid int",
			99999,
			false,
		},
		{
			"invalid type",
			[]string{"12010"},
			false,
		},
		{
			"valid byte slice",
			byteVal,
			true,
		},
		{
			"valid string",
			stringVal,
			true,
		},
		{
			"valid string ptr",
			&stringVal,
			true,
		},
		{
			"valid int",
			intVal,
			true,
		},
		{
			"valid int ptr",
			&intVal,
			true,
		},
		{
			"valid int64",
			int64Val,
			true,
		},
		{
			"valid int64 ptr",
			&int64Val,
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			result := cache.Valid(tc.input)

			if result != tc.valid {
				t.Fatal("result and expectation not equal")
			}
		})
	}
}
