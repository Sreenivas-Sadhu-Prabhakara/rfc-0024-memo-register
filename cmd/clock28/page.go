package main

import "html/template"

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Clock accuracy and delay | RFC 28 lab</title>
<style>
body{margin:0;background:#f3f6fa;color:#182c43;font:1rem/1.55 system-ui,sans-serif}
main{max-width:64rem;margin:auto;padding:clamp(1rem,4vw,3rem)}
h1{font-size:clamp(2rem,5vw,3rem);line-height:1.15}
table{border-collapse:collapse;background:white;width:100%;text-align:left}
th,td{padding:.65rem;border:1px solid #b6c8dc}
thead{background:#dfeaf6} .scroll{overflow-x:auto}
a{color:#174d84}
</style>
</head>
<body>
<main>
<p><a href="https://www.rfc-editor.org/info/rfc28/">RFC 28</a> · Clock accuracy lab</p>
<h1>Can these clocks measure the delay?</h1>
<p>Latency budget {{.BudgetUS}} microseconds. Each delay interval uses the sum of the sender and receiver error bounds.</p>
<div class="scroll">
<table>
<thead><tr><th scope="col">Sample</th><th scope="col">Observed µs</th><th scope="col">Error ±µs</th><th scope="col">Possible delay µs</th><th scope="col">Assessment</th></tr></thead>
<tbody>{{range .Results}}
<tr><th scope="row">{{.Name}}</th><td>{{.ObservedUS}}</td><td>{{.ErrorUS}}</td><td>{{if .Consistent}}{{.MinUS}} to {{.MaxUS}}{{else}}Empty interval{{end}}</td><td>{{.Verdict}}</td></tr>{{end}}
</tbody>
</table>
</div>
<p>Bounds are supplied assumptions, with no probability attached. An empty interval means these timestamps and bounds cannot describe a nonnegative delay.</p>
<p><a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0024-memo-register/tree/main/cmd/clock28">Go source, assumptions, and run instructions</a></p>
</main>
</body>
</html>
`))
