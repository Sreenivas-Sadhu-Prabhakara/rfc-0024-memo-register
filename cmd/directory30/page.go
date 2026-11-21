package main

import "html/template"

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Memo directory checker | RFC 30 lab</title>
<style>
body{margin:0;background:#f3f6fa;color:#182c43;font:1rem/1.55 system-ui,sans-serif}
main{max-width:68rem;margin:auto;padding:clamp(1rem,4vw,3rem)}
h1{font-size:clamp(2rem,5vw,3rem);line-height:1.15}
table{border-collapse:collapse;background:white;width:100%;text-align:left}
th,td{padding:.65rem;border:1px solid #b6c8dc;overflow-wrap:anywhere}
thead{background:#dfeaf6}.scroll{overflow-x:auto}a{color:#174d84}
</style>
</head>
<body>
<main>
<p><a href="https://www.rfc-editor.org/info/rfc30/">RFC 30</a> · Memo directory checker</p>
<h1>Can the notice reach the listed contacts?</h1>
<p>{{.Ready}} ready routes · {{.Unresolved}} unresolved roster entries</p>
<div class="scroll"><table>
<thead><tr><th scope="col">Site</th><th scope="col">Expected contact</th><th scope="col">Directory contact</th><th scope="col">Ready mailbox</th><th scope="col">Status</th></tr></thead>
<tbody>{{range .Routes}}
<tr><th scope="row">{{.Site}}</th><td>{{.Expected}}</td><td>{{if .Listed}}{{.Listed}}{{else}}—{{end}}</td><td>{{if .Mailbox}}{{.Mailbox}}{{else}}—{{end}}</td><td>{{.Status}}</td></tr>{{end}}
</tbody>
</table></div>
<h2>Directory sites outside the current roster</h2>
<ul>{{range .Unused}}<li>{{.}}</li>{{else}}<li>None</li>{{end}}</ul>
<p>This offline check compares the supplied records. Mailbox reachability and contact ownership require separate confirmation.</p>
<p><a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0024-memo-register/tree/main/cmd/directory30">Go source and run instructions</a></p>
</main>
</body>
</html>
`))
