package template

// KubernetesEnv is a Kubernetes configmap manifest template used for
// environment variables in new projects.
var KubernetesEnv = `---

apiVersion: v1
kind: ConfigMap
metadata:
  name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-env
data:
  MICRO_REGISTRY: kubernetes
{{if .Tern}}
  PGHOST: {{ .PostgresAddress }}
	PGUSER: {{lowerhyphen .Service}}{{if .Client}}_client{{end}}
	PGDATABASE: {{lowerhyphen .Service}}{{if .Client}}_client{{end}}
{{end}}
`

// KubernetesClusterRole is a Kubernetes cluster role manifest template
// required for the Kubernetes registry plugin to function correctly.
var KubernetesClusterRole = `---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: micro-registry
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - list
  - patch
  - watch
`

// KubernetesRoleBinding is a Kubernetes role binding manifest template
// required for the Kubernetes registry plugin to function correctly.
var KubernetesRoleBinding = `---

apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: micro-registry
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: micro-registry
subjects:
- kind: ServiceAccount
  name: default
  namespace: {{ .Namespace }}
`

// KubernetesDeployment is a Kubernetes deployment manifest template used for
// new projects.
var KubernetesDeployment = `---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{tohyphen .Service}}{{if .Client}}-client{{end}}
  labels:
    app: {{tohyphen .Service}}{{if .Client}}-client{{end}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{tohyphen .Service}}{{if .Client}}-client{{end}}
  template:
    metadata:
      labels:
        app: {{tohyphen .Service}}{{if .Client}}-client{{end}}
    spec:
		{{- if .Tern}}
      initContainers:
      - name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-migrations
        securityContext:
          allowPrivilegeEscalation: false
        image: golang:alpine
        envFrom:
        - configMapRef:
            name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-env
        - secretRef:
            name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-postgres-env
        volumeMounts:
          - mountPath: /migrations
            name: migrations
        command:
          - sh
          - "-c"
          - |
            go install github.com/jackc/tern@latest
            tern migrate --migrations /migrations
		{{- end}}
      containers:
      - name: {{tohyphen .Service}}{{if .Client}}-client{{end}}
        securityContext:
          allowPrivilegeEscalation: false
        image: {{tohyphen .Service}}{{if .Client}}-client{{end}}:latest
        envFrom:
        - configMapRef:
            name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-env
				{{- if .Tern}}
        - secretRef:
            name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-postgres-env
				{{- end}}
			  {{- if .Health}}
        readinessProbe:
          grpc:
            port: 41888
          initialDelaySeconds: 10
          timeoutSeconds: 5
        livenessProbe:
          grpc:
            port: 41888
          initialDelaySeconds: 10
          timeoutSeconds: 5
				{{- end}}
			{{- if .Tern}}
      volumes:
      - name: migrations
        configMap:
          name: {{tohyphen .Service}}{{if .Client}}-client{{end}}-migrations
			{{- end}}
`
