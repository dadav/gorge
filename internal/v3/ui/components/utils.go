package components

import (
	"encoding/json"
	"sort"
	"strings"

	customMiddleware "github.com/dadav/gorge/internal/middleware"
	model "github.com/dadav/gorge/internal/model"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

// getSortedKeys extracts and sorts endpoint keys from Statistics
// Returns a sorted slice of endpoint strings
func getSortedKeys(stats *customMiddleware.Statistics) []string {
	keys := []string{}
	for k := range stats.ConnectionsPerEndpoint {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// sortModules sorts a slice of modules alphabetically by name
// Returns the sorted slice of modules
func sortModules(modules []*gen.Module) []*gen.Module {
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})
	return modules
}

// deps extracts module dependencies from metadata
// Returns a slice of ModuleDependency or nil if parsing fails
func deps(metadata map[string]interface{}) []model.ModuleDependency {
	var result model.ReleaseMetadata

	jsonStr, err := json.Marshal(metadata)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		return nil
	}

	return result.Dependencies
}

// normalize replaces forward slashes with hyphens in a string
// Returns the normalized string
func normalize(name string) string {
	return strings.Replace(name, "/", "-", 1)
}
