/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package applications

import (
	"context"
	"net/http"
	"net/url"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/corerp/datamodel"
	"github.com/radius-project/radius/pkg/corerp/datamodel/converter"
	"github.com/radius-project/radius/pkg/sdk"
	"github.com/radius-project/radius/pkg/ucp/resources"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	ctrl "github.com/radius-project/radius/pkg/armrpc/frontend/controller"
	"github.com/radius-project/radius/pkg/armrpc/rest"
)

var (
	schemePrefix = "http://"
)

var _ ctrl.Controller = (*GetApplicationGraph)(nil)

// GetApplicationGraph is the controller implementation to get application graph.
type GetApplicationGraph struct {
	ctrl.Operation[*datamodel.Application, datamodel.Application]
}

// NewGetRecipeMetadata creates a new controller for retrieving recipe metadata from an environment.
func NewGetApplicationGraph(opts ctrl.Options) (ctrl.Controller, error) {
	return &GetApplicationGraph{
		ctrl.NewOperation(opts,
			ctrl.ResourceOptions[datamodel.Application]{
				RequestConverter:  converter.ApplicationDataModelFromVersioned,
				ResponseConverter: converter.ApplicationDataModelToVersioned,
			},
		),
	}, nil
}

func (ctrl *GetApplicationGraph) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (rest.Response, error) {

	sCtx := v1.ARMRequestContextFromContext(ctx)

	// Request route for getGraph has name of the operation as suffix which should be removed to get the resource id.
	// route id format: /planes/radius/local/resourcegroups/default/providers/Applications.Core/applications/corerp-resources-application-app/getGraph"
	applicationID := sCtx.ResourceID.Truncate()
	applicationResource, _, err := ctrl.GetResource(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	if applicationResource == nil {
		return rest.NewNotFoundResponse(sCtx.ResourceID), nil
	}
	//Application MUST have an environment id
	environmentID, err := resources.Parse(applicationResource.Properties.Environment)
	if err != nil {
		return nil, err
	}

	if req.URL.Scheme != "" {
		schemePrefix = req.URL.Scheme + "://"
	}

	clientOptions, err := constructClientOptions(req.Host)
	if err != nil {
		return nil, err
	}

	// get all resources in application scope
	applicationResources, err := listAllResourcesByApplication(ctx, applicationID, clientOptions)
	if err != nil {
		return nil, err
	}

	// get all resources in environment scope
	environmentResources, err := listAllResourcesByEnvironment(ctx, environmentID, clientOptions)
	if err != nil {
		return nil, err
	}

	graph := compute(applicationID.Name(), applicationResources, environmentResources)
	if err != nil {
		response := rest.NewInternalServerErrorARMResponse(v1.ErrorResponse{
			Error: v1.ErrorDetails{
				Code:    v1.CodeInternal,
				Message: err.Error(),
			},
		})
		return response, nil
	} else {
		return rest.NewOKResponse(graph), nil
	}

}

// Construct client options from the request
func constructClientOptions(host string) (*policy.ClientOptions, error) {
	baseUrl := schemePrefix + host
	_, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}
	conn, err := sdk.NewDirectConnection(baseUrl)
	if err != nil {
		return nil, err
	}
	clientOptions := sdk.NewClientOptions(conn)

	return clientOptions, nil
}
