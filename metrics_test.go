/*
 Copyright © 2021-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

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

package gounity

import (
	"context"
	"fmt"
	"testing"

	"github.com/dell/gounity/mocks"
	"github.com/dell/gounity/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteRealTimeMetricsQuery(t *testing.T) {
	fmt.Println("Begin - Delete Real Time Metrics Query Test")
	testConf.metricsAPI.client.api.(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()
	queryID := 12345

	testConf.metricsAPI.client.api.(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := testConf.metricsAPI.DeleteRealTimeMetricsQuery(ctx, queryID)
	fmt.Println("Error:", err)
	if err != nil {
		t.Fatalf("Delete Real Time Metrics Query failed: %v", err)
	}

	testConf.metricsAPI.client.api.(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("delete failed")).Once()

	err = testConf.metricsAPI.DeleteRealTimeMetricsQuery(ctx, queryID)
	if err == nil {
		t.Fatalf("Delete Real Time Metrics Query negative case failed: %v", err)
	}

	fmt.Println("Delete Real Time Metrics Query Test - Successful")
}

func TestGetMetricsCollection(t *testing.T) {
	fmt.Println("Begin - Get Metrics Collection Test")
	testConf.metricsAPI.client.api.(*mocks.Client).ExpectedCalls = nil
	ctx := context.Background()
	queryID := 12345

	metricsQueryResult := &types.MetricQueryResult{}
	testConf.metricsAPI.client.api.(*mocks.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*types.MetricQueryResult)
		if resp != nil {
			*resp = *metricsQueryResult
		}
	}).Once()

	result, err := testConf.metricsAPI.GetMetricsCollection(ctx, queryID)
	fmt.Println("Metrics Query Result:", prettyPrintJSON(result), "Error:", err)
	if err != nil {
		t.Fatalf("Get Metrics Collection failed: %v", err)
	}
	assert.NotNil(t, result)

	fmt.Println("Get Metrics Collection Test - Successful")
}
