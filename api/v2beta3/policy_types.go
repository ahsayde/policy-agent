/*
Copyright 2022.

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

package v2beta3

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PolicyResourceName       = "policies"
	PolicyKind               = "Policy"
	PolicyListKind           = "PolicyList"
	TenancyTag               = "tenancy"
	PolicyKubernetesProvider = "kubernetes"
	PolicyTerraformProvider  = "terraform"
)

var (
	PolicyGroupVersionResource = GroupVersion.WithResource(PolicyResourceName)
)

// PolicyParameters defines a needed input in a policy
type PolicyParameters struct {
	// Name is a descriptive name of a policy parameter
	Name string `json:"name"`
	// Type is the type of that parameter, integer, string,...
	Type string `json:"type"`
	// Required specifies if this is a necessary value or not
	Required bool `json:"required"`
	// +optional
	// Value is the value for that parameter
	Value *apiextensionsv1.JSON `json:"value,omitempty"`
}

// PolicyTargets are filters used to determine which resources should be evaluated against a policy
type PolicyTargets struct {
	// Kinds is a list of Kubernetes kinds that are supported by this policy
	Kinds []string `json:"kinds"`
	// +optional
	// Labels is a list of Kubernetes labels that are needed to evaluate the policy against a resource
	// this filter is statisfied if only one label existed, using * for value make it so it will match if the key exists regardless of its value
	Labels []map[string]string `json:"labels"`
	// +optional
	// Namespaces is a list of Kubernetes namespaces that a resource needs to be a part of to evaluate against this policy
	Namespaces []string `json:"namespaces"`
}

type PolicyStandard struct {
	// ID idenitifer of the standarad
	ID string `json:"id"`
	// Controls standard controls
	Controls []string `json:"controls,omitempty"`
}

// PolicyExclusions are the structure which resources should not be evaluated against the policy
type PolicyExclusions struct {
	// +optional
	// Namespaces is a list of Kubernetes namespaces that a resource needs to be a part of to excluded from this policy
	Namespaces []string `json:"namespaces"`
	// +optional
	// Resources is a list of Kubernetes resources that are excluded by this policy (namespace/name)
	Resources []string `json:"resources"`
	// +optional
	// Labels is a list of Kubernetes labels that are needed to excluded the policy against a resource
	// this filter is statisfied if only one label existed, using * for value make it so it will match if the key exists regardless of its value
	Labels []map[string]string `json:"labels"`
}

// PolicySpec defines the desired state of Policy
// It describes all that is needed to evaluate a resource against a rego code
// +kubebuilder:object:generate:true
type PolicySpec struct {
	// Name is the policy name
	Name string `json:"name"`
	// ID is the policy unique identifier
	ID string `json:"id"`
	// Code contains the policy rego code
	Code string `json:"code"`
	// +optional
	// +kubebuilder:default:=true
	// Enforce flag to define whether a policy is enforced via the admission controller or just audited for a violation (default: true)
	Enforce bool `json:"enforce,omitempty"`
	// +optional
	// Parameters are the inputs needed for the policy validation
	Parameters []PolicyParameters `json:"parameters,omitempty"`
	// +optional
	// Targets describes the required metadata that needs to be matched to evaluate a resource against the policy
	// all values specified need to exist in the resource to be considered for evaluation
	Targets PolicyTargets `json:"targets,omitempty"`
	// Description is a summary of what that policy validates
	Description string `json:"description"`
	// HowToSolve is a description of the steps required to solve the issues reported by the policy
	HowToSolve string `json:"how_to_solve"`
	// Category specifies under which grouping this policy should be included
	Category string `json:"category"`
	// +optional
	// Tags is a list of tags associated with that policy
	Tags []string `json:"tags,omitempty"`
	// +kubebuilder:validation:Enum=low;medium;high
	// Severity is a measure of the impact of that policy, can be low, medium or high
	Severity string `json:"severity"`
	// +optional
	// Standards is a list of policy standards that this policy falls under
	Standards []PolicyStandard `json:"standards"`
	//+optional
	//+kubebuilder:default:=kubernetes
	//+kubebuilder:validation:Enum=kubernetes;terraform
	// Provider is policy provider, can be kubernetes, terraform
	Provider string `json:"provider"`

	//+optional
	//+kubebuilder:default:=false
	// Mutate is a flag that indicates whether to enable mutation of resources violating this policy or not
	Mutate bool `json:"mutate"`

	// +optional
	// Exclude describes the policy exclusions on (Namespaces, Labels, Resources)
	// Select one or more by defining the exclusion list
	Exclude PolicyExclusions `json:"exclude,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:printcolumn:name="Severity",type=string,JSONPath=`.spec.severity`
//+kubebuilder:printcolumn:name="Category",type=string,JSONPath=`.spec.category`
//+kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
//+kubebuilder:printcolumn:name="Enforced",type=string,JSONPath=`.spec.enforce`
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:storageversion

// Policy is the Schema for the policies API
type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PolicySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// PolicyList contains a list of Policy
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Policy `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&Policy{},
		&PolicyList{},
	)
}
