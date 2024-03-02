package v3

import (
	"context"
	"errors"
	"net/http"

	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

type SearchFilterOperationsApi struct {
	gen.SearchFilterOperationsAPIServicer
}

func NewSearchFilterOperationsApi() *SearchFilterOperationsApi {
	return &SearchFilterOperationsApi{}
}

// AddSearchFilter - Create search filter
func (s *SearchFilterOperationsApi) AddSearchFilter(ctx context.Context, searchFilterSlug string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	// TODO - update AddSearchFilter with the required logic for this service method.
	// Add api_search_filter_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, SearchFilter{}) or use other options such as http.Ok ...
	// return Response(200, SearchFilter{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return Response(304, nil),nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	// TODO: Uncomment the next line to return response Response(409, AddSearchFilter409Response{}) or use other options such as http.Ok ...
	// return Response(409, AddSearchFilter409Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("AddSearchFilter method not implemented")
}

// DeleteUserSearchFilter - Delete search filter by ID
func (s *SearchFilterOperationsApi) DeleteUserSearchFilter(ctx context.Context, id int32) (gen.ImplResponse, error) {
	// TODO - update DeleteUserSearchFilter with the required logic for this service method.
	// Add api_search_filter_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	// TODO: Uncomment the next line to return response Response(403, DeleteUserSearchFilter403Response{}) or use other options such as http.Ok ...
	// return Response(403, DeleteUserSearchFilter403Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("DeleteUserSearchFilter method not implemented")
}

// GetUserSearchFilters - Get user&#39;s search filters
func (s *SearchFilterOperationsApi) GetUserSearchFilters(ctx context.Context) (gen.ImplResponse, error) {
	// TODO - update GetUserSearchFilters with the required logic for this service method.
	// Add api_search_filter_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, SearchFilterResponse{}) or use other options such as http.Ok ...
	// return Response(200, SearchFilterResponse{}), nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, GetUserSearchFilters401Response{}) or use other options such as http.Ok ...
	// return Response(401, GetUserSearchFilters401Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetUserSearchFilters method not implemented")
}
