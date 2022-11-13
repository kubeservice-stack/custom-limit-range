/*
Copyright 2022 The KubeService-Stack Authors.

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

package webhook

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "custom.cmss.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

type CustomItems struct {
	Ingress resource.Quantity `json:"ingress-bandwidth,omitempty"`
	Egress  resource.Quantity `json:"egress-bandwidth,omitempty"`
}

type LimitRange struct {
	Type    string      `json:"type"`
	Max     CustomItems `json:"max,omitempty"`
	Min     CustomItems `json:"min,omitempty"`
	Default CustomItems `json:"default,omitempty"`
}

// CustomLimitRangeSpec defines the desired state of CustomLimitRange
type CustomLimitRangeSpec struct {
	LRange LimitRange `json:"limitrange"`
}

// CustomLimitRangeStatus defines the observed state of CustomLimitRange
type CustomLimitRangeStatus struct {
}

// CustomLimitRange is the Schema for the customlimitranges API
type CustomLimitRange struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomLimitRangeSpec   `json:"spec"`
	Status CustomLimitRangeStatus `json:"status,omitempty"`
}

// CustomLimitRangeList contains a list of CustomLimitRange
type CustomLimitRangeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomLimitRange `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomLimitRange{}, &CustomLimitRangeList{})
}
