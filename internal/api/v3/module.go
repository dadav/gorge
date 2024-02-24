package v3

import (
	"context"
	"errors"
	"net/http"

	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

type ModuleOperationsApi struct {
	gen.ModuleOperationsAPIServicer
}

func NewModuleOperationsApi() *ModuleOperationsApi {
	return &ModuleOperationsApi{}
}

// DeleteModule - Delete module
func (s *ModuleOperationsApi) DeleteModule(ctx context.Context, moduleSlug string, reason string) (gen.ImplResponse, error) {
	// TODO - update DeleteModule with the required logic for this service method.
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

	return gen.Response(http.StatusNotImplemented, nil), errors.New("DeleteModule method not implemented")
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

// GetModule - Fetch module
func (s *ModuleOperationsApi) GetModule(ctx context.Context, moduleSlug string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	// TODO - update GetModule with the required logic for this service method.
	// Add api_module_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Module{}) or use other options such as http.Ok ...
	// return Response(200, Module{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return Response(304, nil),nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetModule method not implemented")
}

// GetModules - List modules
func (s *ModuleOperationsApi) GetModules(ctx context.Context, limit int32, offset int32, sortBy string, query string, tag string, owner string, withTasks bool, withPlans bool, withPdk bool, premium bool, excludePremium bool, endorsements []string, operatingsystem string, operatingsystemrelease string, peRequirement string, puppetRequirement string, withMinimumScore int32, moduleGroups []string, showDeleted bool, hideDeprecated bool, onlyLatest bool, slugs []string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string, startsWith string, supported bool, withReleaseSince string) (gen.ImplResponse, error) {
	// TODO - update GetModules with the required logic for this service method.
	// Add api_module_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, GetModules200Response{}) or use other options such as http.Ok ...
	// return Response(200, GetModules200Response{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return Response(304, nil),nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetModules method not implemented")
}
