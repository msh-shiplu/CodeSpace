//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

func teacher_deactivates_problemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	ActiveProblems = make(map[int]struct{})
	fmt.Fprintf(w, "Ok")
}
