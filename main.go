package main

import (
	"github.com/go-micro/cli/cmd"

	// register commands
	_ "github.com/go-micro/cli/cmd/call"
	_ "github.com/go-micro/cli/cmd/describe"
	_ "github.com/go-micro/cli/cmd/generate"
	_ "github.com/go-micro/cli/cmd/new"
	_ "github.com/go-micro/cli/cmd/run"
	_ "github.com/go-micro/cli/cmd/services"
	_ "github.com/go-micro/cli/cmd/stream"
)

func main() {
	cmd.Run()
}
