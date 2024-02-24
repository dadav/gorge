package backend

import gen "github.com/dadav/gorge/pkg/gen/v3/openapi"

type Backend interface {
	// LoadModules loads modules into memory
	LoadModules() error

	// GetAllModules returns a list of all modules
	GetAllModules() []*gen.Module

	// GetModuleBySlug contains a map to modules
	GetModuleBySlug(string) (*gen.Module, error)

	// GetAllReleases returns a list of all releases
	GetAllReleases() []*gen.Release

	// GetReleaseBySlug contains a map to modules
	GetReleaseBySlug(string) (*gen.Release, error)

	// Add a new release
	AddRelease(module string, version string, data []byte) error
}
