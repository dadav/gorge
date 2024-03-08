package v3

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/dadav/gorge/internal/log"
	"github.com/dadav/gorge/internal/v3/backend"
	"github.com/dadav/gorge/internal/v3/utils"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

type ModuleOperationsApi struct {
	gen.ModuleOperationsAPIServicer
}

func NewModuleOperationsApi() *ModuleOperationsApi {
	return &ModuleOperationsApi{}
}

type DeleteModule500Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// DeleteModule - Delete module
func (s *ModuleOperationsApi) DeleteModule(ctx context.Context, moduleSlug string, reason string) (gen.ImplResponse, error) {
	if !utils.CheckModuleSlug(moduleSlug) {
		err := errors.New("invalid module slug")
		return gen.Response(
			400,
			DeleteModule500Response{
				Message: err.Error(),
				Errors:  []string{err.Error()},
			},
		), nil
	}

	err := backend.ConfiguredBackend.DeleteModuleBySlug(moduleSlug)
	if err == nil {
		return gen.Response(204, nil), nil
	}

	return gen.Response(
		500,
		DeleteModule500Response{
			Message: err.Error(),
			Errors:  []string{err.Error()},
		},
	), nil
}

// DeprecateModule - Deprecate module
func (s *ModuleOperationsApi) DeprecateModule(ctx context.Context, moduleSlug string, deprecationRequest gen.DeprecationRequest) (gen.ImplResponse, error) {
	// TODO - update DeprecateModule with the required logic for this service method.
	// Add api_module_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

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

	return gen.Response(http.StatusNotImplemented, nil), errors.New("DeprecateModule method not implemented")
}

type GetModule404Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type GetModule500Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// GetModule - Fetch module
func (s *ModuleOperationsApi) GetModule(ctx context.Context, moduleSlug string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	module, err := backend.ConfiguredBackend.GetModuleBySlug(moduleSlug)
	if err != nil {
		log.Log.Error(err)
		return gen.Response(
			http.StatusNotFound,
			GetModule404Response{
				Message: http.StatusText(http.StatusNotFound),
				Errors:  []string{"Module could not be found"},
			}), nil
	}

	return gen.Response(http.StatusOK, module), nil
}

// GetModules - List modules
func (s *ModuleOperationsApi) GetModules(ctx context.Context, limit int32, offset int32, sortBy string, query string, tag string, owner string, withTasks bool, withPlans bool, withPdk bool, premium bool, excludePremium bool, endorsements []string, operatingsystem string, operatingsystemrelease string, peRequirement string, puppetRequirement string, withMinimumScore int32, moduleGroups []string, showDeleted bool, hideDeprecated bool, onlyLatest bool, slugs []string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string, startsWith string, supported bool, withReleaseSince string) (gen.ImplResponse, error) {
	results := []gen.Module{}
	filtered := []gen.Module{}
	allModules, _ := backend.ConfiguredBackend.GetAllModules()
	params := url.Values{}
	params.Add("offset", strconv.Itoa(int(offset)))
	params.Add("limit", strconv.Itoa(int(limit)))

	if int(offset)+1 > len(allModules) {
		return gen.Response(404, GetModule404Response{
			Message: "Invalid offset",
			Errors:  []string{"The given offset is larger than the total number of (filtered) modules."},
		}), nil
	}

	for _, m := range allModules[offset:] {
		var filterMatched, filterSet bool
		if query != "" {
			filterSet = true
			filterMatched = strings.Contains(m.Slug, query) || strings.Contains(m.Owner.Slug, query)
			params.Add("query", query)
		}

		if tag != "" {
			filterSet = true
			filterMatched = slices.Contains(m.CurrentRelease.Tags, tag)
			params.Add("tag", tag)
		}

		if owner != "" {
			filterSet = true
			filterMatched = m.Owner.Username == owner
			params.Add("owner", owner)
		}

		if withTasks {
			filterSet = true
			filterMatched = len(m.CurrentRelease.Tasks) > 0
			params.Add("with_tasks", strconv.FormatBool(withTasks))
		}

		if withPlans {
			filterSet = true
			filterMatched = len(m.CurrentRelease.Plans) > 0
			params.Add("with_plans", strconv.FormatBool(withPlans))
		}

		if withPdk {
			filterSet = true
			filterMatched = m.CurrentRelease.Pdk
			params.Add("with_pdk", strconv.FormatBool(withPdk))
		}

		if premium {
			filterSet = true
			filterMatched = m.Premium
			params.Add("premium", strconv.FormatBool(premium))
		}

		if excludePremium {
			filterSet = true
			filterMatched = !m.Premium
			params.Add("exclude_premium", strconv.FormatBool(excludePremium))
		}

		if len(endorsements) > 0 {
			filterSet = true
			filterMatched = m.Endorsement != nil && slices.Contains(endorsements, *m.Endorsement)
			params.Add("endorsements", "["+strings.Join(endorsements, ",")+"]")
		}

		if !filterSet || filterMatched {
			filtered = append(filtered, *m)
		}
	}

	i := 1
	for _, module := range filtered {
		if i > int(limit) {
			break
		}
		results = append(results, module)
		i++
	}

	base, _ := url.Parse("/v3/modules")
	base.RawQuery = params.Encode()
	currentInf := interface{}(base.String())
	params.Set("offset", "0")
	firstInf := interface{}(base.String())
	params.Set("offset", strconv.Itoa(int(offset)+len(results)))
	nextInf := interface{}(base.String())

	return gen.Response(http.StatusOK, gen.GetModules200Response{
		Pagination: gen.GetModules200ResponsePagination{
			Limit:   limit,
			Offset:  offset,
			First:   &firstInf,
			Current: &currentInf,
			Next:    &nextInf,
			Total:   int32(len(allModules)),
		},
		Results: results,
	}), nil
}
