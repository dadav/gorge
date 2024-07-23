package ui

import (
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/dadav/gorge/internal/log"
	customMiddleware "github.com/dadav/gorge/internal/middleware"
	"github.com/dadav/gorge/internal/v3/backend"
	"github.com/dadav/gorge/internal/v3/ui/components"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"github.com/go-chi/chi/v5"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		w.WriteHeader(500)
		log.Log.Error(err)
		return
	}
	templ.Handler(components.Page("Gorge", components.SearchView("", modules, 1, 1, 10))).ServeHTTP(w, r)
}

type Pagination struct {
	Current int
	Max     int
}

func NewPagination(current, max int) Pagination {
	return Pagination{
		Current: current,
		Max:     max,
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")

	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		w.WriteHeader(500)
		log.Log.Error(err)
		return
	}

	filtered := []*gen.Module{}

	for _, module := range modules {
		ok := true
		for _, q := range strings.Split(query, " ") {
			if !strings.Contains(module.Name, q) && !strings.Contains(module.Owner.Username, q) && !strings.Contains(module.CurrentRelease.Version, q) {
				ok = false
			}
		}
		if ok {
			filtered = append(filtered, module)
		}
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Slug < filtered[j].Slug
	})

	pageInt := 1
	pageSizeInt := 10

	if page != "" {
		pageInt, _ = strconv.Atoi(page)
		if pageInt < 1 {
			pageInt = 1
		}
	}
	if pageSize != "" {
		pageSizeInt, _ = strconv.Atoi(pageSize)
	}

	maxPageInt := int(math.Ceil(float64(len(filtered)) / float64(pageSizeInt)))

	startIdx := (pageInt - 1) * pageSizeInt
	endIdx := startIdx + pageSizeInt

	if startIdx > len(filtered) {
		startIdx = len(filtered)
		w.WriteHeader(404)
		return
	}

	if endIdx > len(filtered) {
		endIdx = len(filtered)
	}

	filtered = filtered[startIdx:endIdx]

	templ.Handler(components.Page("Gorge", components.SearchView(query, filtered, pageInt, maxPageInt, pageSizeInt))).ServeHTTP(w, r)
}

func AuthorHandler(w http.ResponseWriter, r *http.Request) {
	authorSlug := chi.URLParam(r, "author")
	modules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		w.WriteHeader(500)
		log.Log.Error(err)
		return
	}

	filtered := []*gen.Module{}

	for _, module := range modules {
		if module.Owner.Slug == authorSlug {
			filtered = append(filtered, module)
		}
	}

	if len(filtered) > 0 {
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
		w.WriteHeader(500)
		log.Log.Error(err)
		return
	}

	for _, release := range releases {
		if release.Module.Slug == moduleSlug && release.Version == version {
			if release.Version == version {
				templ.Handler(components.Page(release.Slug, components.ReleaseView(release))).ServeHTTP(w, r)
				return
			}
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
