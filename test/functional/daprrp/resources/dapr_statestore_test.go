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

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/radius-project/radius/test/functional"
	"github.com/radius-project/radius/test/functional/shared"
	"github.com/radius-project/radius/test/step"
	"github.com/radius-project/radius/test/validation"
)

func Test_DaprStateStore_Manual(t *testing.T) {
	template := "testdata/daprrp-resources-statestore-manual.bicep"
	name := "daprrp-rs-statestore-manual"
	appNamespace := "default-daprrp-rs-statestore-manual"

	test := shared.NewRPTest(t, name, []shared.TestStep{
		{
			Executor: step.NewDeployExecutor(template, functional.GetMagpieImage(), fmt.Sprintf("namespace=%s", appNamespace)),
			RPResources: &validation.RPResourceSet{
				Resources: []validation.RPResource{
					{
						Name: "daprrp-rs-statestore-manual",
						Type: validation.ApplicationsResource,
					},
					{
						Name: "dapr-sts-manual-ctnr",
						Type: validation.ContainersResource,
						App:  name,
					},
					{
						Name: "dapr-sts-manual",
						Type: validation.DaprStateStoresResource,
						App:  name,
					},
				},
			},
			K8sObjects: &validation.K8sObjectSet{
				Namespaces: map[string][]validation.K8sObject{
					appNamespace: {
						validation.NewK8sPodForResource(name, "dapr-sts-manual-ctnr"),

						// Deployed as supporting resources using Kubernetes Bicep extensibility.
						validation.NewK8sPodForResource(name, "dapr-sts-manual-redis").
							ValidateLabels(false),
						validation.NewK8sServiceForResource(name, "dapr-sts-manual-redis").
							ValidateLabels(false),

						validation.NewDaprComponent(name, "dapr-sts-manual").
							ValidateLabels(false),
					},
				},
			},
		},
	})

	test.RequiredFeatures = []shared.RequiredFeature{shared.FeatureDapr}

	test.PostDeleteVerify = func(ctx context.Context, t *testing.T, test shared.RPTest) {
		verifyDaprComponentsDeleted(ctx, t, test, "Applications.Dapr/stateStores", "dapr-sts-manual", appNamespace)
	}

	test.Test(t)
}

func Test_DaprStateStore_Recipe(t *testing.T) {
	template := "testdata/daprrp-resources-statestore-recipe.bicep"
	name := "daprrp-rs-sts-recipe"
	appNamespace := "daprrp-env-recipes-env"

	test := shared.NewRPTest(t, name, []shared.TestStep{
		{
			Executor: step.NewDeployExecutor(template, functional.GetMagpieImage(), functional.GetBicepRecipeRegistry(), functional.GetBicepRecipeVersion()),
			RPResources: &validation.RPResourceSet{
				Resources: []validation.RPResource{
					{
						Name: "daprrp-env-recipes-env",
						Type: validation.EnvironmentsResource,
					},
					{
						Name: "daprrp-rs-sts-recipe",
						Type: validation.ApplicationsResource,
						App:  name,
					},
					{
						Name: "dapr-sts-recipe-ctnr",
						Type: validation.ContainersResource,
						App:  name,
					},
					{
						Name: "dapr-sts-recipe",
						Type: validation.DaprStateStoresResource,
						App:  name,
					},
				},
			},
			K8sObjects: &validation.K8sObjectSet{
				Namespaces: map[string][]validation.K8sObject{
					appNamespace: {
						validation.NewK8sPodForResource(name, "dapr-sts-recipe-ctnr").
							ValidateLabels(false),

						validation.NewDaprComponent(name, "dapr-sts-recipe").
							ValidateLabels(false),
					},
				},
			},
		},
	})

	test.RequiredFeatures = []shared.RequiredFeature{shared.FeatureDapr}

	test.PostDeleteVerify = func(ctx context.Context, t *testing.T, test shared.RPTest) {
		verifyDaprComponentsDeleted(ctx, t, test, "Applications.Dapr/stateStores", "dapr-sts-recipe", appNamespace)
	}

	test.Test(t)
}
