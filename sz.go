package sz

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
)

// Cache is a type that can be used to only parse the CSV once.
type Cache struct {
	zipCodes map[string]string
}

// NewCache returns a cache that implements the valid function so it can be
// re-used without parsing the CSV.
func NewCache() *Cache {
	return &Cache{
		zipCodes: parseCSV(),
	}
}

// Valid will check if a zip code/postal code is a valid Swedish one.
func (c *Cache) Valid(in interface{}) bool {
	var key string

	switch v := in.(type) {
	case []byte:
		key = string(v)
	case string:
		key = v
	case *string:
		key = *v
	case int:
		key = strconv.Itoa(v)
	case *int:
		key = strconv.Itoa(*v)
	case int64:
		key = strconv.Itoa(int(v))
	case *int64:
		key = strconv.Itoa(int(*v))
	default:
		return false
	}

	_, ok := c.zipCodes[key]

	return ok
}

// Valid will do a one time check if a zip code/postal code is a valid Swedish
// one.
func Valid(in interface{}) bool {
	return NewCache().Valid(in)
}

func parseCSV() map[string]string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("could not read file path")
	}

	filepath := path.Join(path.Dir(filename), "sweden-zipcode/sweden-zipcode.csv")

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	zipCodes := make(map[string]string)
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		zipCodes[record[0]] = record[1]
	}

	return zipCodes
}
