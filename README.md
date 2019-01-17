# Swedish Zipcode

This package provides validation for Swedish zip codes (or postal codes) based
on the [CSV](https://github.com/zegl/sweden-zipcode) uploaded by @zegl.

## HTTPS fallback

Since there're constantly changes in zip codes in Sweden a lot of them are
missing in the CSV submodule. This package supports the possibiity to fallback
to an [HTTPS API from Bring](https://developer.bring.com/api/postal-code/).

If you call `Store()` on the `ZipCodes` type you can store all newly found zip
codes in the existing CSV file for faster and offline execution in the future.

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
