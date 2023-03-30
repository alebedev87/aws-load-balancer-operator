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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:Enum=AWSShield;AWSWAFv1;AWSWAFv2
type AWSAddon string

const (
	AWSAddonShield AWSAddon = "AWSShield"
	AWSAddonWAFv1  AWSAddon = "AWSWAFv1"
	AWSAddonWAFv2  AWSAddon = "AWSWAFv2"
)

// +kubebuilder:validation:Enum=Auto;Manual
type SubnetTaggingPolicy string

const (

	// AutoSubnetTaggingPolicy enables automatic subnet tagging.
	AutoSubnetTaggingPolicy SubnetTaggingPolicy = "Auto"

	// ManualSubnetTaggingPolicy disables automatic subnet tagging.
	ManualSubnetTaggingPolicy SubnetTaggingPolicy = "Manual"
)

// AWSLoadBalancerControllerSpec defines the desired state of AWSLoadBalancerController
type AWSLoadBalancerControllerSpec struct {

	// subnetTagging describes how resource tagging will be done by the operator.
	//
	// When in "Auto", the operator will detect the subnets where the load balancers
	// will be provisioned and have the required resource tags on them. Whereas when
	// set to "Manual", this responsibility lies on the user. The subnets tagged by the operator
	// will be removed when transitioning from "Auto" to "Manual". Whereas the manually tagged subnets
	// will be left intact when transitioning from "Manual" to "Auto".
	//
	// +kubebuilder:default:=Auto
	// +kubebuilder:validation:Optional
	// +optional
	SubnetTagging SubnetTaggingPolicy `json:"subnetTagging,omitempty"`

	// additionalResourceTags are the AWS Tags that will be applied to all AWS resources managed by this
	// controller. The addition of new tags as well as the update or removal of the existing tags
	// will be propagated to the AWS resources. AWS supports a maximum of 50 tags per resource.
	// AWSLoadBalancerController reserves 3 tags for its use, the rest is split in half between
	// the tag annotation which can be set on the ingress or service and this field.
	// Each tag key must be unique.
	//
	// +kubebuilder:validation:MaxItems=24
	// +kubebuilder:validation:Optional
	// +optional
	// +listType=map
	// +listMapKey=key
	// +patchMergeKey=key
	// +patchStrategy=merge
	AdditionalResourceTags []AWSResourceTag `json:"additionalResourceTags,omitempty" patchStrategy:"merge" patchMergeKey:"key"`

	// ingressClass specifies the Ingress class which the controller will reconcile.
	// This Ingress class will be created unless it already exists.
	// The value will default to "alb".
	//
	// The defaulting to "alb" is necessary so that this controller can function as expected
	// in parallel with openshift-router, for more info see
	// https://github.com/openshift/enhancements/blob/master/enhancements/ingress/aws-load-balancer-operator.md#parallel-operation-of-the-openshift-router-and-lb-controller
	//
	// +kubebuilder:default:=alb
	// +kubebuilder:validation:Optional
	// +optional
	IngressClass string `json:"ingressClass,omitempty"`

	// config specifies further customization options for the controller's deployment spec.
	//
	// +kubebuilder:validation:Optional
	// +optional
	Config *AWSLoadBalancerDeploymentConfig `json:"config,omitempty"`

	// enabledAddons describes the AWS services that can be integrated with
	// the AWS Load Balancers created by the controller.
	// Note that when an addon which was previously enabled is disabled
	// the controller does not remove the existing addon attachment for the provisioned load balancers.
	//
	// +kubebuilder:validation:Optional
	// +optional
	EnabledAddons []AWSAddon `json:"enabledAddons,omitempty"`

	// credentials is a reference to a secret containing
	// the AWS credentials to be used by the controller.
	// The secret is required to have `credentials` data key
	// containing the AWS CLI credentials file (static or STS),
	// for examples, see https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html.
	// under `credentials` data key.
	// The secret is required to be in the operator namespace.
	// If this field is empty - the credentials will be
	// requested using the Cloud Credentials API,
	// see https://docs.openshift.com/container-platform/4.11/authentication/managing_cloud_provider_credentials/about-cloud-credential-operator.html.
	//
	// +kubebuilder:validation:Optional
	// +optional
	Credentials *SecretReference `json:"credentials,omitempty"`
}

// AWSResourceTag is a tag to apply to AWS resources created by the controller.
type AWSResourceTag struct {
	// key is the key of the tag.
	// See https://docs.aws.amazon.com/tag-editor/latest/userguide/tagging.html#tag-conventions
	// for information on the tagging conventions.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z_.:/=+-@]+$`
	// +required
	Key string `json:"key"`

	// value is the value of the tag.
	// See https://docs.aws.amazon.com/tag-editor/latest/userguide/tagging.html#tag-conventions
	// for information on the tagging conventions.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z_.:/=+-@]*$`
	// +required
	Value string `json:"value"`
}

type AWSLoadBalancerDeploymentConfig struct {
	// +kubebuilder:default:=2
	// +kubebuilder:validation:Optional
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
}

// SecretReference contains the information to let you locate the desired secret.
// Secret is required to be in the operator namespace.
type SecretReference struct {
	// name is the name of the secret.
	//
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
}

// AWSLoadBalancerControllerStatus defines the observed state of AWSLoadBalancerController.
type AWSLoadBalancerControllerStatus struct {
	// conditions is a list of operator-specific conditions
	// and their status.
	//
	// +kubebuilder:validation:Optional
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// observedGeneration is the most recent generation observed.
	//
	// +kubebuilder:validation:Optional
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// subnets contains details of the subnets of the cluster
	//
	// +kubebuilder:validation:Optional
	// +optional
	Subnets *AWSLoadBalancerControllerStatusSubnets `json:"subnets,omitempty"`

	// ingressClass is the current default Ingress class.
	//
	// +kubebuilder:validation:Optional
	// +optional
	IngressClass string `json:"ingressClass,omitempty"`
}

type AWSLoadBalancerControllerStatusSubnets struct {
	// subnetTagging indicates the current status of the subnet tags
	// +kubebuilder:validation:Optional
	// +optional
	SubnetTagging SubnetTaggingPolicy `json:"subnetTagging,omitempty"`

	// internal is the list of subnet ids which have the tag `kubernetes.io/role/internal-elb`
	//
	// +kubebuilder:validation:Optional
	// +optional
	Internal []string `json:"internal,omitempty"`

	// public is the list of subnet ids which have the tag `kubernetes.io/role/elb`
	//
	// +kubebuilder:validation:Optional
	// +optional
	Public []string `json:"public,omitempty"`

	// tagged is the list of subnet ids which have been tagged by the operator
	//
	// +kubebuilder:validation:Optional
	// +optional
	Tagged []string `json:"tagged,omitempty"`

	// untagged is the list of subnet ids which do not have any role tags
	//
	// +kubebuilder:validation:Optional
	// +optional
	Untagged []string `json:"untagged,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:storageversion

// AWSLoadBalancerController is the Schema for the awsloadbalancercontrollers API
type AWSLoadBalancerController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSLoadBalancerControllerSpec   `json:"spec,omitempty"`
	Status AWSLoadBalancerControllerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AWSLoadBalancerControllerList contains a list of AWSLoadBalancerController
type AWSLoadBalancerControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSLoadBalancerController `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSLoadBalancerController{}, &AWSLoadBalancerControllerList{})
}
