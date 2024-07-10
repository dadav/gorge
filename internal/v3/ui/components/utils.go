package components

import (
	"encoding/json"
	"strings"

	model "github.com/dadav/gorge/internal/model"
)

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

func normalize(name string) string {
	return strings.Replace(name, "/", "-", 1)
}
