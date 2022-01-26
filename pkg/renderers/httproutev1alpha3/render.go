// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package httproutev1alpha3

import (
	"context"
	"errors"
	"fmt"

	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/project-radius/radius/pkg/azure/azresources"
	"github.com/project-radius/radius/pkg/azure/radclient"
	"github.com/project-radius/radius/pkg/kubernetes"
	"github.com/project-radius/radius/pkg/radrp/outputresource"
	"github.com/project-radius/radius/pkg/renderers"
	"github.com/project-radius/radius/pkg/renderers/gateway"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

type Renderer struct {
}

// Need a step to take rendered routes to be usable by resource
func (r Renderer) GetDependencyIDs(ctx context.Context, resource renderers.RendererResource) ([]azresources.ResourceID, []azresources.ResourceID, error) {
	route := radclient.HTTPRouteProperties{}
	err := resource.ConvertDefinition(&route)
	if err != nil {
		return nil, nil, err
	}

	if route.Gateway != nil && route.Gateway.Source != nil && *route.Gateway.Source != "" {
		resourceId, err := azresources.Parse(*route.Gateway.Source)
		if err != nil {
			return nil, nil, err
		}
		return []azresources.ResourceID{resourceId}, nil, nil
	}
	return nil, nil, nil
}

func (r Renderer) Render(ctx context.Context, options renderers.RenderOptions) (renderers.RendererOutput, error) {
	route := radclient.HTTPRouteProperties{}
	resource := options.Resource
	dependencies := options.Dependencies

	err := resource.ConvertDefinition(&route)
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	computedValues := map[string]renderers.ComputedValueReference{
		"host": {
			Value: kubernetes.MakeResourceName(resource.ApplicationName, resource.ResourceName),
		},
		"port": {
			Value: GetEffectivePort(route),
		},
		"url": {
			Value: fmt.Sprintf("http://%s:%d", kubernetes.MakeResourceName(resource.ApplicationName, resource.ResourceName), GetEffectivePort(route)),
		},
		"scheme": {
			Value: "http",
		},
	}

	outputs := []outputresource.OutputResource{}

	service := r.makeService(resource, route)
	outputs = append(outputs, service)

	if route.Gateway != nil {
		gatewayId := route.Gateway.Source
		if gatewayId == nil || *gatewayId == "" {
			gatewayClassName := options.Runtime.Gateway.GatewayClass
			if gatewayClassName == "" {
				return renderers.RendererOutput{}, errors.New("gateway class not found")
			}

			defaultGateway := r.createDefaultGateway()
			gatewayK8s := gateway.MakeGateway(ctx, resource, defaultGateway, gatewayClassName)
			outputs = append(outputs, gatewayK8s)

			httpRoute := r.makeHttpRoute(resource, route, kubernetes.MakeResourceName(resource.ApplicationName, resource.ResourceName))
			outputs = append(outputs, httpRoute)
		} else {
			existingGateway := dependencies[*gatewayId]
			httpRoute := r.makeHttpRoute(resource, route, kubernetes.MakeResourceName(resource.ApplicationName, existingGateway.ResourceID.Name()))
			outputs = append(outputs, httpRoute)
		}
	}

	return renderers.RendererOutput{
		Resources:      outputs,
		ComputedValues: computedValues,
	}, nil
}

func (r *Renderer) createDefaultGateway() gateway.Gateway {
	port := 80
	gateway := gateway.Gateway{
		Listeners: map[string]gateway.Listener{
			"http": {
				Port:     &port,
				Protocol: "HTTP",
			},
		},
	}

	return gateway
}

func (r *Renderer) makeService(resource renderers.RendererResource, route radclient.HTTPRouteProperties) outputresource.OutputResource {
	typeParts := strings.Split(resource.ResourceType, "/")
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubernetes.MakeResourceName(resource.ApplicationName, resource.ResourceName),
			Namespace: resource.ApplicationName,
			Labels:    kubernetes.MakeDescriptiveLabels(resource.ApplicationName, resource.ResourceName),
		},
		Spec: corev1.ServiceSpec{
			Selector: kubernetes.MakeRouteSelectorLabels(resource.ApplicationName, typeParts[len(typeParts)-1], resource.ResourceName),
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       resource.ResourceName,
					Port:       int32(GetEffectivePort(route)),
					TargetPort: intstr.FromString(kubernetes.GetShortenedTargetPortName(resource.ApplicationName + typeParts[len(typeParts)-1] + resource.ResourceName)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	return outputresource.NewKubernetesOutputResource(outputresource.LocalIDService, service, service.ObjectMeta)
}

func (r *Renderer) makeHttpRoute(resource renderers.RendererResource, route radclient.HTTPRouteProperties, gatewayName string) outputresource.OutputResource {

	serviceName := kubernetes.MakeResourceName(resource.ApplicationName, resource.ResourceName)
	pathMatch := gatewayv1alpha2.PathMatchPathPrefix
	var rules []gatewayv1alpha2.HTTPRouteRule
	for _, rule := range route.Gateway.Rules {
		// Default to prefix match
		if rule.Path != nil && rule.Path.Type != nil && strings.EqualFold(*rule.Path.Type, "exact") {
			pathMatch = gatewayv1alpha2.PathMatchExact
		}
		port := gatewayv1alpha2.PortNumber(GetEffectivePort(route))
		rules = append(rules, gatewayv1alpha2.HTTPRouteRule{
			Matches: []gatewayv1alpha2.HTTPRouteMatch{
				{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &pathMatch,
						Value: rule.Path.Value,
					},
				},
			},
			BackendRefs: []gatewayv1alpha2.HTTPBackendRef{
				{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(serviceName),
							Port: &port,
						},
					},
				},
			},
		})
	}

	// Add a default rule which maps to the service if none specified
	if len(rules) == 0 {
		path := "/"
		port := gatewayv1alpha2.PortNumber(GetEffectivePort(route))
		rules = append(rules, gatewayv1alpha2.HTTPRouteRule{
			Matches: []gatewayv1alpha2.HTTPRouteMatch{
				{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &pathMatch,
						Value: &path,
					},
				},
			},
			BackendRefs: []gatewayv1alpha2.HTTPBackendRef{
				{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(serviceName),
							Port: &port,
						},
					},
				},
			},
		})
	}
	var hostnames []gatewayv1alpha2.Hostname
	hostname := ""
	if route.Gateway != nil && route.Gateway.Hostname != nil {
		hostname = *route.Gateway.Hostname
	}
	if hostname != "" && hostname != "*" {
		hostnames = append(hostnames, gatewayv1alpha2.Hostname(hostname))
	}

	namespace := gatewayv1alpha2.Namespace("default")
	httpRoute := &gatewayv1alpha2.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gatewayv1alpha2.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubernetes.MakeResourceName(resource.ApplicationName, resource.ResourceName),
			Namespace: resource.ApplicationName,
			Labels:    kubernetes.MakeDescriptiveLabels(resource.ApplicationName, resource.ResourceName),
		},
		Spec: gatewayv1alpha2.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentRef{
					{
						Name:      gatewayv1alpha2.ObjectName(gatewayName),
						Namespace: &namespace,
					},
				},
			},
			Rules:     rules,
			Hostnames: hostnames,
		},
	}

	return outputresource.NewKubernetesOutputResource(outputresource.LocalIDHttpRoute, httpRoute, httpRoute.ObjectMeta)
}
