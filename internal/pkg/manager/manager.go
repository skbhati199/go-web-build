package manager

import "context"

type PackageManager interface {
	Install(ctx context.Context, packages ...string) error
	Uninstall(ctx context.Context, packages ...string) error
	Update(ctx context.Context) error
	List(ctx context.Context) ([]Package, error)
}

type Package struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	DevMode      bool              `json:"devMode"`
}
