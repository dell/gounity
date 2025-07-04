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
	"errors"
	"fmt"
	"testing"

	mocksapi "github.com/dell/gounity/mocks/api"
	"github.com/dell/gounity/apitypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteRealTimeMetricsQuery(t *testing.T) {
	fmt.Println("Begin - Delete Real Time Metrics Query Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	queryID := 12345

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := testConf.client.DeleteRealTimeMetricsQuery(ctx, queryID)
	fmt.Println("Error:", err)
	if err != nil {
		t.Fatalf("Delete Real Time Metrics Query failed: %v", err)
	}

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("delete failed")).Once()

	err = testConf.client.DeleteRealTimeMetricsQuery(ctx, queryID)
	if err == nil {
		t.Fatalf("Delete Real Time Metrics Query negative case failed: %v", err)
	}

	fmt.Println("Delete Real Time Metrics Query Test - Successful")
}

func TestGetMetricsCollection(t *testing.T) {
	fmt.Println("Begin - Get Metrics Collection Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()
	queryID := 12345

	metricsQueryResult := &apitypes.MetricQueryResult{}
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apitypes.MetricQueryResult)
		if resp != nil {
			*resp = *metricsQueryResult
		}
	}).Once()

	result, err := testConf.client.GetMetricsCollection(ctx, queryID)
	fmt.Println("Metrics Query Result:", prettyPrintJSON(result), "Error:", err)
	if err != nil {
		t.Fatalf("Get Metrics Collection failed: %v", err)
	}
	assert.NotNil(t, result)

	fmt.Println("Get Metrics Collection Test - Successful")
}

func TestGetAllRealTimeMetricPaths(t *testing.T) {
	fmt.Println("Begin - GetAllRealTimeMetricPaths Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apitypes.MetricPaths)
		*resp = apitypes.MetricPaths{
			Entries: []apitypes.MetricEntries{
				{
					Cnt: apitypes.MetricContent{
						ID: 12345,
					},
				},
			},
		}
	}).Once()
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apitypes.MetricInstance)
		*resp = apitypes.MetricInstance{
			Content: apitypes.MetricInfo{
				IsRealtimeAvailable: true,
			},
		}
	}).Once()
	err := testConf.client.GetAllRealTimeMetricPaths(ctx)
	assert.Nil(t, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("get real time metrics failed")).Once()
	err = testConf.client.GetAllRealTimeMetricPaths(ctx)
	assert.Error(t, err)

	fmt.Println("GetAllRealTimeMetricPaths Test - Successful")
}

func TestCreateRealTimeMetricsQuery(t *testing.T) {
	fmt.Println("Begin - CreateRealTimeMetricsQuery Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apitypes.MetricQueryCreateResponse)
		*resp = apitypes.MetricQueryCreateResponse{
			Content: apitypes.MetricQueryResponseContent{
				Paths:    []string{"dummy"},
				Interval: 0,
			},
		}
	}).Once()
	_, err := testConf.client.CreateRealTimeMetricsQuery(ctx, []string{"dummy"}, 0)
	assert.Nil(t, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("create real time metrics failed")).Once()
	_, err = testConf.client.CreateRealTimeMetricsQuery(ctx, []string{"dummy"}, 0)
	assert.Error(t, err)

	fmt.Println("CreateRealTimeMetricsQuery Test - Successful")
}

func TestGetCapacity(t *testing.T) {
	fmt.Println("Begin - GetCapacity Test")
	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).ExpectedCalls = nil
	ctx := context.Background()

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apitypes.SystemCapacityMetricsQueryResult)
		*resp = apitypes.SystemCapacityMetricsQueryResult{
			Entries: []apitypes.SystemCapacityMetricsResultEntry{},
		}
	}).Once()
	_, err := testConf.client.GetCapacity(ctx)
	assert.Nil(t, err)

	testConf.client.(*UnityClientImpl).api.(*mocksapi.Client).On("DoWithHeaders", anyArgs...).Return(errors.New("get capacity failed")).Once()
	_, err = testConf.client.GetCapacity(ctx)
	assert.Error(t, err)

	fmt.Println("GetCapacity Test - Successful")
}
