[![](https://godoc.org/github.com/go-apilayer/ipstack?status.svg)](http://godoc.org/github.com/go-apilayer/ipstack)

## ipstack

`ipstack` is a Go client library for the [apilayer ipstack](https://ipstack.com/) service, which provides an API to identify website visitors by IP address.

To use this client you will need an **API access key**. The free tier supports 10,000 monthly lookups, but over `http-only`. To get HTTPS Encryption you need a Basic plan or higher 🤷‍♂️.

The official documentation can be found [here](https://ipstack.com/documentation).

If you are on a free account, initialize the client with secure `false`. Otherwise you'll get `105` error:

> Access Restricted - Your current Subscription Plan does not support HTTPS Encryption.

---

### Technical bits: 

Users can supply their own HTTP Client implementation through options, otherwise a default client is used.

This library will return a custom error, `*ApiErr`, which callers can assert to get at the raw code, type and info. If using go1.13 use `errors.As` for assertions, otherwise a regular type switch will do. This is especially useful for `104` or `ErrUsageLimitReached` errors: monthly usage limit has been exceeded.

TimeZone, Currency, Connection, and Security are only available for Basic Plan or higher. If you don't see this in the response, check plan. Make sure to check the field is not `nil`.

### Example usage - single lookup (HTTP GET)

```go
client, err := ipstack.NewClient(key, false)
if err != nil {
    // handler err
}

// request param is optional, this is just an example. Some items will be in German.
ipStack, err := client.Lookup("134.201.250.155", ipstack.RequestParam{Language: ipstack.LangGerman})
if err != nil {
    var e *ipstack.ApiErr
    if errors.As(err, &e) { 
        // handler error of type ApiErr
    }
    // handler err
}

fmt.Printf("Country:%s\nCity:%s\n", ipStack.CountryName, ipStack.City)
// Country:Vereinigte Staaten
// City:Los Angeles
```
