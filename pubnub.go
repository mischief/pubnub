// mischief 2013

// PubNub.com REST API version 3.3 bindings for Go
package pubnub

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// A UUID type for use in the PubNub api
type UUID struct {
	bytes []byte
}

// Generate a UUID. Returns nil, err
// if there is an error getting randomness from Go's
// default random number generator
func UUIDGen() (UUID, error) {
	u := UUID{make([]byte, 16)}

	n, err := rand.Read(u.bytes)
	if n != 16 {
		return UUID{}, errors.New("can't read 16 bytes from random reader")
	} else if err != nil {
		return UUID{}, err
	}

	return u, nil
}

// Format UUID as a string
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u.bytes[0:4], u.bytes[4:6], u.bytes[6:8], u.bytes[8:10], u.bytes[10:])
}

// Default origin for PubNub REST requests
const pubnubOrigin = "pubsub.pubnub.com"

// HTTP headers required by the PubNub API
var pubnubClientHeaders = map[string]string{
	"V":          "3.3",
	"User-Agent": "Go-Google",
	"Accept":     "*/*",
}

// Public interface for PubNub
type PubNubInterface interface {
	Time() (string, error)
	Publish(channel string, message interface{}) (string, error)
	Subscribe(channel string, stopChan <-chan struct{}) (<-chan interface{}, error)
}

// Concrete implementation of PubNubInterface
type PubNub struct {
	publish_key, subscribe_key string
	secret_key, cipher_key     string
	ssl                        bool
	session_uuid               UUID
	origin_url                 string

	time_token string
}

// Internal function to hide the details of getting json objects from the PubNub REST API.
//
// It will construct a URL from origin + urlbits.join('/') + urlparams.
//
// Set encode to make members of urlbits be urlencoded before constructing the url.
//
// Returns jsonobject, nil on success, and nil, err on failure.
func (pn *PubNub) request(urlbits []string, origin string, encode bool, urlparams url.Values) ([]interface{}, error) {
	if urlbits == nil {
		return nil, errors.New("empty urlbits")
	}

	if encode {
		for i, bit := range urlbits {
			urlbits[i] = url.QueryEscape(bit)
		}
	}

	url := pn.origin_url + "/" + strings.Join(urlbits, "/")

	if urlparams != nil {
		url = url + "?" + urlparams.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	for header, value := range pubnubClientHeaders {
		req.Header.Set(header, value)
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var out []interface{}

	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// PubNub.Time
func (pn *PubNub) Time() (string, error) {

	resp, err := pn.request([]string{"time", "0"}, pn.origin_url, false, nil)

	if err != nil {
		return "", err
	}

	time, ok := resp[0].(float64)

	if !ok {
		return "", errors.New("PubNub time response is not a float64")
	}

	return fmt.Sprintf("%.0f", time), nil
}

// PubNub.Publish
// Returns timestamp, nil on success or "", error on failure.
func (pn *PubNub) Publish(channel string, message interface{}) (string, error) {

	json, err := json.Marshal(message)

	if err != nil {
		return "", err
	}

	args := []string{"publish", pn.publish_key, pn.subscribe_key, "0", channel, "0", string(json)}

	query := url.Values{}
	query.Add("uuid", pn.session_uuid.String())

	resp, err := pn.request(args, pn.origin_url, false, query)

	if err != nil {
		return "", err
	}

	// check for api error
	if resp[0].(float64) != 1 {
		return "", errors.New(resp[1].(string))
	}

	timestamp := resp[2].(string)

	return timestamp, nil
}

// PubNub.Subscribe
func (pn *PubNub) Subscribe(channel string, stopChan <-chan struct{}) (<-chan interface{}, error) {

	out := make(chan interface{}, 1)

	// begin subscription
	go func() {
		for {
			select {
			case <-stopChan:
				close(out)
				return
			default:
			}

			args := []string{"subscribe", pn.subscribe_key, channel, "0", pn.time_token}

			//  go func() {
			query := url.Values{}
			query.Add("uuid", pn.session_uuid.String())

			resp, err := pn.request(args, pn.origin_url, true, query)

			if err != nil {
				close(out)
			}
			messages := resp[0].([]interface{})

			pn.time_token = resp[1].(string)

			if len(messages) == 0 {
				// timeout
				continue
			}

			for _, msg := range messages {
				select {
				case out <- msg:
				case <-stopChan:
					close(out)
					return
				}
			}

		}
	}()

	return out, nil
}

// Constructor for a new PubNub client.
func NewPubNub(publish_key, subscribe_key, secret_key, cipher_key string, ssl bool) PubNubInterface {
	pn := &PubNub{
		publish_key:   publish_key,
		subscribe_key: subscribe_key,
		secret_key:    secret_key,
		cipher_key:    cipher_key,
		ssl:           ssl,

		time_token: "0",
	}

	uuid, err := UUIDGen()
	if err != nil {
		panic(err)
	}

	pn.session_uuid = uuid

	if ssl {
		pn.origin_url = "https://" + pubnubOrigin
	} else {
		pn.origin_url = "http://" + pubnubOrigin
	}

	return pn
}
