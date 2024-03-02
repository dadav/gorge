package backend

import gen "github.com/dadav/gorge/pkg/gen/v3/openapi"

type Backend interface {
	// LoadModules loads modules into memory
	LoadModules() error

	// GetAllModules returns a list of all modules
	GetAllModules() ([]*gen.Module, error)

	// GetModuleBySlug contains a map to modules
	GetModuleBySlug(slug string) (*gen.Module, error)

	// GetAllReleases returns a list of all releases
	GetAllReleases() ([]*gen.Release, error)

	// GetReleaseBySlug returns a release by slug
	GetReleaseBySlug(slug string) (*gen.Release, error)

	// AddRelease adds a new release
	AddRelease(data []byte) (*gen.Release, error)

	// DeleteModuleBySlug deletes a module
	DeleteModuleBySlug(slug string) error

	// DeleteReleaseBySlug deletes a release by slug
	DeleteReleaseBySlug(slug string) error
}
