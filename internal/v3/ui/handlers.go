package ui

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/dadav/gorge/internal/log"
	customMiddleware "github.com/dadav/gorge/internal/middleware"
	"github.com/dadav/gorge/internal/v3/backend"
	"github.com/dadav/gorge/internal/v3/ui/components"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"github.com/go-chi/chi/v5"
)

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Log.Error(err)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		handleError(w, err)
		return
	}
	templ.Handler(components.Page("Gorge", components.SearchView("", modules))).ServeHTTP(w, r)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("query"))
	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		handleError(w, err)
		return
	}

	filtered := make([]*gen.Module, 0, len(modules))
	queryTerms := strings.Fields(query)

	for _, module := range modules {
		matches := true
		for _, term := range queryTerms {
			if !strings.Contains(strings.ToLower(module.Name), term) &&
				!strings.Contains(strings.ToLower(module.Owner.Username), term) &&
				!strings.Contains(strings.ToLower(module.CurrentRelease.Version), term) {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, module)
		}
	}

	templ.Handler(components.Page("Gorge", components.SearchView(query, filtered))).ServeHTTP(w, r)
}

func AuthorHandler(w http.ResponseWriter, r *http.Request) {
	authorSlug := chi.URLParam(r, "author")
	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		handleError(w, err)
		return
	}

	authorModules := make(map[string][]*gen.Module)
	for _, module := range modules {
		authorModules[module.Owner.Slug] = append(authorModules[module.Owner.Slug], module)
	}

	if filtered, exists := authorModules[authorSlug]; exists && len(filtered) > 0 {
		templ.Handler(components.Page(authorSlug, components.AuthorView(filtered))).ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}

func ReleaseHandler(w http.ResponseWriter, r *http.Request) {
	moduleSlug := chi.URLParam(r, "module")
	version := chi.URLParam(r, "version")
	releases, err := backend.ConfiguredBackend.GetAllReleases()
	if err != nil {
		handleError(w, err)
		return
	}

	for _, release := range releases {
		if release.Module.Slug == moduleSlug && release.Version == version {
			templ.Handler(components.Page(release.Slug, components.ReleaseView(release))).ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func ModuleHandler(w http.ResponseWriter, r *http.Request) {
	moduleSlug := chi.URLParam(r, "module")
	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		w.WriteHeader(500)
		log.Log.Error(err)
		return
	}

	for _, module := range modules {
		if module.Slug == moduleSlug {
			templ.Handler(components.Page(module.Slug, components.ModuleView(module))).ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func StatisticsHandler(stats *customMiddleware.Statistics) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		stats.Mutex.Lock()
		defer stats.Mutex.Unlock()
		templ.Handler(components.Page("Statistics", components.StatisticsView(stats))).ServeHTTP(w, r)
	}
}
