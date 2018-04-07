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
		if err == nil {
			ok := false
			var password string
			if r.FormValue("role") == "teacher" {
				password, ok = Teacher[uid]
				if ok && password != r.FormValue("password") {
					ok = false
				}
			} else {
				_, ok = Students[uid]
				if !ok {
					ok = load_and_authorize_student(uid, r.FormValue("password"))
				} else if Students[uid].Password != r.FormValue("password") {
					ok = false
				}
			}
			if ok {
				fn(w, r, r.FormValue("name"), uid)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Unauthorized access:", r.FormValue("name"))
		fmt.Fprint(w, "Unauthorized access. Please register again.")
	}
}

//-----------------------------------------------------------------
