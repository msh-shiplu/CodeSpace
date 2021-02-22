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
					$("#p1u").html(data["P1Ungraded"]);
					$("#p1g").html(data["P1Graded"]);
					$("#p2u").html(data["P2Unanswered"]);
					$("#p2a").html(data["P2Answered"]);
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
		.bottom_left {
			position: fixed;
			bottom: 0;
			text-align: left,
			width: 50%;
		}
		.bottom_right {
			position: fixed;
			bottom: 0;
			right: 50px;
			text-align: right,
			width: 50%;
		}
		.label{ display: inline; }
		#p1, #p2, #p1g, #p1u, #p2a, #p2u, #ans, #ap, #bu, #at {
			padding: 0.75em;
			display: inline;
		}
		#p1g, #p2a { color: green; }
		#p1u, #p2u { color: red; }
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

	<div class="bottom_left">
		<div id="p2">{{.P2}}</div> <div class="label">Help Requests:</div>
		<div id="p2u">{{.P2Unanswered}}</div> <div class="label">Pending,</div>
		<div id="p2a">{{.P2Answered}}</div> <div class="label">Answered</div>
	</div>
	<div class="bottom_right">
		<div id="p1">{{.P1}}</div> <div class="label">Submissions:</div>
		<div id="p1u">{{.P1Ungraded}}</div> <div class="label">Pending,</div>
		<div id="p1g">{{.P1Graded}}</div> <div class="label">Graded</div>
	</div>
	</body>
</html>
`
