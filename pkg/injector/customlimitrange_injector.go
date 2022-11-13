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
	"encoding/json"
	"net/http"

	"github.com/kubeservice-stack/custom-limit-range/pkg/common"
	"github.com/kubeservice-stack/custom-limit-range/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var customlimitrangelog = logf.Log.WithName("customlimitrange-injector")

type PodAnnotator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func NewPodAnnotatorMutate(c client.Client) admission.Handler {
	return &PodAnnotator{Client: c}
}

// PodAnnotator adds an annotation to every incoming pods.
func (a *PodAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	customlimitrangelog.Info("PodAnnotator", "req", req)
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	ns := req.Namespace
	customlimitrangelog.Info("request", "pod", pod, "namespace", ns, "req.Namespace", req.Namespace)

	if ns == "" {
		ns = "default"
	}

	if val, ok := pod.Annotations["customlimitrange.kubernetes.io/limited"]; ok && val == "disable" {
		return admission.Allowed("Pod customlimitrange.kubernetes.io/limited setting Disable, Pass CustomLimitRange")
	}

	an, err := a.ConfigAnnotation(pod.Annotations, ns)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	pod.Annotations = an
	pod.Namespace = ns

	customlimitrangelog.Info("patch", "pod", pod, "namespace", ns)

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// PodAnnotator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PodAnnotator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
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
