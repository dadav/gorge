package v3

import (
	"context"
	"errors"
	"net/http"

	gen "github.com/dadav/gorge/pkg/gen/v3/openapi"
)

type UserOperationsApi struct {
	gen.UserOperationsAPIServicer
}

func NewUserOperationsApi() *UserOperationsApi {
	return &UserOperationsApi{}
}

// GetUser - Fetch user
func (s *UserOperationsApi) GetUser(ctx context.Context, userSlug string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	// TODO - update GetUser with the required logic for this service method.
	// Add api_user_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, User{}) or use other options such as http.Ok ...
	// return gen.Response(200, User{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return gen.Response(304, nil),nil

	// TODO: Uncomment the next line to return response Response(400, GetFile400Response{}) or use other options such as http.Ok ...
	// return gen.Response(400, GetFile400Response{}), nil

	// TODO: Uncomment the next line to return response Response(404, GetFile404Response{}) or use other options such as http.Ok ...
	// return gen.Response(404, GetFile404Response{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetUser method not implemented")
}

// GetUsers - List users
func (s *UserOperationsApi) GetUsers(ctx context.Context, limit int32, offset int32, sortBy string, withHtml bool, includeFields []string, excludeFields []string, ifModifiedSince string) (gen.ImplResponse, error) {
	// TODO - update GetUsers with the required logic for this service method.
	// Add api_user_operations_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, GetUsers200Response{}) or use other options such as http.Ok ...
	// return gen.Response(200, GetUsers200Response{}), nil

	// TODO: Uncomment the next line to return response Response(304, {}) or use other options such as http.Ok ...
	// return gen.Response(304, nil),nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetUsers method not implemented")
}
