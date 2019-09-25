package ipstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	apiHost = "api.ipstack.com"
)

// defaultClient is an http client with sane defaults.
var defaultClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,

		ExpectContinueTimeout: 10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	},
	Timeout: 60 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// NewClient returns a ipstack client. Users can modify clients with functional options.
//
// Note: if you have a non-paying account, you must specify secure: false. Only paid accounts
// get access to `https`.
func NewClient(apiKey string, secure bool, options ...Option) (*Client, error) {
	if apiKey == "" {
		b := false
		err := ApiErr{Success: &b}
		err.Err.Type = ErrMissingAccessKey
		err.Err.Code = codeFromErrorType(ErrMissingAccessKey)
		err.Err.Info = "No API Key was specified."
		return nil, &err
	}

	c := &Client{client: defaultClient}

	u := &url.URL{
		Scheme: "http", // Unpaid accounts do not have access to https, sadly.
		Host:   apiHost,
	}
	if secure {
		u.Scheme = "https"
	}
	q := u.Query()
	q.Add("access_key", apiKey)
	u.RawQuery = q.Encode()

	c.url = u

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

type Client struct {
	client HTTPClient
	url    *url.URL

	debug bool
}

func (c *Client) debugf(format string, v ...interface{}) {
	if c.debug {
		for i := range v {
			if u, ok := v[i].(*url.URL); ok && u != nil {
				copy := new(url.URL)
				*copy = *u
				q := copy.Query()
				q.Set("access_key", "hidden")
				copy.RawQuery = q.Encode()
				v[i] = copy
			}
		}
		msg := fmt.Sprintf(format, v...)
		log.Println("go-apilayer/ipstack:", msg)
	}
}

type RequestParam struct {
	// Set to your preferred output field(s). Follow docs to compose this.
	Fields string

	// Enable Hostname Lookup.
	HostName bool

	// Enable the Security module.
	Security bool

	// Set to a 2-letter language code to change output language.
	Language LangType
}

// Lookup (standard) looks up the data behind an IP address. IP can be IPv4, IPv6 or
// domain URL (ipstack will resolve domain to the underlying IP address).
//
// Only first param struct is used, do not pass more than one.
func (c *Client) Lookup(IPAddr string, params ...RequestParam) (*Stack, error) {
	// deep copy URL
	u := new(url.URL)
	*u = *c.url

	u.Path = url.PathEscape(IPAddr)
	q := u.Query()

	if len(params) > 0 {
		param := params[0]

		if param.HostName {
			q.Add("hostname", "1")
		}
		if param.Security {
			q.Add("security", "1")
		}
		if param.Language != "" {
			q.Add("language", param.Language.String())
		}
		if param.Fields != "" {
			q.Add("fields", param.Fields)
		}
	}
	u.RawQuery = q.Encode()

	c.debugf("HTTP request: %v", u)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.debugf("HTTP GET:%d header:%+v", resp.StatusCode, resp.Header)

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiErr *ApiErr
	if err := json.Unmarshal(by, &apiErr); err != nil {
		return nil, err
	}
	if apiErr.Success != nil && !*apiErr.Success {
		return nil, apiErr
	}

	var std *Stack
	if err := json.Unmarshal(by, &std); err != nil {
		return nil, err
	}

	return std, nil
}

// .. Look up the data behind multiple IP addresses at once.
// Maxmium allowed values: 50
func (c *Client) BulkLookup() ([]*Stack, error) {
	ss := []*Stack{}
	return ss, nil
}

// .. Look up the data behind the IP address your API request is coming from.
func (c *Client) RequesterLookup() (*Stack, error) {

	return &Stack{}, nil
}

// Stack is the ipstack response object.
//
// Note: some fields (referred to as modules) are only available to certain
// subscription plans. Make sure to check for existence prior to accessing their fields.
type Stack struct {
	IP            string   `json:"ip"`
	Hostname      string   `json:"hostname"`
	Type          string   `json:"type"`
	ContinentCode string   `json:"continent_code"`
	ContinentName string   `json:"continent_name"`
	CountryCode   string   `json:"country_code"`
	CountryName   string   `json:"country_name"`
	RegionCode    string   `json:"region_code"`
	RegionName    string   `json:"region_name"`
	City          string   `json:"city"`
	Zip           string   `json:"zip"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	Location      Location `json:"location"`

	// These fields varying depending on your subscription plan
	// and may be nil.
	TimeZone   *TimeZone   `json:"time_zone,omitempty"`
	Currency   *Currency   `json:"currency,omitempty"`
	Connection *Connection `json:"connection,omitempty"`
	Security   *Security   `json:"security,omitempty"`
}

type Location struct {
	GeonameID int    `json:"geoname_id"`
	Capital   string `json:"capital"`
	Languages []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Native string `json:"native"`
	} `json:"languages"`
	CountryFlag             string `json:"country_flag"`
	CountryFlagEmoji        string `json:"country_flag_emoji"`
	CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
	CallingCode             string `json:"calling_code"`
	IsEu                    bool   `json:"is_eu"`
}
type TimeZone struct {
	ID               string `json:"id"`
	CurrentTime      string `json:"current_time"`
	GmtOffset        int    `json:"gmt_offset"`
	Code             string `json:"code"`
	IsDaylightSaving bool   `json:"is_daylight_saving"`
}
type Currency struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Plural       string `json:"plural"`
	Symbol       string `json:"symbol"`
	SymbolNative string `json:"symbol_native"`
}
type Connection struct {
	Asn int    `json:"asn"`
	Isp string `json:"isp"`
}
type Security struct {
	IsProxy     bool        `json:"is_proxy"`
	ProxyType   interface{} `json:"proxy_type"`
	IsCrawler   bool        `json:"is_crawler"`
	CrawlerName string      `json:"crawler_name"`
	CrawlerType interface{} `json:"crawler_type"`
	IsTor       bool        `json:"is_tor"`
	ThreatLevel string      `json:"threat_level"`
	ThreatTypes []string    `json:"threat_types"`
}
