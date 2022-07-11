package template

// Plugins is the plugins template used for new projects.
var Plugins = `package main

import (
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
)
`
