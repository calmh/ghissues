// Copyright 2015 Jakob Borg

package main

import "html/template"

var indexTpl = template.Must(template.New("index").Parse(`<html>
<body>
	<ul>
	{{range $issue := .issues}}
		<li><a href="issue-{{$issue.Number}}.html">#{{$issue.Number}} - {{$issue.Title}}</a>
		{{range $label := $issue.Labels}}
			{{$label.Name}}
		{{end}}
		</li>
	{{end}}
	</ul>
</body>
</html>`))

var issueTpl = template.Must(template.New("issue").Parse(`<html>
<body>
<h1>#{{.Number}} - {{.Title}}</h1>
<pre>
{{.Body}}

// @{{.User.Login}}
</pre>
</body>
</html>`))
