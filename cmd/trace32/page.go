package main

import "html/template"

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Host delays and clock offsets | RFC 32 lab</title>
<style>
body{margin:0;background:#f3f6fa;color:#182c43;font:1rem/1.55 system-ui,sans-serif}
main{max-width:68rem;margin:auto;padding:clamp(1rem,4vw,3rem)}
h1{font-size:clamp(2rem,5vw,3rem);line-height:1.15}
table{border-collapse:collapse;background:white;width:100%;text-align:left}
th,td{padding:.65rem;border:1px solid #b6c8dc}thead{background:#dfeaf6}
.scroll{overflow-x:auto}a{color:#174d84}
</style>
</head>
<body><main>
<p><a href="https://www.rfc-editor.org/info/rfc32/">RFC 32</a> · Host tracing lab</p>
<h1>Where does the delay happen?</h1>
<p>Actual network transit {{.ActualNetworkUS}} µs · Actual total delay {{.ActualTotalUS}} µs</p>
<div class="scroll"><table>
<thead><tr><th scope="col">Clock pair</th><th scope="col">Local queue µs</th><th scope="col">Local service µs</th><th scope="col">Raw transit µs</th><th scope="col">Raw total µs</th></tr></thead>
<tbody>{{range .Traces}}
<tr><th scope="row">{{.Name}}</th><td>{{.LocalQueueUS}}</td><td>{{.LocalServiceUS}}</td><td>{{.RawTransitUS}}</td><td>{{.RawTotalUS}}</td></tr>{{end}}
</tbody>
</table></div>
<p>Each request follows the same path. Constant clock offsets cancel within a host. Subtracting timestamps from different hosts adds the difference between their offsets.</p>
<p>This synthetic model uses perfect clock rates and exact timestamps. Negative raw results expose the clock mismatch.</p>
<p><a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0024-memo-register/tree/main/cmd/trace32">Go source, assumptions, and run instructions</a></p>
</main></body></html>
`))
