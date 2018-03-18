//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------
// Authorize localhost
//-----------------------------------------------------------------
func AuthorizeLocalhost(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "localhost:8080" {
			fn(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorized access: host is not local.")
			fmt.Fprint(w, "Unauthorized access: host is not local.")
		}
	}
}

//-----------------------------------------------------------------
// Authorize teachers
//-----------------------------------------------------------------
func Authorize(fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password, ok := Teacher[r.FormValue("name")]
		if ok && password == r.FormValue("password") {
			fn(w, r, r.FormValue("name"), r.FormValue("uid"))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorized access:", r.FormValue("name"), r.FormValue("passcode"))
			fmt.Fprint(w, "Unauthorized access")
		}
	}
}

//-----------------------------------------------------------------
