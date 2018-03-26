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
		<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js?autoload=true&skin=desert"></script>
		<meta http-equiv="refresh" content="10" />
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			text-align: center;
			width: 100%;
		}
		.label{ display: inline; }
		.p1, .p2, .item {
			padding: 0.75em;
			display: inline;
		}
		.p1 { color: green; }
		.p2 { color: red; }
		pre {
			font-family: monospace;
			font-size:110%;
			padding:1em;
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
		.remove { text-align:right; padding-right:10px; margin-bottom:-15px; }
		.remove a { text-decoration: none; }
	</style>
	<body>
	<div class="center">
	<div class="pagination">
	{{range $i, $a := .Idx}}
	    <a href="view_bulletin_board?i={{$i}}&pc={{$.PC}}" class="{{$a}}">{{inc $i}}</a>
	{{end}}
	</div>
	</div>
	{{ if .Authenticated }}
	<div class="remove"><a href="remove_bulletin_page?i={{.Tbr}}&pc={{.PC}}">&#x2718;</a></div>
	{{ end }}
	<pre class="prettyprint">
	{{.Code}}
	</pre>

	<div class="bottom">
	<div class="label">Submissions:</div>
	<div class="p2"> &#128546; {{.P2}}</div>
	<div class="p1"> &#128526; {{.P1}}</div>
	<div class="label">Active problems:</div>
	<div class="item">{{.ActiveProblems}}</div>
	<div class="label">Attendance:</div>
	<div class="item">{{.Attendance}}</div>
	</div>
	</body>
</html>
`
