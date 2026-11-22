package main

import "html/template"

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Symbolic fields to binary | RFC 31 lab</title>
<style>
body{margin:0;background:#f3f6fa;color:#182c43;font:1rem/1.55 system-ui,sans-serif}
main{max-width:64rem;margin:auto;padding:clamp(1rem,4vw,3rem)}
h1{font-size:clamp(2rem,5vw,3rem);line-height:1.15}
table{border-collapse:collapse;background:white;width:100%;text-align:left}
th,td{padding:.65rem;border:1px solid #b6c8dc}
thead{background:#dfeaf6}.scroll{overflow-x:auto}code{overflow-wrap:anywhere}a{color:#174d84}
</style>
</head>
<body>
<main>
<p><a href="https://www.rfc-editor.org/info/rfc31/">RFC 31</a> · Binary field lab</p>
<h1>From named fields to bits</h1>
<p>{{.BitCount}} meaningful bits · Packed bytes <code>{{.Hex}}</code></p>
<p>Stream order <code>{{.Stream}}</code></p>
<div class="scroll"><table>
<thead><tr><th scope="col">Field</th><th scope="col">Bit offset</th><th scope="col">Width</th><th scope="col">Value</th><th scope="col">Field stream</th><th scope="col">Decoded</th></tr></thead>
<tbody>{{range .Fields}}
<tr><th scope="row">{{.Name}}</th><td>{{.Offset}}</td><td>{{.Bits}}</td><td>{{.Value}}</td><td><code>{{.Stream}}</code></td><td>{{.Decoded}}</td></tr>{{end}}
</tbody>
</table></div>
<p>Each field emits its least-significant bit first. Stream positions fill byte bits from zero through seven. Unused bits in the final byte are zero.</p>
<p><a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0024-memo-register/tree/main/cmd/fields31">Go source, model scope, and run instructions</a></p>
</main>
</body>
</html>
`))
