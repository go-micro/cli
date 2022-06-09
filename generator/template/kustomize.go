package template

var KustomizationBase = `---

namespace: {{ .Namespace }}

resources:
  - clusterrole.yaml
  - deployment.yaml
  - rolebinding.yaml

configMapGenerator:
  - name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-env
    envs:
      - app.env
{{- if .Tern}}
  - name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-migrations
    files:
      - ../../postgres/migrations/001_create_schema.sql
{{end}}
`
var KustomizationDev = `---

namespace: {{ .Namespace }}

resources:
 - ../base/
`

var KustomizationProd = `---

namespace: {{ .Namespace }}

resources:
 - ../base/
`

var AppEnv = `MICRO_REGISTRY=kubernetes
{{- if .Tern}}
PGHOST={{ .PostgresAddress }}
PGUSER={{lowerhyphen .Service}}{{if .Client}}_client{{end}}
PGDATABASE={{lowerhyphen .Service}}{{if .Client}}_client{{end}}
{{end}}
`
