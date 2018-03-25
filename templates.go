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
		<meta http-equiv="refresh" content="10" />
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
		}
		.label{ display: inline; }
		.p1, .p2, .item {
			padding: 0.75em;
			display: inline;
		}
		.p1 { color: green; }
		.p2 { color: red; }
	</style>
	<body>
	<div class="bottom">
	<div class="label">Submissions:</div>
	<div class="p2"> &#9785; {{.P2}}</div>
	<div class="p1"> &#128526; {{.P1}}</div>
	<div class="label">Active problems:</div>
	<div class="item">{{.ActiveProblems}}</div>
	<div class="label">Attendance:</div>
	<div class="item">{{.Attendance}}</div>
	</div>
	</body>
</html>
`
