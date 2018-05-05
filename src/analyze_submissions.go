//
// Author: Vinhthuy Phan, 2018
//
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"sort"
	"time"
)

//-----------------------------------------------------------------------------------
type SubmissionData struct {
	Flag      string
	Start     int64
	At        int64
	Completed int64
}

//-----------------------------------------------------------------------------------
func analyze_submissionsHandler(w http.ResponseWriter, r *http.Request) {
	// if r.FormValue("pc") != Passcode {
	// 	fmt.Fprintf(w, "Unauthorized")
	// 	return
	// }
	pid := r.FormValue("pid")
	records := make(map[int][]*SubmissionData)
	var sid, priority int
	var start, at, completed time.Time

	row, _ := Database.Query("select at from problem where id=?", pid)
	for row.Next() {
		row.Scan(&start)
	}
	row.Close()

	rows, _ := Database.Query("select sid, priority, at, completed from submission where pid=?", pid)
	for rows.Next() {
		rows.Scan(&sid, &priority, &at, &completed)
		if _, ok := records[sid]; !ok {
			records[sid] = make([]*SubmissionData, 0)
		}
		flag := "unknown"
		if priority == 1 {
			flag = "Got it!"
		} else if priority == 2 {
			flag = "Help!"
		}
		records[sid] = append(
			records[sid],
			&SubmissionData{
				Flag:      flag,
				Start:     start.UnixNano(),
				At:        at.UnixNano(),
				Completed: completed.UnixNano(),
			})
	}
	for sid, _ := range records {
		sort.Slice(records[sid], func(i, j int) bool {
			return records[sid][i].At < records[sid][j].At
		})
	}
	rows.Close()
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(ANALYZE_SUBMISSIONS_TEMPLATE)
	err := t.Execute(w, records)
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------------------------
var ANALYZE_SUBMISSIONS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['timeline','corechart','bar']});
      google.charts.setOnLoadCallback(drawSubmissions);
      google.charts.setOnLoadCallback(drawWaitingTime);
      google.charts.setOnLoadCallback(drawResponseTime);
      google.charts.setOnLoadCallback(drawAttempts);


      function drawResponseTime() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Duration'],
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), ({{.At}} - {{.Start}})/1e9],
				{{ end }}
			{{ end }}
        ]);

        var options = {
          title: 'Student response time to a problem',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('response'));
        chart.draw(data, options);
      }

      //---------------------------------------------------
      function drawWaitingTime() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Waiting time'],
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), ({{.Completed}} - {{.At}})/1e9],
				{{ end }}
			{{ end }}
        ]);

        var options = {
          title: 'Waiting time for response from teacher',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('waiting'));
        chart.draw(data, options);
      }

      //---------------------------------------------------
      function drawAttempts() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Attempts'],
			{{ range $sid, $rec := . }}
				{{ $length := len $rec }}
				[  String({{$sid}}), {{$length}} ],
			{{ end }}
        ]);

        var options = {
			title: 'Solution attempts',
			legend: { position: 'none' },
	        hAxis: {
	            viewWindowMode:'explicit',
		        viewWindow: { min:1 }
	        },
        };

        var chart = new google.visualization.Histogram(document.getElementById('attempts'));
        chart.draw(data, options);
      }


      //---------------------------------------------------
      function drawSubmissions() {
        var container = document.getElementById('timeline');
        var chart = new google.visualization.Timeline(container);
        var dataTable = new google.visualization.DataTable();

        dataTable.addColumn({ type: 'string', id: 'Student' });
        dataTable.addColumn({ type: 'string', id: 'Flag' });
        dataTable.addColumn({ type: 'date', id: 'Start' });
        dataTable.addColumn({ type: 'date', id: 'End' });
        dataTable.addRows([
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), {{.Flag}}, new Date({{.At}}/1000000), new Date({{.Completed}}/1000000)],
				{{ end }}
			{{ end }}
		]);
	    var options = {
	    	title: 'Submissions',
	    };
        chart.draw(dataTable, options);
      }
    </script>
    <style>
    #waiting,#attempts,#response{ margin: auto; width:75%; height:300px; }
    #timeline{ margin: auto; width:75%; height:500px; }
    .spacer{ width:100%; height:40px; }
    </style>
  </head>
  <body>
    <div id="attempts"></div>
    <div class="spacer"></div>
    <div id="response"></div>
    <div class="spacer"></div>
    <div id="waiting"></div>
    <div class="spacer"></div>
    <div id="timeline"></div>
  </body>
</html>
`
