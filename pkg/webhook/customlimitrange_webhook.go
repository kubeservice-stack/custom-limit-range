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
	"github.com/kubeservice-stack/custom-limit-range/pkg/common"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	wk "sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var customlimitrangelog = logf.Log.WithName("customlimitrange-resource")

func (r *CustomLimitRange) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

var _ wk.Defaulter = &CustomLimitRange{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CustomLimitRange) Default() {
	customlimitrangelog.Info("default", "name", r.Name, "request", r)
	customlimitrangelog.Info("Spec", "spec", r.Spec)
}

var _ wk.Validator = &CustomLimitRange{}
var minRsrc = resource.MustParse("1k")
var maxRsrc = resource.MustParse("1P")

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CustomLimitRange) ValidateCreate() error {
	customlimitrangelog.Info("validate create", "name", r.Name, "request", r)
	var allErrs field.ErrorList
	err := bandwidthValidateIsReasonable(r.Spec.LRange.Min, r.Spec.LRange.Default, r.Spec.LRange.Max)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("LRange"),
			r.Spec.LRange,
			err.Error()))
	}
	customlimitrangelog.Info("validate bandwidthValidateIsReasonable", "err", err, "field.ErrorList", allErrs)
	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CustomLimitRange) ValidateUpdate(old runtime.Object) error {
	customlimitrangelog.Info("validate update", "name", r.Name, "request", r)
	var allErrs field.ErrorList
	err := bandwidthValidateIsReasonable(r.Spec.LRange.Min, r.Spec.LRange.Default, r.Spec.LRange.Max)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("LRange"),
			r.Spec.LRange,
			err.Error()))
	}
	customlimitrangelog.Info("validate bandwidthValidateIsReasonable", "err", err, "field.ErrorList", allErrs)
	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CustomLimitRange) ValidateDelete() error {
	customlimitrangelog.Info("validate delete", "name", r.Name)
	return nil
}

func bandwidthValidateIsReasonable(min, def, max CustomItems) error {
	if err := bandwidthValidate(min); err != nil {
		return err
	}

	if err := bandwidthValidate(def); err != nil {
		return err
	}

	if err := bandwidthValidate(max); err != nil {
		return err
	}

	if !min.Egress.IsZero() && !def.Egress.IsZero() && min.Egress.Value() > def.Egress.Value() {
		return common.ErrInvalidBandwidthMaxMin
	}
	if !min.Ingress.IsZero() && !def.Ingress.IsZero() && min.Ingress.Value() > def.Ingress.Value() {
		return common.ErrInvalidBandwidthMaxMin
	}

	if !min.Egress.IsZero() && !max.Egress.IsZero() && min.Egress.Value() > max.Egress.Value() {
		return common.ErrInvalidBandwidthMaxMin
	}
	if !min.Ingress.IsZero() && !max.Ingress.IsZero() && min.Ingress.Value() > max.Ingress.Value() {
		return common.ErrInvalidBandwidthMaxMin
	}

	if !def.Egress.IsZero() && !max.Egress.IsZero() && def.Egress.Value() > max.Egress.Value() {
		return common.ErrInvalidBandwidthMaxMin
	}
	if !def.Ingress.IsZero() && !max.Ingress.IsZero() && def.Ingress.Value() > max.Ingress.Value() {
		return common.ErrInvalidBandwidthMaxMin
	}

	return nil
}

func bandwidthValidate(item CustomItems) error {
	if !item.Egress.IsZero() {
		if err := validateBandwidthIsReasonable(item.Egress); err != nil {
			return err
		}
	}
	if !item.Ingress.IsZero() {
		if err := validateBandwidthIsReasonable(item.Ingress); err != nil {
			return err
		}
	}

	return nil
}

func validateBandwidthIsReasonable(rsrc resource.Quantity) error {
	if rsrc.Value() < minRsrc.Value() {
		return common.ErrInvalidBandwidthRange
	}
	if rsrc.Value() > maxRsrc.Value() {
		return common.ErrInvalidBandwidthRange
	}
	return nil
}
