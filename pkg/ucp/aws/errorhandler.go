// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------
package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	armrpc_v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	armrpc_rest "github.com/project-radius/radius/pkg/armrpc/rest"
)

func IsAWSResourceNotFoundError(err error) bool {
	target := &types.ResourceNotFoundException{}
	return errors.As(err, &target)
}

func HandleAWSError(err error) (armrpc_rest.Response, error) {
	operationErr := &smithy.OperationError{}
	if !errors.As(err, &operationErr) {
		return nil, err
	}

	httpErr := &smithyhttp.ResponseError{}
	if !errors.As(operationErr.Err, &httpErr) {
		return nil, err
	}

	var apiErr smithy.APIError
	if !errors.As(httpErr.Err, &apiErr) {
		return nil, err
	}

	e := armrpc_v1.ErrorResponse{
		Error: armrpc_v1.ErrorDetails{
			Code:    apiErr.ErrorCode(),
			Message: apiErr.ErrorMessage(),
		},
	}

	// We can't always trust apiErr.Fault :-/
	fault := apiErr.ErrorFault()
	if fault == smithy.FaultUnknown {
		switch apiErr.ErrorCode() {
		case "ValidationException":
			fault = smithy.FaultClient
		default:
			fault = smithy.FaultServer
		}
	}

	if fault == smithy.FaultClient {
		return armrpc_rest.NewBadRequestARMResponse(e), nil
	}

	return armrpc_rest.NewInternalServerErrorARMResponse(e), nil
}

// AWSMissingPropertyError is an error type to be returned when the call to UCP CreateWithPost
// is missing values for one of the expected primary identifier properties
type AWSMissingPropertyError struct {
	PropertyName string
}

func (e *AWSMissingPropertyError) Is(target error) bool {
	_, ok := target.(*AWSMissingPropertyError)
	return ok
}

func (e *AWSMissingPropertyError) Error() string {
	return fmt.Sprintf("mandatory property %s is missing", e.PropertyName)
}
