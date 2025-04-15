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
	"context"
	"testing"

	"github.com/kubeservice-stack/custom-limit-range/pkg/common"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestValidateBandwidthIsReasonable(t *testing.T) {
	assert := assert.New(t)
	err := validateBandwidthIsReasonable(resource.MustParse("124"))
	assert.ErrorIs(err, common.ErrInvalidBandwidthRange)

	err = validateBandwidthIsReasonable(resource.MustParse("0.1k"))
	assert.ErrorIs(err, common.ErrInvalidBandwidthRange)

	err = validateBandwidthIsReasonable(resource.MustParse("1124"))
	assert.Nil(err)

	err = validateBandwidthIsReasonable(resource.MustParse("1.1P"))
	assert.ErrorIs(err, common.ErrInvalidBandwidthRange)

	err = validateBandwidthIsReasonable(resource.MustParse("1001T"))
	assert.ErrorIs(err, common.ErrInvalidBandwidthRange)
}

func TestBandwidthValidate(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	testCases := []struct {
		name     string
		ingress  resource.Quantity
		egress   resource.Quantity
		expected error
	}{
		{
			name:     "InReasonable",
			ingress:  resource.MustParse("1P"),
			egress:   resource.MustParse("1P"),
			expected: nil,
		},
		{
			name:     "NotInReasonable",
			ingress:  resource.MustParse("1k"),
			egress:   resource.MustParse("1P"),
			expected: nil,
		},
		{
			name:     "NotIn2Reasonable",
			ingress:  resource.MustParse("0.1k"),
			egress:   resource.MustParse("1P"),
			expected: common.ErrInvalidBandwidthRange,
		},
		{
			name:     "NotIn3Reasonable",
			ingress:  resource.MustParse("0.1k"),
			egress:   resource.MustParse("1.01P"),
			expected: common.ErrInvalidBandwidthRange,
		},
		{
			name:     "NotIn4Reasonable",
			ingress:  resource.MustParse("2k"),
			egress:   resource.MustParse("1.01P"),
			expected: common.ErrInvalidBandwidthRange,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := bandwidthValidate(CustomItems{Ingress: tc.ingress, Egress: tc.egress})
			assert.ErrorIs(err, tc.expected, tc.name)
		})
	}

	err := bandwidthValidate(CustomItems{Ingress: resource.MustParse("2k")})
	assert.Nil(err)
	err = bandwidthValidate(CustomItems{Ingress: resource.MustParse("2P")})
	assert.ErrorIs(err, common.ErrInvalidBandwidthRange)
}

func TestBandwidthValidateIS(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	testCases := []struct {
		name     string
		mix      CustomItems
		def      CustomItems
		max      CustomItems
		expected error
	}{
		{
			name:     "InReasonableIsEmpty",
			mix:      CustomItems{},
			def:      CustomItems{},
			max:      CustomItems{},
			expected: nil,
		},
		{
			name:     "InReasonableIsOne",
			mix:      CustomItems{Ingress: resource.MustParse("1P")},
			def:      CustomItems{},
			max:      CustomItems{},
			expected: nil,
		},
		{
			name:     "InReasonableIsOneRangeOver",
			mix:      CustomItems{Ingress: resource.MustParse("2P")},
			def:      CustomItems{},
			max:      CustomItems{},
			expected: common.ErrInvalidBandwidthRange,
		},
		{
			name:     "InReasonableIsTwo",
			mix:      CustomItems{Ingress: resource.MustParse("1P")},
			def:      CustomItems{Egress: resource.MustParse("1P")},
			max:      CustomItems{},
			expected: nil,
		},
		{
			name:     "InReasonableIsTwo1",
			mix:      CustomItems{Ingress: resource.MustParse("1P")},
			def:      CustomItems{Ingress: resource.MustParse("1P")},
			max:      CustomItems{},
			expected: nil,
		},
		{
			name:     "InReasonableIsTwo3",
			mix:      CustomItems{Ingress: resource.MustParse("1P")},
			def:      CustomItems{Ingress: resource.MustParse("1G")},
			max:      CustomItems{},
			expected: common.ErrInvalidBandwidthMaxMin,
		},
		{
			name:     "InReasonableIsTwo3",
			mix:      CustomItems{Egress: resource.MustParse("1P")},
			def:      CustomItems{Egress: resource.MustParse("1G")},
			max:      CustomItems{},
			expected: common.ErrInvalidBandwidthMaxMin,
		},
		{
			name:     "InReasonableIsTwo3",
			mix:      CustomItems{Egress: resource.MustParse("1P")},
			def:      CustomItems{},
			max:      CustomItems{Egress: resource.MustParse("1G")},
			expected: common.ErrInvalidBandwidthMaxMin,
		},
		{
			name:     "InReasonableIsTwo3",
			mix:      CustomItems{Egress: resource.MustParse("1M")},
			def:      CustomItems{Egress: resource.MustParse("1k")},
			max:      CustomItems{Egress: resource.MustParse("1G")},
			expected: common.ErrInvalidBandwidthMaxMin,
		},
		{
			name:     "InReasonableIsTwo3",
			mix:      CustomItems{Egress: resource.MustParse("1M")},
			def:      CustomItems{Ingress: resource.MustParse("1k")},
			max:      CustomItems{Egress: resource.MustParse("1G")},
			expected: nil,
		},
		{
			name:     "InReasonableIsTwo3",
			mix:      CustomItems{Egress: resource.MustParse("1M"), Ingress: resource.MustParse("1M")},
			def:      CustomItems{Ingress: resource.MustParse("1k")},
			max:      CustomItems{Egress: resource.MustParse("1G")},
			expected: common.ErrInvalidBandwidthMaxMin,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := bandwidthValidateIsReasonable(tc.mix, tc.def, tc.max)
			assert.ErrorIs(err, tc.expected, tc.name)
		})
	}
}

func TestCustomLimitRangeIsEmpty(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	c := &CustomLimitRange{}
	ctx := context.Background()
	_, err := c.ValidateCreate(ctx, nil)
	assert.Nil(err)
	_, err = c.ValidateDelete(ctx, c)
	assert.Nil(err)
	_, err = c.ValidateUpdate(ctx, c, c)
	assert.Nil(err)
	d := c.DeepCopy()
	assert.Equal(c, d)
}

func TestCustomLimitRangeNotEmpty(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	c := &CustomLimitRange{
		Spec: CustomLimitRangeSpec{
			LRange: LimitRange{
				Max: CustomItems{
					Ingress: resource.MustParse("1P"),
					Egress:  resource.MustParse("10M"),
				},
				Min: CustomItems{
					Ingress: resource.MustParse("1G"),
					Egress:  resource.MustParse("1k"),
				},
				Default: CustomItems{
					Ingress: resource.MustParse("1T"),
					Egress:  resource.MustParse("1M"),
				},
			},
		},
	}

	ctx := context.Background()
	_, err := c.ValidateCreate(ctx, nil)
	assert.Nil(err)
	_, err = c.ValidateDelete(ctx, c)
	assert.Nil(err)
	_, err = c.ValidateUpdate(ctx, c, c)
	assert.Nil(err)

	d := c.DeepCopy()
	assert.Equal(c, d)
}
