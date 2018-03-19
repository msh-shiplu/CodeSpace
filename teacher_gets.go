//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
)

//-----------------------------------------------------------------------------------
func teacher_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	BoardsSem.Lock()
	defer BoardsSem.Unlock()
	fmt.Println(priority)
	fmt.Fprintf(w, "Ok")
}

//-----------------------------------------------------------------------------------
