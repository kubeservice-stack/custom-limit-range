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

package common

import (
	"errors"
)

const (
	WebhookName    = "CustomLimitRange"
	WebhookEnable  = "enabled"
	WebhookDisable = "disabled"
	WebhookVersion = "v1"

	WebhookUrlPath = "/node/affinity"

	WebhookPodDisable = "customlimitrange.kubernetes.io/limited"
)

var (
	ErrInvalidAdmissionReview    = errors.New("invalid admission review error")
	ErrInvalidAdmissionReviewObj = errors.New("invalid admission review object error")
	ErrFailedToCreatePatch       = errors.New("failed to create patch")
	ErrMissingConfiguration      = errors.New("missing configuration")
	ErrInvalidConfiguration      = errors.New("invalid configuration error")

	ErrInvalidBandwidthRange                   = errors.New("resource is unreasonably small (< 1kbit) or large (> 1Pbit)")
	ErrInvalidBandwidthMaxMin                  = errors.New("resource must min <= default <= max")
	ErrInvalidPodSettingBandwidthMaxMin        = errors.New("pod annotation must:  min <= [kubernetes.io/ingress-bandwidth]/[kubernetes.io/egress-bandwidth] <= max")
	ErrInvalidCustomLimitRangeCountMoreThanOne = errors.New("Namespace has more than one CustomLimitRange Resource")
)
