// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package azure

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"

	mapset "github.com/deckarep/golang-set"

	"github.com/google/go-github/github"
	uc "github.com/web-platform-tests/wpt.fyi/api/receiver/client"
	"github.com/web-platform-tests/wpt.fyi/shared"
)

// handleCheckRunEvent processes an Azure Pipelines check run "completed" event.
func handleCheckRunEvent(azureAPI API, aeAPI shared.AppEngineAPI, event *github.CheckRunEvent) (bool, error) {
	log := shared.GetLogger(aeAPI.Context())
	status := event.GetCheckRun().GetStatus()
	if status != "completed" {
		log.Infof("Ignoring non-completed status %s", status)
		return false, nil
	}
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	sender := event.GetSender().GetLogin()
	detailsURL := event.GetCheckRun().GetDetailsURL()
	sha := event.GetCheckRun().GetHeadSHA()
	buildID := extractAzureBuildID(detailsURL)
	if buildID == 0 {
		log.Errorf("Failed to extract build ID from details_url \"%s\"", detailsURL)
		return false, nil
	}
	return processAzureBuild(aeAPI, azureAPI, sha, owner, repo, sender, buildID)
}

func processAzureBuild(aeAPI shared.AppEngineAPI, azureAPI API, sha, owner, repo, sender string, buildID int64) (bool, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/build/artifacts/get?view=azure-devops-rest-4.1
	artifactsURL := azureAPI.GetAzureArtifactsURL(owner, repo, buildID)

	log := shared.GetLogger(aeAPI.Context())
	log.Infof("Fetching %s", artifactsURL)

	client := aeAPI.GetHTTPClient()
	resp, err := client.Get(artifactsURL)
	if err != nil {
		log.Errorf("Failed to fetch artifacts for %s/%s build %v", owner, repo, buildID)
		return false, err
	}

	var artifacts BuildArtifacts
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("Failed to read response body")
		return false, err
	} else if err = json.Unmarshal(body, &artifacts); err != nil {
		log.Errorf("Failed to unmarshal JSON")
		return false, err
	}

	uploadedAny := false
	errors := make(chan (error), artifacts.Count)
	for _, artifact := range artifacts.Value {
		if err != nil {
			log.Errorf("Failed to extract report URL: %s", err.Error())
			continue
		}
		log.Infof("Uploading %s for %s/%s build %v...", artifact.Name, owner, repo, buildID)

		labels := mapset.NewSet()
		if sender != "" {
			labels.Add(shared.GetUserLabel(sender))
		}
		switch artifact.Name {
		case "results":
			labels.Add(shared.MasterLabel)
		case "affected-tests":
			labels.Add(shared.PRHeadLabel)
		case "affected-tests-without-changes":
			labels.Add(shared.PRBaseLabel)
		}

		uploader, err := aeAPI.GetUploader("azure")
		if err != nil {
			log.Errorf("Failed to load azure uploader")
			return false, err
		}

		uploadClient := uc.NewClient(aeAPI)
		err = uploadClient.CreateRun(
			sha,
			uploader.Username,
			uploader.Password,
			[]string{artifact.Resource.DownloadURL},
			shared.ToStringSlice(labels))
		if err != nil {
			log.Errorf("Failed to create run: %s", err.Error())
			errors <- err
		} else {
			uploadedAny = true
		}
	}
	close(errors)
	for err := range errors {
		return uploadedAny, err
	}
	return uploadedAny, nil
}

func extractAzureBuildID(detailsURL string) int64 {
	parsedURL, err := url.Parse(detailsURL)
	if err != nil {
		return 0
	}
	id := parsedURL.Query().Get("buildId")
	parsedID, _ := strconv.ParseInt(id, 0, 0)
	return parsedID
}
