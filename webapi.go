package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ShufflerMessage represents an encrypted message for the shuffler.
type ShufflerMessage struct {
	Encrypted []byte `json:"encrypted"`
}

// P3AMessage represents a P3A message as it's sent by Brave clients.  See
// the browser code for how messages are created:
// https://github.com/brave/brave-core/blob/1adaa0bc057a83f432e9c278c7c373ef60a5b766/components/p3a/p3a_message.cc#L70
// P3AMessage also implements the Report interface.
type P3AMessage struct {
	YearOfSurvey  int    `json:"yos"`
	YearOfInstall int    `json:"yoi"`
	WeekOfSurvey  int    `json:"wos"`
	WeekOfInstall int    `json:"woi"`
	MetricValue   int    `json:"metric_value"`
	MetricHash    string `json:"metric_hash"`
	CountryCode   string `json:"country_code"`
	Platform      string `json:"platform"`
	Version       string `json:"version"`
	Channel       string `json:"channel"`
	RefCode       string `json:"refcode"`
}

// String returns a string representation of the P3A message.
func (m P3AMessage) String() string {
	return fmt.Sprintf("P3A message:\n"+
		"\tYear of survey:  %d\n"+
		"\tYear of install: %d\n"+
		"\tWeek of survey:  %d\n"+
		"\tWeek of install: %d\n"+
		"\tMetric value:    %d\n"+
		"\tMetric hash:     %s\n"+
		"\tCountry code:    %s\n"+
		"\tPlatform:        %s\n"+
		"\tVersion:         %s\n"+
		"\tChannel:         %s\n"+
		"\tRefcode:         %s\n",
		m.YearOfSurvey, m.YearOfInstall,
		m.WeekOfSurvey, m.WeekOfInstall,
		m.MetricValue, m.MetricHash,
		m.CountryCode, m.Platform, m.Version,
		m.Channel, m.RefCode)
}

// CrowdID returns the crowd ID (a SHA-1 over the message) of the P3A message.
func (m P3AMessage) CrowdID() CrowdID {
	hash := fmt.Sprintf("%x", sha1.Sum(m.Payload()))
	return CrowdID(hash)
}

// Payload returns the P3A message's payload.
func (m P3AMessage) Payload() []byte {
	return []byte(m.String())
}

// createP3AHandler creates a handler that receives a set of JSON-encoded P3A
// messages.
func createP3AHandler(inbox chan []Report) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ms []P3AMessage

		err := json.NewDecoder(r.Body).Decode(&ms)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rs := []Report{}
		for _, m := range ms {
			rs = append(rs, m)
		}

		inbox <- rs
		log.Printf("WebAPI: Sent %d P3A message to shuffler.", len(ms))
	}
}

// createShufflerHandler creates a handler that receives an encrypted blob
// that, when encrypted, contains a JSON-encoded structure consisting of a
// crowd ID and an encrypted payload that is opaque to the shuffler.
func createShufflerHandler(inbox chan []Report) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Decrypt report and forward it to the shuffler's inbox.
	}
}
