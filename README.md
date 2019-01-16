# Swedish Zipcode

This package provides validation for Swedish zip codes (or postal codes) based
on the [CSV](https://github.com/zegl/sweden-zipcode) uploaded by @zegl.

## Examples

```go
// One time validation
sz.Valid("12010") // true
sz.Valid(12010)   // true

// Multiple validations - create a cache
c := sz.NewCache()

c.Valid("12010") // true
c.Valid(12010)   // true
c.Valid(99999)   // false
```
