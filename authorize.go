//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
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
func Authorize(fn func(http.ResponseWriter, *http.Request, string, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := strconv.Atoi(r.FormValue("uid"))
		unauthorized := false
		if err == nil {
			var password string
			var ok bool
			if r.FormValue("role") == "teacher" {
				password, ok = Teacher[uid]
			} else {
				password, ok = Student[uid]
			}
			if ok && password == r.FormValue("password") {
				fn(w, r, r.FormValue("name"), uid)
			} else {
				unauthorized = true
			}
		} else {
			unauthorized = true
		}
		if unauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorized access:", r.FormValue("name"))
			fmt.Fprint(w, "Unauthorized access")
		}
	}
}

//-----------------------------------------------------------------
