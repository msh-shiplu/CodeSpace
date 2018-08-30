//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"math"
	"net/http"
	"time"
)

type StatsData struct {
	Performance        map[string]int
	ProblemDescription string
	Durations          map[string]float64
}

//-----------------------------------------------------------------------------------
func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	if r.FormValue("problem") != "latest" {
		fmt.Fprintf(w, "Unsupported")
		return
	}
	max_pid := 0
	for _, p := range ActiveProblems {
		if max_pid < p.Info.Pid {
			max_pid = p.Info.Pid
		}
	}
	data := &StatsData{
		Performance: make(map[string]int),
		Durations:   make(map[string]float64),
	}
	rows, err := Database.Query("select score.stid, score.points, score.attempts, problem.at, problem.content, submission.at, submission.completed from score join problem on score.pid=problem.id join submission on score.pid=submission.pid and score.stid=submission.sid where problem.id=?", max_pid)
	var stid, score, attempts int
	var prob_at, sub_at, sub_completed time.Time
	var prob_content string
	var prob_duration float64
	count := 0
	for rows.Next() {
		rows.Scan(&stid, &score, &attempts, &prob_at, &prob_content, &sub_at, &sub_completed)
		prob_duration = math.Ceil(sub_at.Sub(prob_at).Minutes())
		data.Performance[fmt.Sprintf("%d points", score)]++
		data.Durations[fmt.Sprintf("STID %d", stid)] = prob_duration
		count++
		// fmt.Printf("stid %d, score %d, attempts %d, duration %f\n", stid, score, attempts, prob_duration)
	}
	rows.Close()
	data.Performance["Not submitted"] = len(Students) - count - 1
	if data.Performance["Not submitted"] < 0 {
		// something is wrong.
		data.Performance["Not submitted"] = 0
	}
	data.ProblemDescription = prob_content
	// fmt.Println(data.Performance, len(Students))
	w.Header().Set("Content-Type", "text/html")
	t, err := template.New("").Parse(STATS_TEMPLATE)
	if err != nil {
		fmt.Println(err)
	} else {
		err = t.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//-----------------------------------------------------------------------------------
var STATS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['corechart']});
      google.charts.setOnLoadCallback(performancePieChart);
      google.charts.setOnLoadCallback(durationHistogram);

      function performancePieChart() {
        var data = google.visualization.arrayToDataTable([
          	['Category', 'Count'],
			{{ range $key, $val := .Performance }}
				[  {{$key}}, {{$val}} ],
			{{ end }}
        ]);

        var options = {
          title: 'Performance'
        };

        var chart = new google.visualization.PieChart(document.getElementById('performance'));
        chart.draw(data, options);
      }     

      function durationHistogram() {
        var data = google.visualization.arrayToDataTable([
          ['STID', 'Duration'],
          {{ range $key, $val := .Durations }}
          	[ {{$key}}, {{$val}} ],
          {{ end }}
        ]);

        var options = {
          title: 'Durations (minutes)',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('durations'));
        chart.draw(data, options);
       }      
    </script>
  </head>
  <body>
    <div id="performance" style="width: 900px; height: 500px;"></div>
    <div class="spacer" style="width: 100%; height: 40px;"></div>

    <div id="durations" style="width: 900px; height: 500px;"></div>
    <div class="spacer" style="width: 100%; height: 40px;"></div>
    
    <pre style="padding-left:100px;">
    {{.ProblemDescription}}
    </pre>
  </body>
</html>
`
