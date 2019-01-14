// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/web-platform-tests/wpt.fyi/api/query"
	"github.com/web-platform-tests/wpt.fyi/api/query/cache/backfill"
	"github.com/web-platform-tests/wpt.fyi/api/query/cache/poll"
	"github.com/web-platform-tests/wpt.fyi/shared"
	"google.golang.org/api/option"

	"github.com/web-platform-tests/wpt.fyi/api/query/cache/index"
	"github.com/web-platform-tests/wpt.fyi/api/query/cache/monitor"
	cq "github.com/web-platform-tests/wpt.fyi/api/query/cache/query"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/datastore"
	log "github.com/sirupsen/logrus"
)

var (
	port                   = flag.Int("port", 8080, "Port to listen on")
	projectID              = flag.String("project_id", "", "Google Cloud Platform project ID, if different from ID detected from metadata service")
	gcpCredentialsFile     = flag.String("gcp_credentials_file", "", "Path to Google Cloud Platform credentials file, if necessary")
	numShards              = flag.Int("num_shards", runtime.NumCPU(), "Number of shards for parallelizing query execution")
	monitorInterval        = flag.Duration("monitor_interval", time.Second*5, "Polling interval for memory usage monitor")
	monitorMaxIngestedRuns = flag.Uint("monitor_max_ingested_runs", uint(10), "Maximum number of runs that can be ingested before memory monitor must run")
	maxHeapBytes           = flag.Uint64("max_heap_bytes", uint64(2e+11), "Soft limit on heap-allocated bytes before evicting test runs from memory")
	evictRunsPercent       = flag.Float64("evict_runs_percent", 0.1, "Decimal percentage indicating what fraction of runs to evict when soft memory limit is reached")
	updateInterval         = flag.Duration("update_interval", time.Second*10, "Update interval for polling for new runs")
	updateMaxRuns          = flag.Int("update_max_runs", 10, "The maximum number of latest runs to lookup in attempts to update indexes via polling")
	maxRunsPerRequest      = flag.Int("max_runs_per_request", 16, "Maximum number of runs that may be queried per request")

	maxRunsPerRequestMsg string

	idx index.Index
	mon monitor.Monitor
)

func livenessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Alive"))
}

func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ready"))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid HTTP method", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	}
	err = r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to finish reading request body", http.StatusInternalServerError)
	}

	var rq query.RunQuery
	err = json.Unmarshal(data, &rq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(rq.RunIDs) > *maxRunsPerRequest {
		http.Error(w, maxRunsPerRequestMsg, http.StatusBadRequest)
		return
	}

	ids := make([]index.RunID, len(rq.RunIDs))
	for i := range rq.RunIDs {
		ids[i] = index.RunID(rq.RunIDs[i])
	}

	// Ensure runs are loaded before executing query. This is best effort: It is
	// possible, though unlikely, that a run may exist in the cache at this point
	// and be evicted before binding the query to a query execution plan. In such
	// a case, `idx.Bind()` below will return an error.
	runs := make([]shared.TestRun, len(ids))
	for i, id := range ids {
		run, err := idx.Run(id)
		// If getting run metadata fails, attempt write-on-read for this run.
		if err != nil {
			runPtr, err := getTestRun(int64(id))
			if err != nil {
				http.Error(w, fmt.Sprintf("Unknown test run ID: %d", id), http.StatusBadRequest)
				return
			}
			err = idx.IngestRun(*runPtr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to load test run data for run %d: %v", id, err), http.StatusInternalServerError)
				return
			}
			run = *runPtr
		}
		runs[i] = run
	}

	q := cq.PrepareUserQuery(rq.RunIDs, rq.AbstractQuery.BindToRuns(runs))
	plan, err := idx.Bind(runs, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := plan.Execute(runs)
	res, ok := results.([]query.SearchResult)
	if !ok {
		http.Error(w, "Search index returned bad results", http.StatusInternalServerError)
		return
	}

	data, err = json.Marshal(query.SearchResponse{
		Runs:    runs,
		Results: res,
	})
	if err != nil {
		http.Error(w, "Failed to marshal results to JSON", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func getTestRun(id int64) (*shared.TestRun, error) {
	ctx := context.Background()
	var client *datastore.Client
	var err error
	if gcpCredentialsFile != nil && *gcpCredentialsFile != "" {
		client, err = datastore.NewClient(ctx, *projectID, option.WithCredentialsFile(*gcpCredentialsFile))
	} else {
		client, err = datastore.NewClient(ctx, *projectID)
	}
	if err != nil {
		return nil, err
	}
	d := shared.NewCloudDatastore(ctx, client)
	testRun := new(shared.TestRun)
	err = d.Get(d.NewKey("TestRun", id), testRun)
	if err != nil {
		return nil, err
	}

	testRun.ID = id
	return testRun, nil
}

func init() {
	flag.Parse()

	maxRunsPerRequestMsg = fmt.Sprintf("Too many runs specified; maximum is %d.", *maxRunsPerRequest)
}

func main() {
	autoProjectID, err := metadata.ProjectID()
	if err != nil {
		log.Warningf("Failed to get project ID from metadata service")
	} else {
		if *projectID == "" {
			log.Infof(`Using project ID from metadata service: "%s"`, *projectID)
			*projectID = autoProjectID
		} else if *projectID != autoProjectID {
			log.Warningf(`Using project ID from flag: "%s" even though metadata service reports project ID of "%s"`, *projectID, autoProjectID)
		} else {
			log.Infof(`Using project ID: "%s"`, *projectID)
		}
	}

	log.Infof("Serving index with %d shards", *numShards)
	// TODO: Use different field configurations for index, backfiller, monitor?
	logger := log.StandardLogger()

	idx, err = index.NewShardedWPTIndex(index.HTTPReportLoader{}, *numShards)
	if err != nil {
		log.Fatalf("Failed to instantiate index: %v", err)
	}

	fetcher := backfill.NewDatastoreRunFetcher(*projectID, gcpCredentialsFile, logger)
	mon, err = backfill.FillIndex(fetcher, logger, monitor.GoRuntime{}, *monitorInterval, *monitorMaxIngestedRuns, *maxHeapBytes, *evictRunsPercent, idx)
	if err != nil {
		log.Fatalf("Failed to initiate index backkfill: %v", err)
	}

	// Index, backfiller, monitor now in place. Start polling to load runs added
	// after backfilling was started.
	go poll.KeepRunsUpdated(fetcher, logger, *updateInterval, *updateMaxRuns, idx)

	http.HandleFunc("/_ah/liveness_check", livenessCheckHandler)
	http.HandleFunc("/_ah/readiness_check", readinessCheckHandler)
	http.HandleFunc("/api/search/cache", searchHandler)
	log.Infof("Listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
