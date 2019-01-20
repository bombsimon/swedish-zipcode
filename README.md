# Swedish Zipcode

This package provides validation for Swedish zip codes (or postal codes) based
on the provided [CSV](sweden-zipcode.csv). The CSV was generated at 2019-01-19
by traversing the Bring API described below between `10000` and `99999`.

## HTTPS fallback

Since there're constantly changes in zip codes in Sweden a lot of them are
missing in the CSV submodule. This package supports the possibiity to fallback
to an [HTTPS API from Bring](https://developer.bring.com/api/postal-code/).

If you call `Store()` on the `ZipCodes` type you can store all newly found zip
codes in the existing CSV file for faster and offline execution in the future.

## Multiple results

A few zip codes are listed as multiple matches which means they're shared
between multiple locations. Instead of choosing one of them none is chosen so
to see which these are you can use `grep  ',$' sweden-zipcode.csv`.

**Example**
```
set ex (grep  ',$' sweden-zipcode.csv | string sub --length 5 | head -1); \
    curl -sL \
    "https://api.bring.com/shippingguide/api/postalCode.json?clientUrl=ex&country=SE&pnr=$ex" \
    | jq .multipleMatches

[
  "Södertälje",
  "Enhörna"
]
```

## Examples

```go
// One time validation
sz.Valid("12010") // true
sz.Valid(12010)   // true

// Multiple validations - create a cache
httpFallback := true
zc := sz.NewZipCodes(httpFallback)
zc.ClientURL("https://my.url.se")

zc.Valid("12010") // true
zc.Valid(12010)   // true
zc.Valid(99999)   // false

zc.Store()        // Updates the CSV file with new zip codes.
```
