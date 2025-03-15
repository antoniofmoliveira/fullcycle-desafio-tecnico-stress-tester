package entity

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Red struct {
	Target     string
	SentAt     time.Time
	ReceivedAt time.Time
	StatusCode int
	Payload    string
}

// Get sends a GET request to the url in Target and populates the rest of the
// fields in the Red object. It returns the same object.
//
// If an error occurs while sending the request, the error is logged and the
// function will panic.
//
// If an error occurs while reading the response, the error is logged and the
// function will return the object with the ReceivedAt set to the current time and
// the StatusCode set to -1.
func (r *Red) Get(client *http.Client) *Red {
	req, err := http.NewRequest("GET", r.Target, nil)
	if err != nil {
		slog.Error("(*Red).Get", "msg", err.Error())
		panic(err)
	}
	r.SentAt = time.Now()

	res, err := client.Do(req)
	if err != nil {
		r.ReceivedAt = time.Now()
		r.StatusCode = -1
		return r
	}
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		slog.Error("(*Red).Get io.Copy", "msg", err.Error())
	}
	res.Body.Close()

	r.ReceivedAt = time.Now()
	r.StatusCode = res.StatusCode
	return r
}
