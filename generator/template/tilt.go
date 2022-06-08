package template

var Tiltfile = `{{if .PrivateRepo -}}
# Start SSH Agent
if os.getenv('SSH_AUTH_SOCK', '') == "":
    git_key = os.getenv('GIT_SSH_KEY', '')
    local('eval $(ssh-agent) && ssh-add {}'.format(git_key))

{{end -}}
# Build Docker image
docker_build('{{.Service}}{{if .Client}}-client{{end}}',
             context='.',
             dockerfile='./Dockerfile',
{{- if .PrivateRepo}}
             ssh='default',
{{- end}}
)

# Config
{{- if .Kustomize}}KUSTOMIZE_DIR="./resources/dev/"

{{- if .Tern}}# LoadRestrictor option is passed as migrations need to be accessed outside of base directory
# See: https://github.com/kubernetes-sigs/kustomize/issues/865
manifests = local("kustomize build --load-restrictor LoadRestrictionsNone {dir}".format(dir=KUSTOMIZE_DIR), quiet=True)
{{- else}}
manifests = kustomize(KUSTOMIZE_DIR)
{{- end}}
{{- else}}
KUBERNETS_DIR="./resources"
manifests = listdir(KUBERNETES_DIR)
{{- end}}

# Apply Kubernetes manifests
# Allow duplcates is marked true for when you import multiple go-micro Tiltfiles
#  into a single Tiltfile it will mark the clusterrole as duplicate.
k8s_yaml(manifests, allow_duplicates=True)
`
