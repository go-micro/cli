package main

import (
	"github.com/go-micro/go-micro/cmd"

	// register commands
	_ "github.com/go-micro/go-micro/cmd/call"
	_ "github.com/go-micro/go-micro/cmd/describe"
	_ "github.com/go-micro/go-micro/cmd/generate"
	_ "github.com/go-micro/go-micro/cmd/new"
	_ "github.com/go-micro/go-micro/cmd/run"
	_ "github.com/go-micro/go-micro/cmd/services"
	_ "github.com/go-micro/go-micro/cmd/stream"

	// plugins
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
)

func main() {
	cmd.Run()
}
