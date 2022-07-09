package new

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	mcli "github.com/go-micro/cli/cmd"
	"github.com/go-micro/cli/generator"
	tmpl "github.com/go-micro/cli/generator/template"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var flags []cli.Flag = []cli.Flag{
	&cli.BoolFlag{
		Name:  "jaeger",
		Usage: "Generate Jaeger tracer files",
	},
	&cli.BoolFlag{
		Name:  "kubernetes",
		Usage: "Generate Kubernetes resource files",
	},
	&cli.BoolFlag{
		Name:  "skaffold",
		Usage: "Generate Skaffold files",
	},
	&cli.BoolFlag{
		Name:  "tilt",
		Usage: "Generate Tiltfile",
	},
	&cli.BoolFlag{
		Name:  "health",
		Usage: "Generate gRPC Health service used for Kubernetes liveliness and readiness probes",
	},
	&cli.BoolFlag{
		Name:  "kustomize",
		Usage: "Generate kubernetes resouce files in a kustomize structure",
	},
	&cli.BoolFlag{
		Name:  "sqlc",
		Usage: "Generate sqlc resources",
	},
	&cli.BoolFlag{
		Name:  "grpc",
		Usage: "Use gRPC as default server and client",
	},
	&cli.BoolFlag{
		Name:  "buildkit",
		Usage: "Use BuildKit features in Dockerfile",
	},
	&cli.BoolFlag{
		Name:  "tern",
		Usage: "Generate tern resouces; sql migrations templates",
	},
	&cli.BoolFlag{
		Name:  "advanced",
		Usage: "Generate advanced features in main.go server file",
	},
	&cli.BoolFlag{
		Name:  "privaterepo",
		Usage: "Amend Dockerfile to build from private repositories (add ssh-agent)",
	},
	&cli.StringFlag{
		Name:  "namespace",
		Usage: "Default namespace for kubernetes resources, defaults to 'default'",
		Value: "default",
	},
	&cli.StringFlag{
		Name:  "postgresaddress",
		Usage: "Default postgres address for kubernetes resources, defaults to postgres.database.svc",
		Value: "postgres.database.svc",
	},
	&cli.BoolFlag{
		Name:  "complete",
		Usage: "Complete will set the following flags to true; jaeger, health, grpc, sqlc, tern, kustomize, tilt, advanced",
	},
}

// NewCommand returns a new new cli command.
func init() {
	mcli.Register(&cli.Command{
		Name:  "new",
		Usage: "Create a project template",
		Subcommands: []*cli.Command{
			{
				Name:   "client",
				Usage:  "Create a client template, e.g. " + mcli.App().Name + " new client [github.com/auditemarlow/]helloworld",
				Action: Client,
				Flags:  flags,
			},
			{
				Name:   "function",
				Usage:  "Create a function template, e.g. " + mcli.App().Name + " new function [github.com/auditemarlow/]helloworld",
				Action: Function,
				Flags:  flags,
			},
			{
				Name:   "service",
				Usage:  "Create a service template, e.g. " + mcli.App().Name + " new service [github.com/auditemarlow/]helloworld",
				Action: Service,
				Flags:  flags,
			},
		},
	})
}

func Client(ctx *cli.Context) error {
	return createProject(ctx, "client")
}

// Function creates a new function project template. Exits on error.
func Function(ctx *cli.Context) error {
	return createProject(ctx, "function")
}

// Service creates a new service project template. Exits on error.
func Service(ctx *cli.Context) error {
	return createProject(ctx, "service")
}

func createProject(ctx *cli.Context, pt string) error {
	arg := ctx.Args().First()
	if len(arg) == 0 {
		return cli.ShowSubcommandHelp(ctx)
	}

	client := pt == "client"
	name, vendor := getNameAndVendor(arg)

	dir := name
	if client {
		dir += "-client"
	}

	if path.IsAbs(dir) {
		fmt.Println("must provide a relative path as service name")
		return nil
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", dir)
	}

	fmt.Printf("creating %s %s\n", pt, name)

	g := generator.New(
		generator.Service(name),
		generator.Vendor(vendor),
		generator.Directory(dir),
		generator.Client(client),
		generator.Jaeger(ctx.Bool("jaeger") || ctx.Bool("complete")),
		generator.Skaffold(ctx.Bool("skaffold")),
		generator.Tilt(ctx.Bool("tilt") || ctx.Bool("complete")),
		generator.Health(ctx.Bool("health") || ctx.Bool("complete")),
		generator.Kustomize(ctx.Bool("kustomize") || ctx.Bool("complete")),
		generator.Sqlc(ctx.Bool("sqlc") || ctx.Bool("complete")),
		generator.GRPC(ctx.Bool("grpc") || ctx.Bool("health") || ctx.Bool("complete")),
		generator.Buildkit(ctx.Bool("buildkit") || ctx.Bool("privaterepo") || ctx.Bool("complete")),
		generator.Tern(ctx.Bool("tern") || ctx.Bool("complete")),
		generator.Advanced(ctx.Bool("advanced") || ctx.Bool("complete")),
		generator.PrivateRepo(ctx.Bool("privaterepo")),
		generator.Namespace(ctx.String("namespace")),
		generator.PostgresAddress(ctx.String("postgresaddress")),
	)

	files := []generator.File{
		{Path: ".dockerignore", Template: tmpl.DockerIgnore},
		{Path: ".gitignore", Template: tmpl.GitIgnore},
		{Path: "Dockerfile", Template: tmpl.Dockerfile},
		{Path: "Makefile", Template: tmpl.Makefile},
		{Path: "go.mod", Template: tmpl.Module},
	}

	switch pt {
	case "client":
		files = append(files, []generator.File{
			{Path: "main.go", Template: tmpl.MainCLT},
		}...)
	case "function":
		files = append(files, []generator.File{
			{Path: "handler/" + name + ".go", Template: tmpl.HandlerFNC},
			{Path: "main.go", Template: tmpl.MainFNC},
			{Path: "proto/" + name + ".proto", Template: tmpl.ProtoFNC},
		}...)
	case "service":
		files = append(files, []generator.File{
			{Path: "handler/" + name + ".go", Template: tmpl.HandlerSRV},
			{Path: "main.go", Template: tmpl.MainSRV},
			{Path: "proto/" + name + ".proto", Template: tmpl.ProtoSRV},
		}...)
	default:
		return fmt.Errorf("%s project type not supported", pt)
	}

	opts := g.Options()
	if opts.Sqlc {
		files = append(files, []generator.File{
			{Path: "postgres/sqlc.yaml", Template: tmpl.Sqlc},
			{Path: "postgres/postgres.go", Template: tmpl.Postgres},
			{Path: "postgres/queries/example.sql", Template: tmpl.QueryExample},
			{Path: "postgres/migrations/", Template: ""},
		}...)
	}

	if opts.Tern {
		files = append(files, []generator.File{
			{Path: "postgres/migrations/001_create_schema.sql", Template: tmpl.TernSql},
		}...)
	}

	if opts.Health {
		files = append(files, []generator.File{
			{Path: "proto/health.proto", Template: tmpl.ProtoHEALTH},
			{Path: "handler/health.go", Template: tmpl.HealthSRV},
		}...)
	}

	if (ctx.Bool("kubernetes") || opts.Skaffold || opts.Tilt) && !opts.Kustomize {
		files = append(files, []generator.File{
			{Path: "plugins.go", Template: tmpl.Plugins},
			{Path: "resources/clusterrole.yaml", Template: tmpl.KubernetesClusterRole},
			{Path: "resources/configmap.yaml", Template: tmpl.KubernetesEnv},
			{Path: "resources/deployment.yaml", Template: tmpl.KubernetesDeployment},
			{Path: "resources/rolebinding.yaml", Template: tmpl.KubernetesRoleBinding},
		}...)
	}

	if opts.Kustomize {
		files = append(files, []generator.File{
			{Path: "plugins.go", Template: tmpl.Plugins},
			{Path: "resources/base/clusterrole.yaml", Template: tmpl.KubernetesClusterRole},
			{Path: "resources/base/app.env", Template: tmpl.AppEnv},
			{Path: "resources/base/deployment.yaml", Template: tmpl.KubernetesDeployment},
			{Path: "resources/base/rolebinding.yaml", Template: tmpl.KubernetesRoleBinding},
			{Path: "resources/base/kustomization.yaml", Template: tmpl.KustomizationBase},
			{Path: "resources/dev/kustomization.yaml", Template: tmpl.KustomizationDev},
			{Path: "resources/prod/kustomization.yaml", Template: tmpl.KustomizationProd},
		}...)
	}

	if opts.Skaffold {
		files = append(files, []generator.File{
			{Path: "skaffold.yaml", Template: tmpl.SkaffoldCFG},
		}...)
	}

	if opts.Tilt {
		files = append(files, []generator.File{
			{Path: "Tiltfile", Template: tmpl.Tiltfile},
		}...)
	}

	if err := g.Generate(files); err != nil {
		return err
	}

	var comments []string
	if client {
		comments = clientComments(name, dir)
	} else {
		var err error
		comments, err = protoComments(name, dir, opts.Sqlc)
		if err != nil {
			return err
		}
	}

	for _, comment := range comments {
		fmt.Println(comment)
	}

	return nil
}

func clientComments(name, dir string) []string {
	return []string{
		"\ninstall dependencies:",
		"\ncd " + dir,
		"make update tidy",
	}
}

func protoComments(name, dir string, sqlc bool) ([]string, error) {
	tmp := `
download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases/latest

download protobuf for go-micro:

go get -u google.golang.org/protobuf/proto
go install github.com/golang/protobuf/protoc-gen-go@latest
go install go-micro.dev/v4/cmd/protoc-gen-micro@v4

compile the proto file {{ .Name }}.proto and install dependencies:

cd {{ .Dir }}
make proto {{ if .Sqlc }}sqlc {{ end }}update tidy`

	t, err := template.New("comments").Parse(tmp)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse comments template")
	}

	var b bytes.Buffer
	if err := t.Execute(&b, map[string]interface{}{
		"Name": name,
		"Dir":  dir,
		"Sqlc": sqlc,
	}); err != nil {
		return nil, errors.Wrap(err, "Failed to execute proto comments template")
	}
	return []string{b.String()}, nil
}

func getNameAndVendor(s string) (string, string) {
	var n string
	var v string

	if i := strings.LastIndex(s, "/"); i == -1 {
		n = s
		v = ""
	} else {
		n = s[i+1:]
		v = s[:i+1]
	}

	return n, v
}
