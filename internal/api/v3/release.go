package v3

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dadav/gorge/internal/backend"
	"github.com/dadav/gorge/internal/config"
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
	// TODO - update AddRelease with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(201, ReleaseMinimal{}) or use other options such as http.Ok ...
	// return Response(201, ReleaseMinimal{}), nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, GetUserSearchFilters401Response{}) or use other options such as http.Ok ...
	// return Response(401, GetUserSearchFilters401Response{}), nil

	// TODO: Uncomment the next line to return response Response(403, DeleteUserSearchFilter403Response{}) or use other options such as http.Ok ...
	// return Response(403, DeleteUserSearchFilter403Response{}), nil

	// TODO: Uncomment the next line to return response Response(409, AddSearchFilter409Response{}) or use other options such as http.Ok ...
	// return Response(409, AddSearchFilter409Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("AddRelease method not implemented")
}

// DeleteRelease - Delete module release
func (s *ReleaseOperationsApi) DeleteRelease(ctx context.Context, releaseSlug string, reason string) (gen.ImplResponse, error) {
	// TODO - update DeleteRelease with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, GetUserSearchFilters401Response{}) or use other options such as http.Ok ...
	// return Response(401, GetUserSearchFilters401Response{}), nil

	// TODO: Uncomment the next line to return response Response(403, DeleteUserSearchFilter403Response{}) or use other options such as http.Ok ...
	// return Response(403, DeleteUserSearchFilter403Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("DeleteRelease method not implemented")
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
	metadata, readme, err := backend.ReadReleaseMetadata(releaseSlug)
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

	return gen.Response(http.StatusOK, gen.Release{
		Slug:   releaseSlug,
		Module: gen.ReleaseModule{Name: metadata.Name},
		Readme: readme,
	}), nil
}

// GetReleasePlan - Fetch module release plan
func (s *ReleaseOperationsApi) GetReleasePlan(ctx context.Context, releaseSlug string, planName string) (gen.ImplResponse, error) {
	// TODO - update GetReleasePlan with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, ReleasePlan{}) or use other options such as http.Ok ...
	// return Response(200, ReleasePlan{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetReleasePlan method not implemented")
}

// GetReleasePlans - List module release plans
func (s *ReleaseOperationsApi) GetReleasePlans(ctx context.Context, releaseSlug string) (gen.ImplResponse, error) {
	// TODO - update GetReleasePlans with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, GetReleasePlans200Response{}) or use other options such as http.Ok ...
	// return Response(200, GetReleasePlans200Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetReleasePlans method not implemented")
}

// GetReleases - List module releases
func (s *ReleaseOperationsApi) GetReleases(ctx context.Context, limit int32, offset int32, sortBy string, module string, owner string, withPdk bool, operatingsystem string, operatingsystemrelease string, peRequirement string, puppetRequirement string, moduleGroups []string, showDeleted bool, hideDeprecated bool, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string, supported bool) (gen.ImplResponse, error) {
	results := []gen.Release{}
	filtered := []gen.Release{}
	allReleases := backend.ConfiguredBackend.GetAllReleases()
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
