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
	"github.com/dadav/gorge/internal/v3/utils"
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

// Add custom error types for better error handling
type ReleaseError struct {
	Code    int
	Message string
	Errors  []string
}

// AddRelease - Create module release
func (s *ReleaseOperationsApi) AddRelease(ctx context.Context, addReleaseRequest gen.AddReleaseRequest) (gen.ImplResponse, error) {
	if addReleaseRequest.File == "" {
		return gen.Response(400, gen.GetFile400Response{
			Message: "No file data provided",
			Errors:  []string{"file data is required"},
		}), nil
	}

	decodedTarball, err := base64.StdEncoding.DecodeString(addReleaseRequest.File)
	if err != nil {
		return gen.Response(400, gen.GetFile400Response{
			Message: "Invalid base64 encoded data",
			Errors:  []string{err.Error()},
		}), nil
	}

	release, err := backend.ConfiguredBackend.AddRelease(decodedTarball)
	if err != nil {
		return gen.Response(400, gen.GetFile400Response{
			Message: "Failed to add release",
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
	if !utils.CheckReleaseSlug(releaseSlug) {
		err := errors.New("invalid release slug")
		return gen.Response(
			400,
			DeleteRelease500Response{
				Message: err.Error(),
				Errors:  []string{err.Error()},
			},
		), nil
	}
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

type GetFile400Response struct {
	Message string `json:"message,omitempty"`

	Errors []string `json:"errors,omitempty"`
}

// GetFile - Download module release
func (s *ReleaseOperationsApi) GetFile(ctx context.Context, filename string) (gen.ImplResponse, error) {

	if filename == "" {
		return gen.Response(400, gen.GetFile400Response{
			Message: "No filename provided",
			Errors:  []string{"filename is required"},
		}), nil
	}

	// Validate the filename to ensure it does not contain any path separators or parent directory references
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "..") {
		return gen.Response(400, gen.GetFile400Response{
			Message: "Invalid filename",
			Errors:  []string{"filename contains invalid characters"},
		}), nil
	}

	releaseSlug := strings.TrimSuffix(filename, ".tar.gz")
	if !utils.CheckReleaseSlug(releaseSlug) {
		return gen.Response(400, gen.GetFile400Response{
			Message: "Invalid release slug format",
			Errors:  []string{"release slug is invalid"},
		}), nil
	}

	filePath := filepath.Join(config.ModulesDir, ReleaseToModule(filename), filename)

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return gen.Response(http.StatusNotFound, gen.GetFile404Response{
				Message: "File not found",
				Errors:  []string{"the file does not exist"},
			}), nil
		}
		return gen.Response(http.StatusInternalServerError, GetRelease500Response{
			Message: "Failed to open file",
			Errors:  []string{err.Error()},
		}), nil
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
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

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
