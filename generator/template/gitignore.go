package template

// GitIgnore is the .gitignore template used for new projects.
var GitIgnore = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# don't commit the service binary to vcs
{{.Service}}{{if .Client}}-client{{end}}
`
