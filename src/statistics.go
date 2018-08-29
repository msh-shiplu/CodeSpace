//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"time"
)

//-----------------------------------------------------------------------------------
// type DailyActivityData struct {
// 	Pids     map[int]bool
// 	Sids     map[int]bool
// 	Count    int
// 	PidCount int
// 	SidCount int
// }

//-----------------------------------------------------------------------------------
func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	// if r.FormValue("pc") != Passcode {
	// 	fmt.Fprintf(w, "Unauthorized")
	// 	return
	// }
	if r.FormValue("problem") != "latest" {
		fmt.Fprintf(w, "Unsupported")
		return
	}
	fmt.Println("Students", Students)
	max_pid := 0
	for _, p := range ActiveProblems {
		if max_pid < p.Info.Pid {
			max_pid = p.Info.Pid
		}
	}
	fmt.Println(">", max_pid)
	rows, err := Database.Query("select score.stid, score.points, score.attempts, problem.at, submission.at, submission.completed from score join problem on score.pid = problem.id join submission on score.pid==submission.pid where problem.id=?", max_pid)
	var stid, score, attempts int
	var prob_at, sub_at, sub_completed time.Time
	var prob_duration float64
	for rows.Next() {
		rows.Scan(&stid, &score, &attempts, &prob_at, &sub_at, &sub_completed)
		prob_duration = sub_at.Sub(prob_at).Seconds()
		fmt.Println(stid, score, attempts, prob_duration)
		fmt.Println(prob_at)
		fmt.Println(sub_at)
	}
	rows.Close()
	// rows, _ := Database.Query("select pid, sid, at from submission")
	// var at time.Time
	// var pid, sid int
	// data := make(map[int64]*DailyActivityData)
	// for rows.Next() {
	// 	rows.Scan(&pid, &sid, &at)
	// 	date := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, at.Location()).UnixNano()
	// 	if _, ok := data[date]; !ok {
	// 		data[date] = &DailyActivityData{
	// 			Pids:  make(map[int]bool),
	// 			Sids:  make(map[int]bool),
	// 			Count: 0,
	// 		}
	// 	}
	// 	data[date].Count++
	// 	data[date].Pids[pid] = true
	// 	data[date].Sids[sid] = true
	// }
	// rows.Close()
	// for d, _ := range data {
	// 	data[d].PidCount = len(data[d].Pids)
	// 	data[d].SidCount = len(data[d].Sids)
	// }
	w.Header().Set("Content-Type", "text/html")
	t, err := template.New("").Parse(STATS_TEMPLATE)
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------------------------
var STATS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','line']});
    </script>
    <style>
    #chart1_div,#chart2_div{ margin: auto; width:75%; }
    .spacer{ width:100%; height:40px; }
    </style>
  </head>

  <body>
    <div id="chart1_div"></div>
    <div class="spacer"></div>
    <div id="chart2_div"></div>
  </body>
</html>
`
