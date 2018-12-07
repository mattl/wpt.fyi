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
	"runtime"
	"strconv"
	"time"

	"google.golang.org/api/option"

	"github.com/web-platform-tests/wpt.fyi/api/query"
	"github.com/web-platform-tests/wpt.fyi/api/query/cache/auth"
	"github.com/web-platform-tests/wpt.fyi/api/query/cache/backfill"
	"github.com/web-platform-tests/wpt.fyi/shared"

	"github.com/web-platform-tests/wpt.fyi/api/query/cache/index"
	"github.com/web-platform-tests/wpt.fyi/api/query/cache/monitor"

	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/datastore"
	log "github.com/sirupsen/logrus"
)

// PushID is a unique identifier for a request to push a test run to the search
// cache.
type PushID struct {
	Time  time.Time `json:"time"`
	RunID int64     `json:"run_id"`
}

var (
	port               = flag.Int("port", 8080, "Port to listen on")
	projectID          = flag.String("project_id", "", "Google Cloud Platform project ID, if different from ID detected from metadata service")
	gcpCredentialsFile = flag.String("gcp_credentials_file", "", "Path to Google Cloud Platform credentials file, if necessary")
	numShards          = flag.Int("num_shards", runtime.NumCPU(), "Number of shards for parallelizing query execution")
	monitorFrequency   = flag.Duration("monitor_frequency", time.Second*5, "Polling frequency for memory usage monitor")
	maxHeapBytes       = flag.Uint64("max_heap_bytes", uint64(1e+11), "Soft limit on heap-allocated bytes before evicting test runs from memory")

	idx   index.Index
	mon   monitor.Monitor
	authr auth.Authenticator
)

func livenessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Alive"))
}

func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ready"))
}

func pushRunHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Only PUT is supported", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	if !authr.Authenticate(ctx, r) {
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.URL.Query().Get("run_id"), 10, 0)
	if err != nil {
		http.Error(w, `Missing or invalid query param: "run_id"`, http.StatusBadRequest)
		return
	}

	t := time.Now().UTC()
	pushID := PushID{t, id}
	pushCtx := context.WithValue(context.Background(), shared.DefaultLoggerCtxKey(), log.WithFields(log.Fields{
		"search_cache_push_run_id": pushID,
	}))

	var client *datastore.Client
	if gcpCredentialsFile != nil && *gcpCredentialsFile != "" {
		client, err = datastore.NewClient(ctx, *projectID, option.WithCredentialsFile(*gcpCredentialsFile))
	} else {
		client, err = datastore.NewClient(ctx, *projectID)
	}
	if err != nil {
		http.Error(w, "Failed to connect to test run metadata service", http.StatusInternalServerError)
		return
	}

	var run shared.TestRun
	err = client.Get(ctx, &datastore.Key{Kind: "TestRun", ID: id}, &run)
	if err != nil {
		http.Error(w, fmt.Sprintf(`Unknown test run ID: %d`, id), http.StatusBadRequest)
		return
	}

	go func() {
		err := idx.IngestRun(run)
		if err != nil {
			shared.GetLogger(pushCtx).Errorf("Error pushing run: %v: %v", run, err)
		} else {
			shared.GetLogger(pushCtx).Infof("Cached new run: %v", run)
		}
	}()

	data, err := json.Marshal(pushID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(data)
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

	ids := make([]index.RunID, len(rq.RunIDs))
	for i := range rq.RunIDs {
		ids[i] = index.RunID(rq.RunIDs[i])
	}
	runs, err := idx.Runs(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	plan, err := idx.Bind(runs, rq.AbstractQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func init() {
	flag.Parse()
}

func main() {
	autoProjectID, err := metadata.ProjectID()
	if err != nil {
		log.Warningf("Failed to get project ID from metadata service; disabling authentication for pushing runs")
		authr = auth.NewNopAuthenticator()
	} else {
		if *projectID == "" {
			log.Infof(`Using project ID from metadata service: "%s"`, *projectID)
			*projectID = autoProjectID
		} else if *projectID != autoProjectID {
			log.Warningf(`Using project ID from flag: "%s" even though metadata service reports project ID of "%s"`, *projectID, autoProjectID)
		} else {
			log.Infof(`Using project ID: "%s"`, *projectID)
		}
		authr = auth.NewDatastoreAuthenticator(*projectID)
	}

	log.Infof("Serving index with %d shards", *numShards)
	// TODO: Use different field configurations for index, backfiller, monitor?
	logger := log.StandardLogger()

	idx, err = index.NewShardedWPTIndex(index.HTTPReportLoader{}, *numShards)
	if err != nil {
		log.Fatalf("Failed to instantiate index: %v", err)
	}

	mon, err = backfill.FillIndex(backfill.NewDatastoreRunFetcher(*projectID, gcpCredentialsFile, logger), logger, monitor.GoRuntime{}, *monitorFrequency, *maxHeapBytes, idx)
	if err != nil {
		log.Fatalf("Failed to initiate index backkfill: %v", err)
	}

	http.HandleFunc("/_ah/liveness_check", livenessCheckHandler)
	http.HandleFunc("/_ah/readiness_check", readinessCheckHandler)
	http.HandleFunc("/api/search/cache", searchHandler)
	http.HandleFunc("/api/search/cache/push_run", pushRunHandler)
	log.Infof("Listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
