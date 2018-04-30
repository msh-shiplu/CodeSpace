//
// Author: Vinhthuy Phan, 2018
//
package main

var STUDENT_MESSAGING_TEMPLATE = `
<html>
	<head>
  		<title>Student messaging</title>
		<meta http-equiv="refresh" content="10" />
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			font-size: 150%;
			color: red;
		}
	</style>
	<body>
	<div class="bottom">{{.Message}}</div>
	</body>
</html>
`

var TEACHER_MESSAGING_TEMPLATE = `
<html>
	<head>
  		<title>Teacher messaging</title>
		<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js?autoload=true&skin=sons-of-obsidian"></script>
  		<script src="http://code.jquery.com/jquery-3.1.1.min.js"></script>
	    <script type="text/javascript">
			var updateInterval = 5000;		// 5 sec update interval
			var maxUpdateTime =  1800000;   // no longer update after 30 min.
			var totalUpdateTime = 0;
			function getData() {
				var url = "http://{{.Address}}/bulletin_board_data";
				$.getJSON(url, function( data ) {
					console.log(data);
					$("#p1").html(data["P1"]);
					$("#p2").html(data["P2"]);
					$("#ap").html(data["ActiveProblems"]);
					$("#bu").html(data["BulletinItems"]);
					$("#at").html(data["Attendance"]);
				});
			}
			$(document).ready(function(){
				getData();
				handle = setInterval(getData, updateInterval);
			});
	    </script>
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			text-align: center;
			width: 100%;
		}
		.label{ display: inline; }
		#p1, #p2, #ap, #bu, #at {
			padding: 0.75em;
			display: inline;
		}
		#p1 { color: green; }
		#p2 { color: red; }
		pre {
			font-family: monospace;
			font-size:120%;
			margin-top:50px;
			padding-left:2em;
			overflow-x:scroll;
			overflow-y:scroll;
			tab-size: 4;
			-moz-tab-size: 4;
		}
		.center {
		    text-align: center;
		}
		.pagination {
		    display: inline-block;
		    padding-bottom: 20px;
		}
		.pagination a {
		    color: black;
		    float: left;
		    padding: 8px 16px;
		    text-decoration: none;
		    transition: background-color .3s;
		    border: 1px solid #ddd;
		    margin: 0 4px;
		    border-radius: 5px;
		}
		.pagination a.active {
		    background-color: #4CAF50;
		    color: white;
		    border: 1px solid #4CAF50;
		    border-radius: 5px;
		}
		.pagination a:hover:not(.active) {background-color: #ddd;}
		.nav a { text-decoration: none; padding:3px;}
		.nav { display: inline-block; vertical-align: baseline;}
		#navWrap{position:absolute;top:20;right:10;}
	</style>
	<body>
	<div id="navWrap">
	{{ if .Authenticated }}
	<div class="nav"><a href="view_bulletin_board?i=0&pc={{.PC}}">First<a></div>
	<div class="nav"><a href="view_bulletin_board?i={{.PrevI}}&pc={{.PC}}">Prev<a></div>
	<div class="nav"><a href="view_bulletin_board?i={{.NextI}}&pc={{.PC}}">Next<a></div>
	<div class="nav"><a href="remove_bulletin_page?i={{.I}}&pc={{.PC}}">&#x2718;</a></div>
	{{ end }}
	</div>
	<pre class="prettyprint linenums">{{.Code}}</pre>

	<div class="bottom">
	<div class="label">&#128546;</div><div id="p2">{{.P2}}</div>
	<div class="label">&#128526;</div><div id="p1">{{.P1}}</div>
	<div class="label">Problems:</div><div id="ap">{{.ActiveProblems}}</div>
	<div class="label">Bulletin:</div><div id="bu">{{.BulletinItems}}</div>
	<div class="label">Attendance:</div><div id="at">{{.Attendance}}</div>
	</div>
	</body>
</html>
`

var VIEW_ANSWERS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar']});
      google.charts.setOnLoadCallback(drawChart);
      function drawChart() {
        var data = google.visualization.arrayToDataTable([
          ["Answer", "Count", {role: 'annotation'}],
          	{{$total := .Total}}
			{{ range $key, $value := .Counts }}
				[{{ $key }}, {{ $value }}, Math.round(100 * {{$value}} / {{ $total }})  + '%'],
			{{ end }}
		]);
        var options = {
        	'title':'Total votes: {{ .Total }}',
        	'height':300,
        	'legend': {position: "none"},
        	'fontSize': 24,
            vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: "", textStyle: { fontSize: 20} },
        };
        var chart = new google.charts.Bar(document.getElementById('chart_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }
    </script>
    <style>
    #chart_div{ margin: auto; width:70%; }
    pre{ margin: auto; width:60%}
    .spacer{ width:100%; height:40px; }
    </style>
  </head>

  <body>
    <div id="chart_div"></div>
    <div class="spacer"></div>
    <pre id="content">{{ .Content }}</pre>
    <div class="spacer"></div>
  </body>
</html>
`

var TAG_REPORT_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','line']});
      google.charts.setOnLoadCallback(drawSuccess_vs_Activity);
      google.charts.setOnLoadCallback(drawSuccess_vs_Time);
      function drawSuccess_vs_Activity() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Success Rate', 'Activity'],
			{{ range $pid, $rec := .Performance }}
				[ '{{$pid}}', {{$rec.Success}}, {{$rec.Activity}}/100 ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Analyze running time of nested loops.',
        	height: 350,
	        vAxis: { format: 'percent', textStyle: { fontSize: 20} },
            hAxis: { title: 'Problem ID', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart1_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }

      function drawSuccess_vs_Time() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Success Rate'],
			{{ range $pid, $rec := .Performance }}
				[ new Date({{$rec.Timestamp}} / 1000000), {{$rec.Success}} ],
			{{ end }}
	      ]);
	      var options = {
        	height: 350,
	        vAxis: { format: 'percent', textStyle: { fontSize: 20}},
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Line(document.getElementById('chart2_div'));
        chart.draw(data, google.charts.Line.convertOptions(options));
      }
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

var TAGS_VIEW_TEMPLATE = `
<html>
  <head>
    <style>
    body { font-size: 16pt;}
    </style>
  </head>
  <body>
  	<h1>Learning objectives</h1>
  	<ul>
	{{ range . }}
	<li><a href="report_tag?pc={{.PC}}&tag_id={{.Id}}">{{.Description}}</a></li>
	{{ end }}
	</ul>
  </body>
</html>
`

var ACTIVITY_VIEW_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','line']});
      google.charts.setOnLoadCallback(drawStudentCount);
      google.charts.setOnLoadCallback(drawExerciseCount);

      function drawStudentCount() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Student Count'],
			{{ range $day, $val := . }}
				[  new Date({{$day}} / 1000000), {{$val.SidCount}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Daily student participation',
        	height: 350,
	        vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart1_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }

      function drawExerciseCount() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Excercise Count'],
			{{ range $day, $val := . }}
				[  new Date({{$day}} / 1000000), {{$val.PidCount}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Daily exercise',
        	height: 350,
	        vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart2_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }
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
