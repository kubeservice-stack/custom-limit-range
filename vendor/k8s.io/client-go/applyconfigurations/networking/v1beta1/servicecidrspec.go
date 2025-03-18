/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

// ServiceCIDRSpecApplyConfiguration represents a declarative configuration of the ServiceCIDRSpec type for use
// with apply.
type ServiceCIDRSpecApplyConfiguration struct {
	CIDRs []string `json:"cidrs,omitempty"`
}

// ServiceCIDRSpecApplyConfiguration constructs a declarative configuration of the ServiceCIDRSpec type for use with
// apply.
func ServiceCIDRSpec() *ServiceCIDRSpecApplyConfiguration {
	return &ServiceCIDRSpecApplyConfiguration{}
}

// WithCIDRs adds the given value to the CIDRs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the CIDRs field.
func (b *ServiceCIDRSpecApplyConfiguration) WithCIDRs(values ...string) *ServiceCIDRSpecApplyConfiguration {
	for i := range values {
		b.CIDRs = append(b.CIDRs, values[i])
	}
	return b
}
