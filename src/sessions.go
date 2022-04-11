package main

import "time"

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
