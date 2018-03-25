//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

func student_checks_inHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// attendance is taken automatically by authorization
	fmt.Fprint(w, "Ok")
}
