package template

// Dockerfile is the Dockerfile template used for new projects.
var Dockerfile = `FROM golang:alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/{{.Service}}{{if .Client}}-client{{end}}

# Install dependencies
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git{{if .PrivateRepo}} openssh-client{{end}}
{{- if .PrivateRepo}}
{{if ne .Vendor ""}}
# Env config for private repo
ENV GOPRIVATE="{{ gitorg .Vendor }}/*"
RUN git config --global url."ssh://git@{{ gitorg .Vendor }}".insteadOf "https://{{ gitorg .Vendor }}" 
{{else}}
# Configure these values
# ENV GOPRIVATE="github.com/<your private org>/*"
# RUN git config --global url."ssh://git@github.com/<your private org>".insteadOf "https://github.com/<your private org>" 
{{end}}
# Authorize SSH Host
RUN mkdir -p /root/.ssh && \
	chmod 0700 /root/.ssh && \
	ssh-keyscan github.com > /root/.ssh/known_hosts &&\
	chmod 644 /root/.ssh/known_hosts && touch /root/.ssh/config
{{end}}
{{if .Health}}
# Download grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.11 && \
wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
chmod +x /bin/grpc_health_probe
{{end}}
# Build Go binary
COPY {{if not .Client}}Makefile {{end}}go.mod go.sum ./
RUN {{if .PrivateRepo}}--mount=type=ssh {{end}}{{if .Buildkit}}--mount=type=cache,mode=0755,target=/go/pkg/mod {{end}}{{if not .Client}}make init && {{end}}go mod download 
COPY . .
RUN {{if .Buildkit}}--mount=type=cache,target=/root/.cache/go-build --mount=type=cache,mode=0755,target=/go/pkg/mod {{end}}make {{if not .Client}}proto {{end}}tidy build

# Deployment container
FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
{{- if .Health}}
COPY --from=build /bin/grpc_health_probe /bin/
{{- end}}
COPY --from=builder /go/src/{{.Service}}{{if .Client}}-client{{end}}/{{.Service}}{{if .Client}}-client{{end}} /{{.Service}}{{if .Client}}-client{{end}}
ENTRYPOINT ["/{{.Service}}{{if .Client}}-client{{end}}"]
CMD []
`

// DockerIgnore is the .dockerignore template used for new projects.
var DockerIgnore = `.gitignore
Dockerfile{{if or .Skaffold .Kustomize}}
resources/
skaffold.yaml{{end}}
`
