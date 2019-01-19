package sz

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	bringAPIURL = "https://api.bring.com/shippingguide/api/postalCode.json"
)

// ZipCodes is a type that can be used to only parse the CSV once or use the
// HTTPS API from Bring.
type ZipCodes struct {
	clientURL    string
	filepath     string
	httpClient   *http.Client
	httpFallback bool
	mu           sync.RWMutex
	zipCodes     map[string]string
}

// NewZipCodes returns a cache that implements the valid function so it can be
// re-used without parsing the CSV.
func NewZipCodes(httpFallback bool) *ZipCodes {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("could not read file path")
	}

	filepath := path.Join(path.Dir(filename), "sweden-zipcode.csv")

	z := &ZipCodes{
		clientURL: "https://change.me.se",
		filepath:  filepath,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		httpFallback: httpFallback,
		mu:           sync.RWMutex{},
	}

	z.Read()

	return z
}

// ClientURL sets the querying client URL according to the Bring API
// documentation.
func (z *ZipCodes) ClientURL(url string) {
	z.clientURL = url
}

// Valid will check if a zip code/postal code is a valid Swedish one.
func (z *ZipCodes) Valid(in interface{}) bool {
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

	z.mu.RLock()
	_, ok := z.zipCodes[key]
	z.mu.RUnlock()

	if !ok && z.httpFallback {
		ok = z.queryBring(key)
	}

	return ok
}

// Valid will do a one time check if a zip code/postal code is a valid Swedish
// one.
func Valid(in interface{}) bool {
	return NewZipCodes(false).Valid(in)
}

// queryBring will query the Bring HTTPS API. See
// https://developer.bring.com/api/postal-code/.
func (z *ZipCodes) queryBring(zipCode string) bool {
	apiURL, err := url.Parse(bringAPIURL)
	if err != nil {
		return false
	}

	query := url.Values{}
	query.Add("clientUrl", z.clientURL)
	query.Add("country", "SE")
	query.Add("pnr", zipCode)

	apiURL.RawQuery = query.Encode()

	response, err := z.httpClient.Get(apiURL.String())
	if err != nil {
		return false
	}

	if response.StatusCode != http.StatusOK {
		return false
	}

	apiResult := struct {
		Result string `json:"result"`
		Valid  bool   `json:"valid"`
		Typ    string `json:"postalCodeType"`
	}{}

	if err = json.NewDecoder(response.Body).Decode(&apiResult); err != nil {
		return false
	}

	if apiResult.Valid {
		z.mu.Lock()
		defer z.mu.Unlock()

		z.zipCodes[zipCode] = apiResult.Result
	}

	return apiResult.Valid
}

func (z *ZipCodes) Read() {
	file, err := os.Open(z.filepath)
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

	z.zipCodes = zipCodes
}

// Store will store all newly found zip codes via external APIs and update the
// CSV file so it's always expanding.
func (z *ZipCodes) Store() error {
	file, err := os.Create(z.filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	zipCodes := []string{}
	records := [][]string{{"Zip", "City"}}

	for k := range z.zipCodes {
		zipCodes = append(zipCodes, k)
	}

	sort.Strings(zipCodes)

	for _, zipCode := range zipCodes {
		records = append(records, []string{zipCode, z.zipCodes[zipCode]})
	}

	w := csv.NewWriter(file)
	if err := w.WriteAll(records); err != nil {
		return err
	}

	return nil
}
