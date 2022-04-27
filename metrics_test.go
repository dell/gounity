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
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestMetrics(t *testing.T) {
	ctx = context.Background()

	getVolumeMetrics(t)
}

func getVolumeMetrics(t *testing.T) {
	debugOn, _ := strconv.ParseBool(os.Getenv("GOUNITY_SHOWHTTP"))
	if debugOn {
		level, _ := log.ParseLevel("debug")
		log.SetLevel(level)
	}

	queryID := -1

	fmt.Println("Begin - Realtime Volume Metrics Query")
	defer func() {
		fmt.Println("End - Realtime Volume Metrics Query")
		// Clean up the query if it was created
		if queryID != -1 {
			err := testConf.metricsAPI.DeleteRealTimeMetricsQuery(ctx, queryID)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	var err error

	err = testConf.metricsAPI.GetAllRealTimeMetricPaths(ctx)
	if err != nil {
		t.Fatalf("Get all real time Metric Paths failed: %v", err)
	}

	paths := []string{
		"sp.*.storage.lun.*.reads",
		"sp.*.storage.lun.*.writes",
		"sp.*.cpu.summary.busyTicks",
		"sp.*.cpu.summary.idleTicks",
	}

	interval := 5 // seconds
	query, err := testConf.metricsAPI.CreateRealTimeMetricsQuery(ctx, paths, interval)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Example result:
	//	============== 1 ==============
	//	Timestamp: 2021-04-08T13:42:50.000Z
	//  QueryID:   71
	//  Path:      sp.*.storage.lun.*.reads [spa = map[sv_108:0 sv_18:0 sv_19:0 sv_22:0 sv_23:0 sv_24:0 sv_25:0 sv_26:0 sv_27:0 sv_28:0 sv_29:0 sv_42:0 sv_43:0]]
	//	Path:      sp.*.storage.lun.*.writes [spa = map[sv_108:0 sv_18:0 sv_19:0 sv_22:0 sv_23:0 sv_24:0 sv_25:0 sv_26:0 sv_27:0 sv_28:0 sv_29:0 sv_42:0 sv_43:0]]
	//	Path:      sp.*.cpu.summary.busyTicks [spa = 243675336]
	//	Path:      sp.*.cpu.summary.idleTicks [spa = 615488915]
	//	================================
	queryID = query.Content.ID
	fmt.Printf("Created MetricsQuery %d. Waiting %d seconds before trying queries\n", queryID, interval)
	for i := 1; i <= 2; i++ {
		time.Sleep(time.Duration(interval) * time.Second)
		timeMetrics, err2 := testConf.metricsAPI.GetMetricsCollection(ctx, queryID)
		if err2 != nil {
			t.Fatal(err2)
			return
		}
		fmt.Printf("============== %d ==============\n", i)
		doOnce := true
		for _, entry := range timeMetrics.Entries {
			if doOnce {
				fmt.Printf("Timestamp: %s\n", entry.Content.Timestamp)
				fmt.Printf("QueryID:   %d\n", entry.Content.QueryID)
				doOnce = false
			}
			keyValues := make([]string, 0)
			for k, v := range entry.Content.Values {
				keyValues = append(keyValues, fmt.Sprintf("%s = %s", k, v))
			}
			fmt.Printf("Path:      %s [%s]\n", entry.Content.Path, strings.Join(keyValues, ","))
		}
		fmt.Println("================================")
	}
}
