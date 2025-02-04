/*
 * Copyright (c) 2021. Dell Inc., or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 */

package gounity

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
)

// GetAllRealTimeMetricPaths gets all the Unity Metric paths. Consider using for debugging
// or enumerating metrics. This will take a bit of time to complete.
// - /api/types/metric/instances?compact=true&filter=isRealtimeAvailable eq true
func (c *UnityClientImpl) GetAllRealTimeMetricPaths(ctx context.Context) error {
	log := util.GetRunIDLogger(ctx)
	filter := "isRealtimeAvailable eq true"

	query := fmt.Sprintf("%s&compact=true", url.QueryEscape(filter))
	queryURI := fmt.Sprintf(api.UnityInstancesFilter, api.UnityMetric, query)
	log.Info("GetAllRealTimeMetricPaths: ", queryURI)

	result := &types.MetricPaths{}

	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, queryURI, nil, result)
	if err != nil {
		return err
	}

	metricInstance := &types.MetricInstance{}
	for _, entry := range result.Entries {
		instanceURI := fmt.Sprintf(api.UnityAPIGetResourceURI, api.UnityMetric, strconv.Itoa(entry.Cnt.ID))
		c.executeWithRetryAuthenticate(ctx, http.MethodGet, instanceURI, nil, metricInstance)
		fmt.Printf("%d - %s - %s\n", metricInstance.Content.ID, metricInstance.Content.Path, metricInstance.Content.Description)
	}

	return nil
}

// GetMetricsCollection gets Unity MetricsCollection of the provided 'queryID'.
// - The MetricCollection should exist already or you can create one using CreateXXXMetricsQuery.
// - Example: GET /api/types/metricQueryResult/instances?filter=queryId eq 37
func (c *UnityClientImpl) GetMetricsCollection(ctx context.Context, queryID int) (*types.MetricQueryResult, error) {
	log := util.GetRunIDLogger(ctx)

	filter := fmt.Sprintf("queryId eq %d", queryID)
	queryURI := fmt.Sprintf(api.UnityInstancesFilter, api.UnityMetricQueryResult, url.QueryEscape(filter))
	log.Info("GetMetricsCollection: ", queryURI)

	metricsQueryResult := &types.MetricQueryResult{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, queryURI, nil, metricsQueryResult)
	if err != nil {
		return nil, err
	}

	return metricsQueryResult, nil
}

// CreateRealTimeMetricsQuery create an MetricRealTime Collection of the given metric paths and collection interval.
//   - The GetMetricsCollection interface can be called to retrieve results.
//   - Example: POST api/types/metricRealTimeQuery/instances
//     BODY:  {
//     "paths": ["sp.*.cpu.summary.busyTicks" ,"sp.*.cpu.summary.idleTicks"],
//     "interval": 5
//     }
func (c *UnityClientImpl) CreateRealTimeMetricsQuery(ctx context.Context, metricPaths []string, interval int) (*types.MetricQueryCreateResponse, error) {
	log := util.GetRunIDLogger(ctx)

	createURI := fmt.Sprintf(api.UnityAPIInstanceTypeResources, api.UnityMetricRealTimeQuery)
	log.Info("CreateRealTimeMetricQuery: ", createURI)

	metricQueryResponse := &types.MetricQueryCreateResponse{}
	metricQuery := types.MetricRealTimeQuery{
		Interval: interval,
		Paths:    metricPaths,
	}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodPost, createURI, metricQuery, metricQueryResponse)
	if err != nil {
		return nil, err
	}

	return metricQueryResponse, nil
}

// DeleteRealTimeMetricsQuery deletes the MetricRealTime Collection of the given queryID.
// - Example: DELETE /api/instances/metricRealTimeQuery/37
func (c *UnityClientImpl) DeleteRealTimeMetricsQuery(ctx context.Context, queryID int) error {
	log := util.GetRunIDLogger(ctx)
	deleteURI := fmt.Sprintf(api.UnityAPIGetResourceURI, api.UnityMetricRealTimeQuery, strconv.Itoa(queryID))
	log.Info("DeleteRealTimeMetricsQuery:", deleteURI)

	err := c.executeWithRetryAuthenticate(ctx, http.MethodDelete, deleteURI, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetCapacity gets Unity capacity metrics at the system level.
// - Example: GET /api/types/systemCapacity/instances?fields=id,sizeFree,sizeTotal,sizeUsed,sizePreallocated,sizeSubscribed,totalLogicalSize
func (c *UnityClientImpl) GetCapacity(ctx context.Context) (*types.SystemCapacityMetricsQueryResult, error) {
	log := util.GetRunIDLogger(ctx)

	queryURI := fmt.Sprintf(api.UnityAPIInstanceTypeResourcesWithFields, api.UnitySystemCapacity, SystemCapacityFields)
	log.Info("GetSystemCapacityMetrics: ", queryURI)

	systemCapacityMetricsQueryResult := &types.SystemCapacityMetricsQueryResult{}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, queryURI, nil, systemCapacityMetricsQueryResult)
	if err != nil {
		return nil, err
	}

	return systemCapacityMetricsQueryResult, nil
}
