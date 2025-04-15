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

package injector

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kubeservice-stack/custom-limit-range/pkg/common"
	"github.com/kubeservice-stack/custom-limit-range/pkg/webhook"
)

// log is for logging in this package.
var customlimitrangelog = logf.Log.WithName("customlimitrange-injector")

type PodAnnotator struct {
	Client client.Client
}

// PodAnnotator adds an annotation to every incoming pods.
func (a *PodAnnotator) Default(ctx context.Context, obj runtime.Object) error {
	customlimitrangelog.Info("PodAnnotator", "obj", obj)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	ns := pod.Namespace
	customlimitrangelog.Info("request", "pod", pod, "namespace", ns)

	if ns == "" {
		ns = "default"
	}

	if val, ok := pod.Annotations["customlimitrange.kubernetes.io/limited"]; ok && val == "disable" {
		return nil
	}

	an, err := a.ConfigAnnotation(pod.Annotations, ns)
	if err != nil {
		return err
	}

	pod.Annotations = an
	pod.Namespace = ns

	customlimitrangelog.Info("patch", "pod", pod, "namespace", ns)

	return nil
}

func (a *PodAnnotator) ConfigAnnotation(an map[string]string, namespace string) (map[string]string, error) {
	clrl := &webhook.CustomLimitRangeList{}
	err := a.Client.List(context.Background(), clrl, client.InNamespace(namespace))
	if err != nil {
		customlimitrangelog.Info("Get CustomLimitRange Resource Error", "namespace", namespace, "resource name", common.WebhookName)
		if errors.IsNotFound(err) {
			return an, nil
		}
		return nil, common.ErrMissingConfiguration
	}

	if len(clrl.Items) > 1 {
		customlimitrangelog.Info("Namespace has more than one CustomLimitRange Resource", "count", len(clrl.Items))
		return nil, common.ErrInvalidCustomLimitRangeCountMoreThanOne
	} else if len(clrl.Items) <= 0 {
		// Namespace not found CustomLimitRange Resource
		customlimitrangelog.Info("Namespace not found CustomLimitRange Resource")
		return an, nil
	} else {
		customlimitrangelog.Info("PodAnnotator get CustomLimitRange", "CustomLimitRange", clrl.Items[0].Spec)
		ingress, ok1 := an["kubernetes.io/ingress-bandwidth"]
		if ok1 {
			ig := resource.MustParse(ingress)
			if (!clrl.Items[0].Spec.LRange.Max.Ingress.IsZero() && ig.Value() > clrl.Items[0].Spec.LRange.Max.Ingress.Value()) ||
				(!clrl.Items[0].Spec.LRange.Min.Ingress.IsZero() && ig.Value() < clrl.Items[0].Spec.LRange.Min.Ingress.Value()) {
				return nil, common.ErrInvalidPodSettingBandwidthMaxMin
			}
		} else {
			if !clrl.Items[0].Spec.LRange.Default.Ingress.IsZero() {
				an["kubernetes.io/ingress-bandwidth"] = clrl.Items[0].Spec.LRange.Default.Ingress.String()
			}
		}
		egress, ok2 := an["kubernetes.io/egress-bandwidth"]
		if ok2 {
			eg := resource.MustParse(egress)
			if (!clrl.Items[0].Spec.LRange.Max.Egress.IsZero() && eg.Value() > clrl.Items[0].Spec.LRange.Max.Egress.Value()) ||
				(!clrl.Items[0].Spec.LRange.Min.Egress.IsZero() && eg.Value() < clrl.Items[0].Spec.LRange.Min.Egress.Value()) {
				return nil, common.ErrInvalidPodSettingBandwidthMaxMin
			}
		} else {
			if !clrl.Items[0].Spec.LRange.Default.Egress.IsZero() {
				an["kubernetes.io/egress-bandwidth"] = clrl.Items[0].Spec.LRange.Default.Egress.String()
			}
		}

	}

	return an, nil
}
