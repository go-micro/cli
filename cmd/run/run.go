package run

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
	mcli "github.com/go-micro/cli/cmd"
	"go-micro.dev/v4/runtime"
	"go-micro.dev/v4/runtime/local/git"
)

var (
	DefaultRetries = 3

	flags []cli.Flag = []cli.Flag{
		&cli.StringFlag{
			Name:  "command",
			Usage: "command to execute",
		},
		&cli.StringFlag{
			Name:  "args",
			Usage: "command arguments",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "the type of service to operate on",
		},
	}
)

func init() {
	mcli.Register(&cli.Command{
		Name:   "run",
		Usage:  "Build and run a service continuously, e.g. " + mcli.App().Name + " run [github.com/auditemarlow/helloworld]",
		Flags:  flags,
		Action: Run,
	})
}

// Run runs a service and watches the project directory for change events. On
// write, the service is restarted. Exits on error.
func Run(ctx *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	source, err := git.ParseSourceLocal(wd, ctx.Args().First())
	if err != nil {
		return err
	}

	tmpSrc := source.RuntimeSource()
	if os.Getenv("GOOS") == "windows" {
		tmpPath = strings.Replace(source.RuntimeSource(), ":", "", 1)
	}
	
	svc := &runtime.Service{
		Name:     source.RuntimeName(),
		Source:   tmpSrc,
		Version:  source.Ref,
		Metadata: make(map[string]string),
	}

	typ := ctx.String("type")
	command := strings.TrimSpace(ctx.String("command"))
	args := strings.TrimSpace(ctx.String("args"))

	r := *mcli.DefaultCLI.Options().Runtime

	var retries = DefaultRetries
	if ctx.IsSet("retries") {
		retries = ctx.Int("retries")
	}

	opts := []runtime.CreateOption{
		runtime.WithOutput(os.Stdout),
		runtime.WithRetries(retries),
		runtime.CreateType(typ),
	}

	if len(command) > 0 {
		opts = append(opts, runtime.WithCommand(command))
	}

	if len(args) > 0 {
		opts = append(opts, runtime.WithArgs(args))
	}

	if err := r.Create(svc, opts...); err != nil {
		return err
	}

	done := make(chan bool)
	if r.String() == "local" {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		go func() {
			<-sig
			r.Delete(svc)
			done <- true
		}()
	}

	if source.Local {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println(err)
		}
		defer watcher.Close()

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						r.Update(svc)
					}
					if event.Op&fsnotify.Create == fsnotify.Create {
						watcher.Add(event.Name)
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						watcher.Remove(event.Name)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					fmt.Println("ERROR", err)
				}
			}
		}()

		var files []string
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})

		for _, file := range files {
			if err := watcher.Add(file); err != nil {
				return err
			}
		}
	}

	if r.String() == "local" {
		<-done
	}

	return nil
}
