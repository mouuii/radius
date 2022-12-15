// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package datamodel

import (
	"context"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	daprComponentCRD = "components.dapr.io"

	// DaprMissingError is an error message that can be used when Dapr is not installed on the cluster.
	DaprMissingError = "Dapr is not installed in your Kubernetes cluster. Please install Dapr by following the instructions at https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/."
)

// IsDaprInstalled will check for Dapr to be installed in the deployment environment and return
// and true if it is installed. Callers of this function can use DaprMissingError for a friendly error
// message to send back to users.
//
// This check is based on the Dapr Component CRD, and only supports Kubernetes.
func IsDaprInstalled(ctx context.Context, kubeClient client.Client) (bool, error) {
	crd := &apiextv1.CustomResourceDefinition{}
	err := kubeClient.Get(ctx, client.ObjectKey{Name: daprComponentCRD}, crd)
	if apierrors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
