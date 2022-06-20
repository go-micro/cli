# Go Micro CLI

Go Micro CLI is the command line interface for developing [Go Micro][1] projects.

## Getting Started

[Download][2] and install **Go**. Version `1.16` or higher is required.

Installation is done by using the [`go install`][3] command.

```bash
go install github.com/go-micro/cli/cmd/go-micro@v1.1.1
```

Let's create a new service using the `new` command.

```bash
go-micro new service helloworld
```

Follow the on-screen instructions. Next, we can run the program.

```bash
cd helloworld
make proto tidy
go-micro run
```

Finally we can call the service.

```bash
go-micro call helloworld Helloworld.Call '{"name": "John"}'
```

That's all you need to know to get started. Refer to the [Go Micro][1]
documentation for more info on developing services.

## Dependencies

You will need protoc-gen-micro for code generation

* [protobuf][4]
* [protoc-gen-go][5]
* [protoc-gen-micro][6]

```bash
# Download latest proto release
# https://github.com/protocolbuffers/protobuf/releases
go get -u google.golang.org/protobuf/proto
go install github.com/golang/protobuf/protoc-gen-go@latest
go install go-micro.dev/v4/cmd/protoc-gen-micro@v4
```

## Creating A Service

To create a new service, use the `micro new service` command, and provide either a bare
service name, or a full GitHub repo module name.

```bash
$ go-micro new service github.com/<org>/<repo>/helloworld
...
```

```bash
$ go-micro new service helloworld
creating service helloworld

download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases/latest

download protobuf for go-micro:

go get -u google.golang.org/protobuf/proto
go install github.com/golang/protobuf/protoc-gen-go@latest
go install go-micro.dev/cmd/protoc-gen-micro/v4@latest

compile the proto file helloworld.proto:

cd helloworld
make proto tidy
```

To create a new function, use the `micro new function` command. Functions differ
from services in that they exit after returning.

```bash
$ go-micro new function helloworld
creating function helloworld

download protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:

visit https://github.com/protocolbuffers/protobuf/releases/latest

download protobuf for go-micro:

go get -u google.golang.org/protobuf/proto
go install github.com/golang/protobuf/protoc-gen-go@latest
go install go-micro.dev/cmd/protoc-gen-micro/v4@latest

compile the proto file helloworld.proto:

cd helloworld
make proto tidy
```

### Jaeger

To create a new service with [Jaeger][7] integration, pass the `--jaeger` flag
to the `micro new service` or `micro new function` commands. You may configure
the Jaeger client using [environment variables][8].

```bash
go-micro new service --jaeger helloworld
```

You may invoke `trace.NewSpan(context.Context).Finish()` to nest spans. For
example, consider the following handler implementing a greeter.

`handler/helloworld.go`

```go
package helloworld

import (
    "context"

    log "go-micro.dev/v4/logger"

    "helloworld/greeter"
    pb "helloworld/proto"
)

type Helloworld struct{}

func (e *Helloworld) Call(ctx context.Context, req pb.CallRequest, rsp *pb.CallResponse) error {
    log.Infof("Received Helloworld.Call request: %v", req)
    rsp.Msg = greeter.Greet(ctx, req.Name)
    return nil
}
```

`greeter/greeter.go`

```go
package greeter

import (
    "context"
    "fmt"

    "go-micro.dev/v4/cmd/micro/debug/trace"
)

func Greet(ctx context.Context, name string) string {
    defer trace.NewSpan(ctx).Finish()
    return fmt.Sprint("Hello " + name)
}
```

### gRPC Server/Client

By default, go-micro uses an JSON/HTTP RPC server. Many microservice use
cases require a gRPC server or client, therefore, go-micro offers a gRPC server 
built in.

To create a new service with a gRPC server pass the `--grpc` flag to
the `micro new service` or `micro new function` commands.

```bash
go-micro new service --grpc helloworld
```

### Tern - Postgres Migrations

Tern can be used to create and manage Postgres migrations. Go-micro can set the
service up to use Tern SQL migrations.

To create a new service with Tern pass the `--tern` flag to
the `micro new service` or `micro new function` commands.

To locally run `tern migrate` it is recommended to create a `.env` file with
connection details as shown below, and manually run `source .env`, as these environment variables
will also be picked up by pgx.

```env
PGHOST=localhost
PGUSER=helloworld
PGDATABASE=helloworld
PGPASSWORD=<empty for localhost>
```

Setting the `--tern` flag in combination with any of the Kubernetes flags will
also create a `InitContainer` in the deployment manifest to automatically apply
all migrations upon deployment. You will need to create a `helloworld-postgres-env`
secret with a `PGPASSWORD` to allow tern to connect to your database instance.

The default database address used is `postgres.database.svc`, and the user and 
database are set to the service name. To specify a different database address
use the `--postgresaddress="my.namespace.svc"` flag

If you are also using Kustomize to manage your Kubernetes resources, be aware
that you will have to manually add every migration to `resouces/base/kustomization.yaml`.
And that you will have to pass the `--load-restrictor LoadRestrictionsNone` flag 
to `kustomize build`, to allow Kustomize to access resources outside of the `base` 
folder. Tilt and Skaffold will do this for you automatically.

```bash
go-micro new service --tern helloworld
```

### sqlc - SQL Code Generation

[Sqlc](14) can compile SQL queries into boilerplate Go code that allows you to easily
create and manage your database layer. Go-micro can set your service up for use
with sqlc, and used Postgres as a default backend. Sqlc works well in combination 
with [Tern](#tern---postgres-migrations). 

Place your SQL queries in `postgres/queries/*.sql` and run `make sqlc` to compile.
Be sure you have your SQL schema defined in `postgres/migrations/*.sql`, as can 
be done with Tern.

After compilation, you can create your database layer in `postgres/*.go` with the
sqlc connector. An example is provided in `postgres/postgres.go`.

To create a new service with sqlc pass the `--sqlc` flag to
the `micro new service` or `micro new function` commands.

```bash
go-micro new service --sqlc helloworld
```

### Docker BuildKit

[Docker BuildKit](11) is a new container build engine that provides new useful 
features, such as the ability to cache specific directories across builds. This
can prevent Go from having to re-download modules every build.

To create a new service with the BuildKit engine pass the `--buildkit` flag to
the `micro new service` or `micro new function` commands.

```bash
go-micro new service --buildkit helloworld
```

### Private Git Repository

If you plan on hosting the service in a private Git repository, the docker file
needs some tweaks to allow Go to access and clone private repositories.

For this, SSH Git access needs to be set up, and an SSH agent needs to be running,
with the Git SSH key added. You can manually start an SSH agent and add the SSH key
by running:

```bash
$ eval $(ssh-agent) && ssh-add <optional: path to ssh key, default: ~/.shh/id_rsa>
```

Alternatively, you can use the generated [Tiltfile](#tilt---kubernets-deployment) 
to let Tilt set one up for you.

To create a new service with a private Git repository pass the `--privaterepo` flag to
the `micro new service` or `micro new function` commands. This implies the `--buildkit`
flag.

```bash
go-micro new service --privaterepo helloworld
```

### Kubernetes

Micro can automatically generate Kubernetes manifests for a service template.

To create a new service with Kubernetes resources pass the `--kubernetes` flag to
the `micro new service` or `micro new function` commands.

```bash
go-micro new service --kubernetes helloworld
```

### Kustomize - Kubernetes Resource Management

[Kustomize](15) can be used to manage more complex Kubernetes manifests for various
deployments, such as a development and production environment.

To create a new service with Kubernetes resources organized in a Kustomize structure
pass the `--kustomize` flag to the `micro new service` or `micro new function` 
commands.

```bash
go-micro new service --kustomize helloworld
```

### gRPC Health Protocol - Kubernetes Probes

Since Kubernetes [1.24](12), probes can make use of the [gRPC Health Protocol](13).
This allows you to directly probe the go-micro service in a Kubernetes container
if it implements the health protocol.

By passing the `--health` flag the gRPC protocol will be implemented, and if 
Kubernetes manifests are generated through any of the flags, it will add probes
to the deployment manifest.

To use this feature, the `GRPCContainerProbe` feature gate needs to be enabled 
inside your cluster. In version 1.24 this is enabled by default, in version 1.23
you need to manually enable the feature gate.

To create a new service with the gRPC health protocol implemented, pass the 
`--health` flag to the `micro new service` or `micro new function` commands. This
implies the `--grpc` flag.

```bash
go-micro new service --health helloworld
```

### Kubernetes Options

#### Namespace

Kubernetes manifests and Kustomize files set an explicit namespace by default.
The default Kubernetes namespace is `default`. You can manually specify a different 
namespace during service creation.

```bash
go-micro new service --kustomize --namespace=custom helloworld
```

#### Postgres Address

If you create the service with the `--tern` flag, the default Postgres address
used is `postgres.database.svc`. To specify a different address use the 
`--postgresaddress` flag

```bash
go-micro new service --kustomize --tern --postgresaddress="my.namespace.svc" helloworld
```

### Tilt - Kubernetes Deployment

Tilt can be used to set up a local Kubernetes deployment pipeline.

To create a new service with a [Tiltfile][9] file, pass the `--tilt` flag to
the `micro new service` or `micro new function` commands.

This implies the `--kubernetes` flag.

```bash
go-micro new service --tilt helloworld
```

### Skaffold - Kubernetes Deployment

Skaffold can be used to locally deploy a service.

To create a new service with [Skaffold][9] files, pass the `--skaffold` flag to
the `micro new service` or `micro new function` commands.

This implies the `--kubernetes` flag.

```bash
go-micro new service --skaffold helloworld
```

### Advanced

Some patterns will often occur in more complex services. Such as the need to 
gracefully shutdown go routines, and pass down a context to provide the cancelation  
signal.

To prevent you from having to rewrite them for every service, you can pass the 
`--advanced` flag. This will generate a waitgroup, context, and define functions
for `BeforeStart`, `BeforeStop` and `AfterStop`.

```bash
go-micro new service --advanced helloworld
```

### Complete

With so many possible flags to create a service, the `--complete` flag will
set the following flags to true:

```bash
go-micro new service --jaeger --health --grpc --sqlc --tern --buildkit --kustomize --tilt --advanced
```

## Running A Service

To run a service, use the `micro run` command to build and run your service
continuously.

```bash
$ go-micro run
2021-08-20 14:05:54  file=v3@v3.5.2/service.go:199 level=info Starting [service] helloworld
2021-08-20 14:05:54  file=server/rpc_server.go:820 level=info Transport [http] Listening on [::]:34531
2021-08-20 14:05:54  file=server/rpc_server.go:840 level=info Broker [http] Connected to 127.0.0.1:44975
2021-08-20 14:05:54  file=server/rpc_server.go:654 level=info Registry [mdns] Registering node: helloworld-45f43a6f-5fc0-4b0d-af73-e4a10c36ef54
```

### With Docker

To run a service with Docker, build the Docker image and run the Docker
container.

```bash
$ make docker
$ docker run helloworld:latest
2021-08-20 12:07:31  file=v3@v3.5.2/service.go:199 level=info Starting [service] helloworld
2021-08-20 12:07:31  file=server/rpc_server.go:820 level=info Transport [http] Listening on [::]:36037
2021-08-20 12:07:31  file=server/rpc_server.go:840 level=info Broker [http] Connected to 127.0.0.1:46157
2021-08-20 12:07:31  file=server/rpc_server.go:654 level=info Registry [mdns] Registering node: helloworld-31f58714-72f5-4d12-b2eb-98f66aea7a34
```

### With Skaffold

When you've created your service using the `--skaffold` flag, you may run the
Skaffold pipeline using the `skaffold` command.

```bash
skaffold dev
```

### With Tilt

When you've created your service using the `--tilt` flag, you may run the
Tilt pipeline using the `tilt` command.

```bash
tilt up --stream
```

If you don't want to stream logs, but do want to exit on errors, you can run

```bash
tilt ci
```


## Creating A Client

To create a new client, use the `micro new client` command. The name is the
service you'd like to create a client project for.

```bash
$ go-micro new client helloworld
creating client helloworld
cd helloworld-client
make tidy
```

You may optionally pass the fully qualified package name of the service you'd
like to create a client project for.

```bash
$ go-micro new client github.com/auditemarlow/helloworld
creating client helloworld
cd helloworld-client
make tidy
```

## Running A Client

To run a client, use the `micro run` command to build and run your client
continuously.

```bash
$ go-micro run
2021-09-03 12:52:23  file=helloworld-client/main.go:33 level=info msg:"Hello John"
```

## Generating Files

To generate Go Micro project template files after the fact, use the `micro
generate` command. It will place the generated files in the current working
directory.

```bash
$ go-micro generate skaffold
skaffold project template files generated
```

## Listing Services

To list services, use the `micro services` command.

```bash
$ go-micro services
helloworld
```

## Describing A Service

To describe a service, use the `micro describe service` command.

```bash
$ go-micro describe service helloworld
{
  "name": "helloworld",
  "version": "latest",
  "metadata": null,
  "endpoints": [
    {
      "name": "Helloworld.Call",
      "request": {
        "name": "CallRequest",
        "type": "CallRequest",
        "values": [
          {
            "name": "name",
            "type": "string",
            "values": null
          }
        ]
      },
      "response": {
        "name": "CallResponse",
        "type": "CallResponse",
        "values": [
          {
            "name": "msg",
            "type": "string",
            "values": null
          }
        ]
      }
    }
  ],
  "nodes": [
    {
      "id": "helloworld-9660f06a-d608-43d9-9f44-e264ff63c554",
      "address": "172.26.165.161:45059",
      "metadata": {
        "broker": "http",
        "protocol": "mucp",
        "registry": "mdns",
        "server": "mucp",
        "transport": "http"
      }
    }
  ]
}
```

You may pass the `--format=yaml` flag to output a YAML formatted object.

```bash
$ go-micro describe service --format=yaml helloworld
name: helloworld
version: latest
metadata: {}
endpoints:
- name: Helloworld.Call
  request:
    name: CallRequest
    type: CallRequest
    values:
    - name: name
      type: string
      values: []
  response:
    name: CallResponse
    type: CallResponse
    values:
    - name: msg
      type: string
      values: []
nodes:
- id: helloworld-9660f06a-d608-43d9-9f44-e264ff63c554
  address: 172.26.165.161:45059
  metadata:
    broker: http
    protocol: mucp
    registry: mdns
    server: mucp
    transport: http
```

## Calling A Service

To call a service, use the `micro call` command. This will send a single request
and expect a single response.

```bash
$ go-micro call helloworld Helloworld.Call '{"name": "John"}'
{"msg":"Hello John"}
```

To call a service's server stream, use the `micro stream server` command. This
will send a single request and expect a stream of responses.

```bash
$ go-micro stream server helloworld Helloworld.ServerStream '{"count": 10}'
{"count":0}
{"count":1}
{"count":2}
{"count":3}
{"count":4}
{"count":5}
{"count":6}
{"count":7}
{"count":8}
{"count":9}
```

To call a service's bidirectional stream, use the `micro stream bidi` command.
This will send a stream of requests and expect a stream of responses.

```bash
$ go-micro stream bidi helloworld Helloworld.BidiStream '{"stroke": 1}' '{"stroke": 2}' '{"stroke": 3}'
{"stroke":1}
{"stroke":2}
{"stroke":3}
```

[1]: https://go-micro.dev
[2]: https://golang.org/dl/
[3]: https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies
[4]: https://grpc.io/docs/protoc-installation/
[5]: https://micro.mu/github.com/golang/protobuf/protoc-gen-go
[6]: https://go-micro.dev/tree/master/cmd/protoc-gen-micro
[7]: https://www.jaegertracing.io/
[8]: https://github.com/jaegertracing/jaeger-client-go#environment-variables
[9]: https://skaffold.dev/
[10]: https://docs.tilt.dev/
[11]: https://docs.docker.com/develop/develop-images/build_enhancements/
[12]: https://kubernetes.io/blog/2022/05/13/grpc-probes-now-in-beta/
[13]: https://github.com/grpc/grpc/blob/master/doc/health-checking.md
[14]: https://github.com/kyleconroy/sqlc
[15]: https://kubectl.docs.kubernetes.io/references/kustomize/
