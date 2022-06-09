package template

// MainCLT is the main template used for new client projects.
var MainCLT = `package main

import (
	"context"
	"time"

	pb "{{.Vendor}}{{lower .Service}}/proto"

	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
{{if .GRPC}}
	"github.com/asim/go-micro/plugins/client/grpc/v4"
{{end}}
)

var (
	service = "{{lower .Service}}"
	version = "latest"
)

func main() {
	// Create service
	{{if .GRPC}}
	srv := micro.NewService(
		micro.Client(grpc.NewClient()),
	)
	{{else}}
	srv := micro.NewService()
	{{end}}
	srv.Init()

	// Create client
	c := pb.NewHelloworldService(service, srv.Client())

	for {
		// Call service
		rsp, err := c.Call(context.Background(), &pb.CallRequest{Name: "John"})
		if err != nil {
			log.Fatal(err)
		}

		log.Info(rsp)

		time.Sleep(1 * time.Second)
	}
}
`

// MainFNC is the main template used for new function projects.
var MainFNC = `package main

import (
	"{{.Vendor}}{{.Service}}/handler"

{{if .Jaeger}}	ot "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
{{end}}	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"{{if .Jaeger}}

	"go-micro.dev/v4/cmd/micro/debug/trace/jaeger"{{end}}
)

var (
	service = "{{lower .Service}}"
	version = "latest"
)

func main() {
{{if .Jaeger}}	// Create tracer
	tracer, closer, err := jaeger.NewTracer(
		jaeger.Name(service),
		jaeger.FromEnv(true),
		jaeger.GlobalTracer(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

{{end}}	// Create function
	fnc := micro.NewFunction(
		micro.Name(service),
		micro.Version(version),
{{if .Jaeger}}		micro.WrapCall(ot.NewCallWrapper(tracer)),
		micro.WrapClient(ot.NewClientWrapper(tracer)),
		micro.WrapHandler(ot.NewHandlerWrapper(tracer)),
		micro.WrapSubscriber(ot.NewSubscriberWrapper(tracer)),
{{end}}	)
	fnc.Init()

	// Handle function
	fnc.Handle(new(handler.{{title .Service}}))

	// Run function
	if err := fnc.Run(); err != nil {
		log.Fatal(err)
	}
}
`

// MainSRV is the main template used for new service projects.
var MainSRV = `package main

import (
{{- if .Advanced}}
	"context"
	"sync"
{{- end}}

	"{{.Vendor}}{{.Service}}/handler"
	pb "{{.Vendor}}{{.Service}}/proto"

{{if .Jaeger}}	ot "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
{{end}}	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"{{if .Jaeger}}
{{- if .Advanced}}
	"go-micro.dev/v4/server"
{{- end}}

	"github.com/go-micro/go-micro/debug/trace/jaeger"{{end}}
{{if .GRPC}}
	grpcc "github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
{{- end}}
)

var (
	service = "{{lower .Service}}"
	version = "latest"
)

func main() {
{{if .Jaeger}}	// Create tracer
	tracer, closer, err := jaeger.NewTracer(
		jaeger.Name(service),
		jaeger.FromEnv(true),
		jaeger.GlobalTracer(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
{{ if .Advanced }}
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
{{- end }}

{{end}}	// Create service
	srv := micro.NewService(
{{- if .GRPC}}
		micro.Server(grpc.NewServer()),
		micro.Client(grpcc.NewClient()),
{{- end}}
		micro.Name(service),
		micro.Version(version),
{{- if .Advanced}}
		micro.BeforeStart(func() error {
			log.Infof("Starting service %s", service)
			return nil
		}),
		micro.BeforeStop(func() error {
			log.Infof("Shutting down service %s", service)
			cancel()
			return nil
		}),
		micro.AfterStop(func() error {
			wg.Wait()
			return nil
		}),
{{- end}}
{{if .Jaeger}}		micro.WrapCall(ot.NewCallWrapper(tracer)),
		micro.WrapClient(ot.NewClientWrapper(tracer)),
		micro.WrapHandler(ot.NewHandlerWrapper(tracer)),
		micro.WrapSubscriber(ot.NewSubscriberWrapper(tracer)),
{{end}}	)
	srv.Init()
{{- if .Advanced}}
	srv.Server().Init(
		server.Wait(&wg),
	)

	ctx = server.NewContext(ctx, srv.Server())
{{- end}}

	// Register handler
	pb.Register{{title .Service}}Handler(srv.Server(), new(handler.{{title .Service}}))
{{- if .Health}}
	pb.RegisterHealthHandler(srv.Server(), new(handler.Health))
{{end}}
	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
`
