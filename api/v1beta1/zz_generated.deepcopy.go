//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSLoadBalancerController) DeepCopyInto(out *AWSLoadBalancerController) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSLoadBalancerController.
func (in *AWSLoadBalancerController) DeepCopy() *AWSLoadBalancerController {
	if in == nil {
		return nil
	}
	out := new(AWSLoadBalancerController)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSLoadBalancerController) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSLoadBalancerControllerList) DeepCopyInto(out *AWSLoadBalancerControllerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AWSLoadBalancerController, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSLoadBalancerControllerList.
func (in *AWSLoadBalancerControllerList) DeepCopy() *AWSLoadBalancerControllerList {
	if in == nil {
		return nil
	}
	out := new(AWSLoadBalancerControllerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSLoadBalancerControllerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSLoadBalancerControllerSpec) DeepCopyInto(out *AWSLoadBalancerControllerSpec) {
	*out = *in
	if in.AdditionalResourceTags != nil {
		in, out := &in.AdditionalResourceTags, &out.AdditionalResourceTags
		*out = make([]AWSResourceTag, len(*in))
		copy(*out, *in)
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(AWSLoadBalancerDeploymentConfig)
		**out = **in
	}
	if in.EnabledAddons != nil {
		in, out := &in.EnabledAddons, &out.EnabledAddons
		*out = make([]AWSAddon, len(*in))
		copy(*out, *in)
	}
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(v1.SecretNameReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSLoadBalancerControllerSpec.
func (in *AWSLoadBalancerControllerSpec) DeepCopy() *AWSLoadBalancerControllerSpec {
	if in == nil {
		return nil
	}
	out := new(AWSLoadBalancerControllerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSLoadBalancerControllerStatus) DeepCopyInto(out *AWSLoadBalancerControllerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Subnets != nil {
		in, out := &in.Subnets, &out.Subnets
		*out = new(AWSLoadBalancerControllerStatusSubnets)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSLoadBalancerControllerStatus.
func (in *AWSLoadBalancerControllerStatus) DeepCopy() *AWSLoadBalancerControllerStatus {
	if in == nil {
		return nil
	}
	out := new(AWSLoadBalancerControllerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSLoadBalancerControllerStatusSubnets) DeepCopyInto(out *AWSLoadBalancerControllerStatusSubnets) {
	*out = *in
	if in.Internal != nil {
		in, out := &in.Internal, &out.Internal
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Public != nil {
		in, out := &in.Public, &out.Public
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Tagged != nil {
		in, out := &in.Tagged, &out.Tagged
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Untagged != nil {
		in, out := &in.Untagged, &out.Untagged
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSLoadBalancerControllerStatusSubnets.
func (in *AWSLoadBalancerControllerStatusSubnets) DeepCopy() *AWSLoadBalancerControllerStatusSubnets {
	if in == nil {
		return nil
	}
	out := new(AWSLoadBalancerControllerStatusSubnets)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSLoadBalancerDeploymentConfig) DeepCopyInto(out *AWSLoadBalancerDeploymentConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSLoadBalancerDeploymentConfig.
func (in *AWSLoadBalancerDeploymentConfig) DeepCopy() *AWSLoadBalancerDeploymentConfig {
	if in == nil {
		return nil
	}
	out := new(AWSLoadBalancerDeploymentConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSResourceTag) DeepCopyInto(out *AWSResourceTag) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSResourceTag.
func (in *AWSResourceTag) DeepCopy() *AWSResourceTag {
	if in == nil {
		return nil
	}
	out := new(AWSResourceTag)
	in.DeepCopyInto(out)
	return out
}
