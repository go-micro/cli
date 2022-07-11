package template

// Module is the go.mod template used for new projects.
var Module = `module {{.Vendor}}{{.Service}}{{if .Client}}-client{{end}}

go 1.18

require (
	go-micro.dev/v4 v4.7.0
)
{{if eq .Vendor ""}}
replace {{lower .Service}} => ./
{{end}}
`
