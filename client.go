/*Package icinga2 provides Icinga2 API client functionality.

See the tests for basic usage examples. */
package icinga2

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/bytemine/go-icinga2/event"
)

const icingaAPI = "v1"

// Client is a Icinga2 client.
type Client struct {
	url                *url.URL
	user               string
	password           string
	insecureSkipVerify bool
}

// NewClient prepares a Client for usage.
//
// The icingaURL should contain the path up to the API-version part, eg. http://example.org:5665/ when the API lives at
// http://example.org:5665/v1 . User and password are the credentials of the API user. Set insecureSkipVerify to true
// if you have an invalid certificate chain.
func NewClient(icingaURL string, user, password string, insecureSkipVerify bool) (*Client, error) {
	x, err := url.Parse(icingaURL)
	if err != nil {
		return nil, err
	}
	return &Client{url: x, user: user, password: password, insecureSkipVerify: insecureSkipVerify}, nil
}

// EventStream opens an Icinga2 event stream.
//
// Queue is an "unique" queue name, though multiple clients can use the same name if they use the same parameters for filter and streamtype.
//
// Filter is an Icinga API filter, described at: http://docs.icinga.org/icinga2/snapshot/doc/module/icinga2/chapter/icinga2-api#icinga2-api-filters
//
// Streamtype selects the type of events the EventStreamer should listen for, see the package constants.
//
// The returned io.Reader can directly be used with a json.Decoder if only one StreamType is requested. Otherwise event.Mux can be used to
// split the stream into seperate steams for each type.
func (c *Client) EventStream(queue string, filter string, streamtype ...event.StreamType) (io.Reader, error) {
	x := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.insecureSkipVerify}}, Timeout: 0}

	var q url.Values = make(map[string][]string)
	types := []string{}
	for _, v := range streamtype {
		types = append(types, string(v))
	}
	q["types"] = types

	if filter != "" {
		q.Set("filter", filter)
	}

	if queue == "" {
		return nil, errors.New("queue name can't have zero value")
	}

	q.Set("queue", queue)

	u := url.URL{Scheme: "https", Host: c.url.Host, Path: filepath.Join(c.url.Path, icingaAPI, "events"), RawQuery: q.Encode()}

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.user, c.password)

	req.Header.Add("Accept", "application/json")

	res, err := x.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
