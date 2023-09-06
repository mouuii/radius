//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220901privatepreview

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ResourcesClient contains the methods for the Resources group.
// Don't use this type directly, use NewResourcesClient() instead.
type ResourcesClient struct {
	internal *arm.Client
}

// NewResourcesClient creates a new instance of ResourcesClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewResourcesClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ResourcesClient, error) {
	cl, err := arm.NewClient(moduleName+".ResourcesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ResourcesClient{
	internal: cl,
	}
	return client, nil
}

// NewListPager - List resources in a resource group
//
// Generated from API version 2022-09-01-privatepreview
//   - planeType - The plane type.
//   - planeName - The name of the plane
//   - resourceGroupName - The name of resource group
//   - options - ResourcesClientListOptions contains the optional parameters for the ResourcesClient.NewListPager method.
func (client *ResourcesClient) NewListPager(planeType string, planeName string, resourceGroupName string, options *ResourcesClientListOptions) (*runtime.Pager[ResourcesClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ResourcesClientListResponse]{
		More: func(page ResourcesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ResourcesClientListResponse) (ResourcesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, planeType, planeName, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ResourcesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ResourcesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ResourcesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ResourcesClient) listCreateRequest(ctx context.Context, planeType string, planeName string, resourceGroupName string, options *ResourcesClientListOptions) (*policy.Request, error) {
	urlPath := "/planes/{planeType}/{planeName}/resourcegroups/{resourceGroupName}/resources"
	if planeType == "" {
		return nil, errors.New("parameter planeType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{planeType}", url.PathEscape(planeType))
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", planeName)
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ResourcesClient) listHandleResponse(resp *http.Response) (ResourcesClientListResponse, error) {
	result := ResourcesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GenericResourceListResult); err != nil {
		return ResourcesClientListResponse{}, err
	}
	return result, nil
}

