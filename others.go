//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	"net/http"
)

//-----------------------------------------------------------------------------------
func teacher_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------
func testHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Show content of boards
	fmt.Println("Boards", len(Boards))
	for uid, board_pages := range Boards {
		fmt.Printf("Uid: %d has %d pages.\n", uid, len(board_pages))
		for i := 0; i < len(board_pages); i++ {
			b := board_pages[i]
			fmt.Printf("Attempts: %d, Ext: %s, Pid: %d, Answer: %s\n",
				b.Attempts, b.Ext, b.Pid, b.Answer)
			fmt.Println(b.Content)
			fmt.Println("********")
		}
		fmt.Println("--------------------------------")
	}
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------
