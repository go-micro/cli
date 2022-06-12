package completion

import (
	"bytes"
	"fmt"
	"text/template"

	mcli "github.com/go-micro/cli/cmd"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	mcli.Register(&cli.Command{
		Name:  "completion",
		Usage: "Output shell completion code for the specified shell (bash or zsh)",
		Subcommands: []*cli.Command{
			{
				Name:   "bash",
				Usage:  "Create completion script for bash shell. Usage: [[ /sbin/" + mcli.App().Name + " ]] && source <(" + mcli.App().Name + " completion bash)",
				Action: BashCompletion,
			},
			{
				Name:   "zsh",
				Usage:  "Create completion script for zsh shell. Usage: [[ /sbin/" + mcli.App().Name + " ]] && source <(" + mcli.App().Name + " completion zsh)",
				Action: ZshCompletion,
			},
		},
	})
}

func ZshCompletion(ctx *cli.Context) error {
	return renderTemplate(zshTemplate)
}

func BashCompletion(ctx *cli.Context) error {
	return renderTemplate(bashTemplate)
}

func renderTemplate(t string) error {
	tmpl, err := template.New("completionTemplate").Parse(t)
	if err != nil {
		return errors.Wrap(err, "Failed to parse completion template")
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, map[string]interface{}{
		"Prog": mcli.App().Name,
	}); err != nil {
		return errors.Wrap(err, "Failed to render completion template")
	}

	fmt.Println(b.String())

	return nil
}
