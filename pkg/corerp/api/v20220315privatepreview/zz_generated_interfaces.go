//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220315privatepreview

// ApplicationExtensionClassification provides polymorphic access to related types.
// Call the interface's GetApplicationExtension() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ApplicationExtension, *ApplicationKubernetesMetadataExtension, *ApplicationKubernetesNamespaceExtension
type ApplicationExtensionClassification interface {
	ExtensionClassification
	// GetApplicationExtension returns the ApplicationExtension content of the underlying type.
	GetApplicationExtension() *ApplicationExtension
}

// ContainerExtensionClassification provides polymorphic access to related types.
// Call the interface's GetContainerExtension() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ContainerExtension, *ContainerKubernetesMetadataExtension, *DaprSidecarExtension, *ManualScalingExtension
type ContainerExtensionClassification interface {
	ExtensionClassification
	// GetContainerExtension returns the ContainerExtension content of the underlying type.
	GetContainerExtension() *ContainerExtension
}

// EnvironmentComputeClassification provides polymorphic access to related types.
// Call the interface's GetEnvironmentCompute() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EnvironmentCompute, *KubernetesCompute
type EnvironmentComputeClassification interface {
	// GetEnvironmentCompute returns the EnvironmentCompute content of the underlying type.
	GetEnvironmentCompute() *EnvironmentCompute
}

// EnvironmentExtensionClassification provides polymorphic access to related types.
// Call the interface's GetEnvironmentExtension() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EnvironmentExtension, *EnvironmentKubernetesMetadataExtension
type EnvironmentExtensionClassification interface {
	ExtensionClassification
	// GetEnvironmentExtension returns the EnvironmentExtension content of the underlying type.
	GetEnvironmentExtension() *EnvironmentExtension
}

// EnvironmentRecipePropertiesClassification provides polymorphic access to related types.
// Call the interface's GetEnvironmentRecipeProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *BicepRecipeProperties, *EnvironmentRecipeProperties, *TerraformRecipeProperties
type EnvironmentRecipePropertiesClassification interface {
	// GetEnvironmentRecipeProperties returns the EnvironmentRecipeProperties content of the underlying type.
	GetEnvironmentRecipeProperties() *EnvironmentRecipeProperties
}

// ExtensionClassification provides polymorphic access to related types.
// Call the interface's GetExtension() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ApplicationExtension, *ApplicationKubernetesNamespaceExtension, *ContainerExtension, *ContainerKubernetesMetadataExtension,
// - *DaprSidecarExtension, *EnvironmentExtension, *Extension, *ManualScalingExtension
type ExtensionClassification interface {
	// GetExtension returns the Extension content of the underlying type.
	GetExtension() *Extension
}

// HealthProbePropertiesClassification provides polymorphic access to related types.
// Call the interface's GetHealthProbeProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ExecHealthProbeProperties, *HTTPGetHealthProbeProperties, *HealthProbeProperties, *TCPHealthProbeProperties
type HealthProbePropertiesClassification interface {
	// GetHealthProbeProperties returns the HealthProbeProperties content of the underlying type.
	GetHealthProbeProperties() *HealthProbeProperties
}

// VolumeClassification provides polymorphic access to related types.
// Call the interface's GetVolume() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EphemeralVolume, *PersistentVolume, *Volume
type VolumeClassification interface {
	// GetVolume returns the Volume content of the underlying type.
	GetVolume() *Volume
}

// VolumePropertiesClassification provides polymorphic access to related types.
// Call the interface's GetVolumeProperties() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AzureKeyVaultVolumeProperties, *VolumeProperties
type VolumePropertiesClassification interface {
	// GetVolumeProperties returns the VolumeProperties content of the underlying type.
	GetVolumeProperties() *VolumeProperties
}

