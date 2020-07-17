//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	// "math"
	"net/http"
	"strconv"
	"time"
)

type StatsData struct {
	Performance        map[string]int
	ProblemDescription string
	Durations          map[string][]float64
	NextPid            int
	PrevPid            int
	Date               string
	PC                 string
	// Durations          map[string]float64
}

//-----------------------------------------------------------------------------------
func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	pid, err := strconv.Atoi(r.FormValue("pid"))
	if err != nil {
		fmt.Println("Unknown problem")
		fmt.Fprintf(w, "Unknown problem")
		return
	}
	if pid <= 0 { // select the last problem (max id)
		// for _, p := range ActiveProblems {
		// 	if pid < p.Info.Pid {
		// 		pid = p.Info.Pid
		// 	}
		// }
		row, err := Database.Query("select id from problem order by id desc limit 1")
		if err != nil {
			fmt.Println("Error retrieving latest problem", err)
			return
		}
		for row.Next() {
			row.Scan(&pid)
		}
		row.Close()
	}
	data := &StatsData{
		Performance: make(map[string]int),
		Durations:   make(map[string][]float64),
		PC:          Passcode,
		NextPid:     pid + 1,
		PrevPid:     pid - 1,
	}
	if pid > 0 {
		rows, err := Database.Query("select score.student_id, score.score, score.graded_submission_number, problem.at, problem.problem_description, submission.id, submission.code_submitted_at, submission.completed from score join problem on score.problem_id=problem.id join submission on score.problem_id=submission.problem_id and score.student_id=submission.student_id where problem.id=? order by submission.id desc", pid)
		if err != nil {
			fmt.Println("Error retrieving problem statistics", pid, err)
			return
		}
		var student_id, score, attempts, sub_id int
		var prob_at, sub_at, sub_completed time.Time
		var prob_content string
		var prob_duration float64
		participants := make(map[int]int)
		for rows.Next() {
			rows.Scan(&student_id, &score, &attempts, &prob_at, &prob_content, &sub_id, &sub_at, &sub_completed)
			// Submission id is ordered descendingly.
			// Therefore, only the last submission of student is looked at.
			if _, ok := participants[student_id]; !ok {
				participants[student_id] = sub_id
				prob_duration = sub_at.Sub(prob_at).Minutes()
				key := fmt.Sprintf("%d points", score)
				data.Performance[key]++
				if _, ok := data.Durations[key]; !ok {
					data.Durations[key] = make([]float64, 0)
				}
				data.Durations[key] = append(data.Durations[key], prob_duration)
				// fmt.Println(data.Performance)
				// fmt.Println(data.Durations)
			}
		}
		rows.Close()

		data.ProblemDescription = prob_content

		the_date := prob_at.Format("2006-01-02")
		rows, err = Database.Query("select student_id, attendance_at from attendance where DATE(at) = ?", the_date)
		var at time.Time
		attendants := make(map[int]int)
		for rows.Next() {
			rows.Scan(&student_id, &at)
			attendants[student_id] = 0
		}
		rows.Close()

		data.Performance["Inactive"] = len(attendants) - len(participants)
		data.Date = the_date
		// data.Performance["Inactive"] = len(Students) - count - 1
		// if data.Performance["Inactive"] < 0 {
		// 	// something is wrong.
		// 	data.Performance["Inactive"] = 0
		// }
	}

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
  <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/numeric/1.2.6/numeric.min.js"></script>
    <style>
    #main { width:1000px; margin: 0 auto;}
    #performance {
      width: 500px;
      height: 500px;
      float: left;
    }
    #durations {
      width: 500px;
      height: 500px;
      margin-left: 500px;
    }
    #pre{ width:100%; display:block;}
    .spacer{ width:100%; height:40px; margin: 0 auto;}
    .pager{ font-size:120%; text-align: center; }
    .pager a{padding:25px; text-decoration: none;}
    .pager a:visited{color:blue}
    </style>
  </head>
  <body>
    <div id="main">
    <div id="performance"></div>
	<div id="durations"></div>
	<script>
  	var perf = [];
	{{ range $key, $val := .Performance }}
		perf.push([{{$key}}, {{$val}}]);
	{{ end }}
	perf.sort();
	var values = [];
	var labels = [];
	for (i=0; i<perf.length; i++ ){
		values.push(perf[i][1]);
		labels.push(perf[i][0]);
	}
	var data = [{
	  values: values,
	  labels: labels,
	  type: 'pie'
	}];
	Plotly.newPlot('performance', data, {'title':'Points'});

	var data2 = [];
	{{ range $key, $val := .Durations }}
		data2.push({
			type:'violin',
			name: {{$key}},
			y: {{$val}},
			box: { visible: true },
			line: { color: 'blue' },
			meanline: { visible: true }
		});
	{{ end }}
	Plotly.newPlot('durations', data2, {
		title:'Time (min)',
		yaxis: {zeroline: false},
	});
    </script>

    <div class="spacer"></div>
    <pre style="padding-left:100px;">{{.Date}}
{{.ProblemDescription}}</pre>
    <div class="spacer"></div>
    <div class="pager">
    <a href="statistics?pc={{.PC}}&pid={{.PrevPid}}">Previous</a>
    <a href="statistics?pc={{.PC}}&pid={{.NextPid}}">Next</a>
    </div>
    </div>
  </body>
</html>
`

var STATS_TEMPLATE_OLD = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['corechart']});
      google.charts.setOnLoadCallback(performancePieChart);
      google.charts.setOnLoadCallback(durationHistogram);

      function performancePieChart() {
      	var perf = [];
		{{ range $key, $val := .Performance }}
			perf.push([{{$key}}, {{$val}}]);
		{{ end }}
		perf.sort();
	    perf.unshift(['Category', 'Count']);
	    console.log(perf);
        var data = google.visualization.arrayToDataTable(perf);
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
    <style>
    #main { width:1000px; margin: 0 auto;}
    #performance {
      width: 500px;
      height: 500px;
      float: left;
    }
    #durations {
      width: 500px;
      height: 500px;
      margin-left: 400px;
    }
    #pre{ width:100%; display:block;}
    .spacer{ width:100%; height:40px; margin: 0 auto;}
    .pager{ font-size:120%; text-align: center; }
    .pager a{padding:25px; text-decoration: none;}
    .pager a:visited{color:blue}
    </style>
  </head>
  <body>
    <div id="main">
    <div id="performance"></div>
    <div id="durations"></div>
    <div class="spacer"></div>
    <pre style="padding-left:100px;">{{.Date}}
{{.ProblemDescription}}</pre>
    <div class="spacer"></div>
    <div class="pager">
    <a href="statistics?pc={{.PC}}&pid={{.PrevPid}}">Previous</a>
    <a href="statistics?pc={{.PC}}&pid={{.NextPid}}">Next</a>
    </div>
    </div>
  </body>
</html>
`
