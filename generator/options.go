package generator

// Options represents the options for the generator.
type Options struct {
	// Service is the name of the service the generator will generate files
	// for.
	Service string
	// Vendor is the service vendor.
	Vendor string
	// Directory is the directory where the files will be generated to.
	Directory string

	// Client determines whether or not the project is a client project.
	Client bool
	// Jaeger determines whether or not Jaeger integration is enabled.
	Jaeger bool
	// Skaffold determines whether or not Skaffold integration is enabled.
	Skaffold bool
	// Tilt determines whether or not Tilt integration is enabled.
	Tilt bool
	// Health determines whether or not health proto service is enabled.
	Health bool
	// Kustomize determines whether or not Kustomize integration is enabled.
	Kustomize bool
	// Sqlc determines whether or not Sqlc integration is enabled.
	Sqlc bool
	// GRPC determines whether or not GRPC integration is enabled.
	GRPC bool
	// Buildkit determines whether or not Buildkit integration is enabled.
	Buildkit bool
	// Tern directory whether or not Tern integration is enabled.
	Tern bool
	// Advanced directory whether or not Advanced integration is enabled.
	Advanced bool
	// PrivateRepo
	PrivateRepo bool
	// Namespace sets the default namespace
	Namespace string
	// PostgresAddress sets the default postgres address
	PostgresAddress string
}

// Option manipulates the Options passed.
type Option func(o *Options)

// Service sets the service name.
func Service(s string) Option {
	return func(o *Options) {
		o.Service = s
	}
}

// Vendor sets the service vendor.
func Vendor(v string) Option {
	return func(o *Options) {
		o.Vendor = v
	}
}

// Directory sets the directory in which files are generated.
func Directory(d string) Option {
	return func(o *Options) {
		o.Directory = d
	}
}

// Client sets whether or not the project is a client project.
func Client(c bool) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Jaeger sets whether or not Jaeger integration is enabled.
func Jaeger(j bool) Option {
	return func(o *Options) {
		o.Jaeger = j
	}
}

// Skaffold sets whether or not Skaffold integration is enabled.
func Skaffold(s bool) Option {
	return func(o *Options) {
		o.Skaffold = s
	}
}

// Tilt sets whether or not Tilt integration is enabled.
func Tilt(s bool) Option {
	return func(o *Options) {
		o.Tilt = s
	}
}

// Health determines whether or not health proto service is enabled.
func Health(s bool) Option {
	return func(o *Options) {
		o.Health = s
	}
}

// Kustomize determines whether or not Kustomize integration is enabled.
func Kustomize(s bool) Option {
	return func(o *Options) {
		o.Kustomize = s
	}
}

// Sqlc determines whether or not Sqlc integration is enabled.
func Sqlc(s bool) Option {
	return func(o *Options) {
		o.Sqlc = s
	}
}

// GRPC determines whether or not GRPC integration is enabled.
func GRPC(s bool) Option {
	return func(o *Options) {
		o.GRPC = s
	}
}

// Buildkit determines whether or not Buildkit integration is enabled.
func Buildkit(s bool) Option {
	return func(o *Options) {
		o.Buildkit = s
	}
}

// Tern determines whether or not Tern integration is enabled.
func Tern(s bool) Option {
	return func(o *Options) {
		o.Tern = s
	}
}

// Advanced determines whether or not Advanced integration is enabled.
func Advanced(s bool) Option {
	return func(o *Options) {
		o.Advanced = s
	}
}

// PrivateRepo
func PrivateRepo(s bool) Option {
	return func(o *Options) {
		o.PrivateRepo = s
	}
}

// Namespace
func Namespace(s string) Option {
	return func(o *Options) {
		o.Namespace = s
	}
}

// PostgresAddress
func PostgresAddress(s string) Option {
	return func(o *Options) {
		o.PostgresAddress = s
	}
}
