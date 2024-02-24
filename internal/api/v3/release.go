package v3

import (
	"context"
	"errors"
	"net/http"

	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

type ReleaseOperationsApi struct {
	gen.ReleaseOperationsAPIServicer
}

func NewReleaseOperationsApi() *ReleaseOperationsApi {
	return &ReleaseOperationsApi{}
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

// GetFile - Download module release
func (s *ReleaseOperationsApi) GetFile(ctx context.Context, filename string) (gen.ImplResponse, error) {
	// TODO - update GetFile with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, *os.File{}) or use other options such as http.Ok ...
	// return Response(200, *os.File{}), nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetFile method not implemented")
}

// GetRelease - Fetch module release
func (s *ReleaseOperationsApi) GetRelease(ctx context.Context, releaseSlug string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	// TODO - update GetRelease with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, Release{}) or use other options such as http.Ok ...
	// return Response(200, Release{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return Response(304, nil),nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetRelease method not implemented")
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
	// TODO - update GetReleases with the required logic for this service method.
	// Add api_release_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, GetReleases200Response{}) or use other options such as http.Ok ...
	// return Response(200, GetReleases200Response{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return Response(304, nil),nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetReleases method not implemented")
}
