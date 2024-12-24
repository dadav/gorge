package v3

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dadav/gorge/internal/v3/backend"
	"github.com/dadav/gorge/internal/v3/utils"
	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

const (
	defaultLimit  = 20
	maxLimit      = 100
	defaultOffset = 0
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
	// Validate module slug
	if !utils.CheckModuleSlug(moduleSlug) {
		err := errors.New("invalid module slug")
		return gen.Response(
			http.StatusBadRequest,
			DeleteModule500Response{
				Message: err.Error(),
				Errors:  []string{err.Error()},
			},
		), nil
	}

	// Check if module exists
	module, err := backend.ConfiguredBackend.GetModuleBySlug(moduleSlug)
	if err != nil {
		return gen.Response(
			http.StatusNotFound,
			GetModule404Response{
				Message: "Module not found",
				Errors:  []string{"Module could not be found"},
			},
		), nil
	}

	// Update module deprecation status
	deprecatedAt := time.Now().UTC().Format(time.RFC3339)
	module.DeprecatedAt = &deprecatedAt
	module.DeprecatedFor = deprecationRequest.Params.Reason

	if *deprecationRequest.Params.ReplacementSlug != "" {
		module.SupersededBy = gen.ModuleSupersededBy{
			Slug: *deprecationRequest.Params.ReplacementSlug,
		}
	}

	// Save the updated module
	err = backend.ConfiguredBackend.UpdateModule(module)
	if err != nil {
		return gen.Response(
			http.StatusInternalServerError,
			DeleteModule500Response{
				Message: "Failed to deprecate module",
				Errors:  []string{err.Error()},
			},
		), nil
	}

	return gen.Response(http.StatusNoContent, nil), nil
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
	// Validate and set defaults for limit/offset
	if limit <= 0 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = defaultOffset
	}

	// Get all modules with error handling
	allModules, err := backend.ConfiguredBackend.GetAllModules()
	if err != nil {
		return gen.Response(
			http.StatusInternalServerError,
			GetModule500Response{
				Message: "Failed to fetch modules",
				Errors:  []string{err.Error()},
			}), nil
	}

	// Check offset validity early
	if int(offset) >= len(allModules) {
		return gen.Response(
			http.StatusNotFound,
			GetModule404Response{
				Message: "Invalid offset",
				Errors:  []string{"The given offset is larger than the total number of modules"},
			}), nil
	}

	// Create filter function type and map of filters
	type filterFunc func(m *gen.Module) bool
	filters := make(map[string]filterFunc)

	// Add filters conditionally
	if query != "" {
		filters["query"] = func(m *gen.Module) bool {
			return strings.Contains(m.Slug, query) || strings.Contains(m.Owner.Slug, query)
		}
	}
	if tag != "" {
		filters["tag"] = func(m *gen.Module) bool {
			return slices.Contains(m.CurrentRelease.Tags, tag)
		}
	}
	if owner != "" {
		filters["owner"] = func(m *gen.Module) bool {
			return m.Owner.Username == owner
		}
	}
	if withTasks {
		filters["with_tasks"] = func(m *gen.Module) bool {
			return len(m.CurrentRelease.Tasks) > 0
		}
	}
	if withPlans {
		filters["with_plans"] = func(m *gen.Module) bool {
			return len(m.CurrentRelease.Plans) > 0
		}
	}
	if withPdk {
		filters["with_pdk"] = func(m *gen.Module) bool {
			return m.CurrentRelease.Pdk
		}
	}
	if premium {
		filters["premium"] = func(m *gen.Module) bool {
			return m.Premium
		}
	}
	if excludePremium {
		filters["exclude_premium"] = func(m *gen.Module) bool {
			return !m.Premium
		}
	}
	if len(endorsements) > 0 {
		filters["endorsements"] = func(m *gen.Module) bool {
			return m.Endorsement != nil && slices.Contains(endorsements, *m.Endorsement)
		}
	}

	// Apply filters
	filtered := make([]gen.Module, 0)
	for _, m := range allModules[offset:] {
		passes := true
		for _, filter := range filters {
			if !filter(m) {
				passes = false
				break
			}
		}
		if passes {
			filtered = append(filtered, *m)
		}
	}

	// Apply pagination
	endIndex := int(limit)
	if endIndex > len(filtered) {
		endIndex = len(filtered)
	}
	results := filtered[:endIndex]

	base, _ := url.Parse("/v3/modules")
	params := url.Values{}
	params.Add("offset", strconv.Itoa(int(offset)))
	params.Add("limit", strconv.Itoa(int(limit)))

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
