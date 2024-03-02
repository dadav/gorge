package v3

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dadav/gorge/internal/config"
	"github.com/dadav/gorge/internal/v3/backend"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

type ReleaseOperationsApi struct {
	gen.ReleaseOperationsAPIServicer
}

func NewReleaseOperationsApi() *ReleaseOperationsApi {
	return &ReleaseOperationsApi{}
}

type GetRelease404Response struct {
	Message string `json:"message,omitempty"`

	Errors []string `json:"errors,omitempty"`
}

// AddRelease - Create module release
func (s *ReleaseOperationsApi) AddRelease(ctx context.Context, addReleaseRequest gen.AddReleaseRequest) (gen.ImplResponse, error) {
	base64EncodedTarball := addReleaseRequest.File

	decodedTarball, err := base64.StdEncoding.DecodeString(base64EncodedTarball)
	if err != nil {
		return gen.Response(400, gen.GetFile400Response{
			Message: "Could not decode provided data",
			Errors:  []string{err.Error()},
		}), nil
	}

	release, err := backend.ConfiguredBackend.AddRelease(decodedTarball)
	if err != nil {
		return gen.Response(400, gen.GetFile400Response{
			Message: "could not add release",
			Errors:  []string{err.Error()},
		}), nil
	}

	return gen.Response(201, gen.ReleaseMinimal{
		Uri:     release.Uri,
		FileUri: release.FileUri,
		Slug:    release.Slug,
	}), nil
}

type DeleteRelease500Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// DeleteRelease - Delete module release
func (s *ReleaseOperationsApi) DeleteRelease(ctx context.Context, releaseSlug string, reason string) (gen.ImplResponse, error) {
	err := backend.ConfiguredBackend.DeleteReleaseBySlug(releaseSlug)
	if err == nil {
		return gen.Response(204, nil), nil
	}

	return gen.Response(
		500,
		DeleteRelease500Response{
			Message: err.Error(),
			Errors:  []string{err.Error()},
		},
	), nil
}

func ReleaseToModule(releaseSlug string) string {
	return releaseSlug[:strings.LastIndex(releaseSlug, "-")]
}

// GetFile - Download module release
func (s *ReleaseOperationsApi) GetFile(ctx context.Context, filename string) (gen.ImplResponse, error) {
	f, err := os.Open(filepath.Join(config.ModulesDir, ReleaseToModule(filename), filename))
	if err != nil {
		if os.IsNotExist(err) {
			return gen.Response(http.StatusNotFound, gen.GetFile404Response{
				Message: http.StatusText(http.StatusNotFound),
				Errors:  []string{"The file does not exist."},
			}), nil
		}
	}

	return gen.Response(http.StatusOK, f), nil
}

type GetRelease500Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// GetRelease - Fetch module release
func (s *ReleaseOperationsApi) GetRelease(ctx context.Context, releaseSlug string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	release, err := backend.ConfiguredBackend.GetReleaseBySlug(releaseSlug)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return gen.Response(http.StatusNotFound, gen.GetFile404Response{
				Message: http.StatusText(http.StatusNotFound),
				Errors:  []string{"release not found"},
			}), nil
		}
		return gen.Response(http.StatusInternalServerError, GetRelease500Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Errors:  []string{"error while reading release metadata"},
		}), nil
	}

	return gen.Response(http.StatusOK, release), nil
}

func abbrReleaseToFullReleasePlan(abbrReleasePlan gen.ReleasePlanAbbreviated) gen.ReleasePlan {
	planFile := fmt.Sprintf("plans/%s.pp", strings.Join(strings.Split(abbrReleasePlan.Name, "::")[1:], "/"))
	return gen.ReleasePlan{
		Uri:      abbrReleasePlan.Uri,
		Name:     abbrReleasePlan.Name,
		Private:  abbrReleasePlan.Private,
		Filename: planFile,
		PlanMetadata: gen.ReleasePlanPlanMetadata{
			Name:    abbrReleasePlan.Name,
			Private: abbrReleasePlan.Private,
			File:    planFile,
		},
	}
}

// GetReleasePlan - Fetch module release plan
func (s *ReleaseOperationsApi) GetReleasePlan(ctx context.Context, releaseSlug string, planName string) (gen.ImplResponse, error) {
	release, err := backend.ConfiguredBackend.GetReleaseBySlug(releaseSlug)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return gen.Response(http.StatusNotFound, gen.GetFile404Response{
				Message: http.StatusText(http.StatusNotFound),
				Errors:  []string{"plan not found"},
			}), nil
		}
		return gen.Response(http.StatusInternalServerError, GetRelease500Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Errors:  []string{"error while reading release metadata"},
		}), nil
	}

	for _, plan := range release.Plans {
		if plan.Name == planName {
			// modulename::foo becomes plans/foo.pp
			return gen.Response(200, abbrReleaseToFullReleasePlan(plan)), nil
		}
	}
	return gen.Response(http.StatusNotFound, gen.GetFile404Response{
		Message: http.StatusText(http.StatusNotFound),
		Errors:  []string{"plan not found"},
	}), nil
}

// GetReleasePlans - List module release plans
func (s *ReleaseOperationsApi) GetReleasePlans(ctx context.Context, releaseSlug string) (gen.ImplResponse, error) {
	release, err := backend.ConfiguredBackend.GetReleaseBySlug(releaseSlug)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return gen.Response(http.StatusNotFound, gen.GetFile404Response{
				Message: http.StatusText(http.StatusNotFound),
				Errors:  []string{"plan not found"},
			}), nil
		}
		return gen.Response(http.StatusInternalServerError, GetRelease500Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Errors:  []string{"error while reading release metadata"},
		}), nil
	}

	results := []gen.ReleasePlan{}
	for _, plan := range release.Plans {
		results = append(results, abbrReleaseToFullReleasePlan(plan))
	}

	return gen.Response(200, gen.GetReleasePlans200Response{
		Pagination: gen.GetReleasePlans200ResponsePagination{},
		Results:    results,
	}), nil
}

// GetReleases - List module releases
func (s *ReleaseOperationsApi) GetReleases(ctx context.Context, limit int32, offset int32, sortBy string, module string, owner string, withPdk bool, operatingsystem string, operatingsystemrelease string, peRequirement string, puppetRequirement string, moduleGroups []string, showDeleted bool, hideDeprecated bool, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string, supported bool) (gen.ImplResponse, error) {
	results := []gen.Release{}
	filtered := []gen.Release{}
	allReleases, _ := backend.ConfiguredBackend.GetAllReleases()
	params := url.Values{}
	params.Add("offset", strconv.Itoa(int(offset)))
	params.Add("limit", strconv.Itoa(int(limit)))

	if int(offset)+1 > len(allReleases) {
		return gen.Response(404, GetRelease404Response{
			Message: "Invalid offset",
			Errors:  []string{"The given offset is larger than the total number of modules."},
		}), nil
	}

	for _, r := range allReleases[offset:] {
		var filterMatched, filterSet bool

		if module != "" && r.Module.Slug != module {
			filterSet = true
			filterMatched = r.Module.Slug == module
			params.Add("module", module)
		}
		if owner != "" && r.Module.Owner.Slug != owner {
			filterSet = true
			filterMatched = r.Module.Owner.Slug == owner
			params.Add("owner", owner)
		}

		if !filterSet || filterMatched {
			filtered = append(filtered, *r)
		}
	}

	i := 1
	for _, release := range filtered {
		if i > int(limit) {
			break
		}
		results = append(results, release)
		i++
	}

	base, _ := url.Parse("/v3/releases")
	base.RawQuery = params.Encode()
	currentInf := interface{}(base.String())
	params.Set("offset", "0")
	firstInf := interface{}(base.String())

	var nextInf interface{}
	nextOffset := int(offset) + len(results)
	if nextOffset < len(filtered) {
		params.Set("offset", strconv.Itoa(nextOffset))
		nextInf = interface{}(base.String())
	}

	var prevInf *string
	prevOffset := int(offset) - int(limit)
	if prevOffset >= 0 {
		prevStr := base.String()
		prevInf = &prevStr
	}

	return gen.Response(http.StatusOK, gen.GetReleases200Response{
		Pagination: gen.GetReleases200ResponsePagination{
			Limit:    limit,
			Offset:   offset,
			First:    &firstInf,
			Previous: prevInf,
			Current:  &currentInf,
			Next:     &nextInf,
			Total:    int32(len(filtered)),
		},
		Results: results,
	}), nil
}
