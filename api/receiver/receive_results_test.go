// +build small

// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package receiver

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/web-platform-tests/wpt.fyi/shared/sharedtest"
	"google.golang.org/appengine/taskqueue"
)

// regexMatcher is a gomock.Matcher that verifies whether a string argument
// matches the predefined regular expression.
//
// This is used to match arguments containing random strings (e.g. UUID).
type regexMatcher struct {
	regex *regexp.Regexp
}

func (r *regexMatcher) Matches(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return false
	}
	return r.regex.MatchString(s)
}

func (r *regexMatcher) String() string {
	return "matches " + r.regex.String()
}

func matchRegex(r string) *regexMatcher {
	return &regexMatcher{regex: regexp.MustCompile(r)}
}

// An empty (default) extraParams
var emptyParams = map[string]string{
	"browser_name":    "",
	"labels":          "",
	"revision":        "",
	"browser_version": "",
	"os_name":         "",
	"os_version":      "",
	"callback_url":    "",
}

func TestHandleResultsUpload_not_admin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	req := httptest.NewRequest("POST", "/api/results/upload", nil)
	resp := httptest.NewRecorder()
	mockAE := NewMockAppEngineAPI(mockCtrl)
	mockAE.EXPECT().IsAdmin().Return(false)

	HandleResultsUpload(mockAE, resp, req)
	assert.Equal(t, resp.Code, http.StatusUnauthorized)
}

func TestHandleResultsUpload_http_basic_auth_invalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	req := httptest.NewRequest("POST", "/api/results/upload", nil)
	req.SetBasicAuth("not_a_user", "123")
	resp := httptest.NewRecorder()
	mockAE := NewMockAppEngineAPI(mockCtrl)
	gomock.InOrder(
		mockAE.EXPECT().IsAdmin().Return(false),
		mockAE.EXPECT().authenticateUploader("not_a_user", "123").Return(false),
	)

	HandleResultsUpload(mockAE, resp, req)
	assert.Equal(t, resp.Code, http.StatusUnauthorized)
}

func TestHandleResultsUpload_extra_params(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	payload := url.Values{
		"result_url":    {"http://wpt.fyi/test.json.gz"},
		"browser_name":  {"firefox"},
		"labels":        {"stable"},
		"invalid_param": {"should be ignored"},
	}
	req := httptest.NewRequest("POST", "/api/results/upload", strings.NewReader(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("blade-runner", "123")
	resp := httptest.NewRecorder()

	f := &os.File{}
	extraParams := map[string]string{
		"browser_name":    "firefox",
		"labels":          "stable",
		"revision":        "",
		"browser_version": "",
		"os_name":         "",
		"os_version":      "",
		"callback_url":    "",
	}
	task := &taskqueue.Task{Name: "task"}
	mockAE := NewMockAppEngineAPI(mockCtrl)
	mockAE.EXPECT().Context().Return(sharedtest.NewTestContext()).AnyTimes()
	gomock.InOrder(
		mockAE.EXPECT().IsAdmin().Return(false),
		mockAE.EXPECT().authenticateUploader("blade-runner", "123").Return(true),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(f, nil),
		mockAE.EXPECT().uploadToGCS(matchRegex(`^/wptd-results-buffer/blade-runner/.*\.json$`), f, true).Return(nil),
		mockAE.EXPECT().scheduleResultsTask("blade-runner", gomock.Any(), "single", extraParams).Return(task, nil),
	)

	HandleResultsUpload(mockAE, resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
}

func TestHandleResultsUpload_url(t *testing.T) {
	var urls []string
	for i := 1; i <= 2; i++ {
		urls = append(urls, fmt.Sprintf("http://wpt.fyi/wpt_report_%d.json.gz", i))
		t.Run(fmt.Sprintf("%d url(s)", len(urls)), func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var payloadType string
			if len(urls) == 1 {
				payloadType = "single"
			} else {
				payloadType = "multiple"
			}
			payload := url.Values{"result_url": urls}
			req := httptest.NewRequest("POST", "/api/results/upload", strings.NewReader(payload.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth("blade-runner", "123")
			resp := httptest.NewRecorder()

			files := make([]os.File, len(urls))
			task := &taskqueue.Task{Name: "task"}
			mockAE := NewMockAppEngineAPI(mockCtrl)
			mockAE.EXPECT().Context().Return(sharedtest.NewTestContext()).AnyTimes()
			gomock.InOrder(
				mockAE.EXPECT().IsAdmin().Return(false),
				mockAE.EXPECT().authenticateUploader("blade-runner", "123").Return(true),
				mockAE.EXPECT().scheduleResultsTask("blade-runner", gomock.Any(), payloadType, emptyParams).Return(task, nil),
			)
			for i, url := range urls {
				gomock.InOrder(
					mockAE.EXPECT().fetchWithTimeout(url, DownloadTimeout).Return(&files[i], nil),
					mockAE.EXPECT().uploadToGCS(matchRegex(`^/wptd-results-buffer/blade-runner/.*\.json$`), &files[i], true).Return(nil),
				)
			}

			HandleResultsUpload(mockAE, resp, req)
			assert.Equal(t, resp.Code, http.StatusOK)
		})
	}
}

func TestHandleResultsUpload_file(t *testing.T) {
	var filenames []string
	for i := 1; i <= 2; i++ {
		filenames = append(filenames, fmt.Sprintf("wpt_report_%d.json.gz", i))
		t.Run(fmt.Sprintf("%d file(s)", len(filenames)), func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var payloadType string
			if len(filenames) == 1 {
				payloadType = "single"
			} else {
				payloadType = "multiple"
			}
			buffer := new(bytes.Buffer)
			writer := multipart.NewWriter(buffer)
			for _, filename := range filenames {
				writer.CreateFormFile("result_file", filename)
			}
			writer.Close()
			req := httptest.NewRequest("POST", "/api/results/upload", buffer)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.SetBasicAuth("blade-runner", "123")
			resp := httptest.NewRecorder()

			task := &taskqueue.Task{Name: "task"}
			mockAE := NewMockAppEngineAPI(mockCtrl)
			mockAE.EXPECT().Context().Return(sharedtest.NewTestContext()).AnyTimes()
			gomock.InOrder(
				mockAE.EXPECT().IsAdmin().Return(false),
				mockAE.EXPECT().authenticateUploader("blade-runner", "123").Return(true),
				mockAE.EXPECT().scheduleResultsTask("blade-runner", gomock.Any(), payloadType, emptyParams).Return(task, nil),
			)
			for range filenames {
				mockAE.EXPECT().uploadToGCS(matchRegex(`^/wptd-results-buffer/blade-runner/.*\.json$`), gomock.Any(), true).Return(nil)
			}

			HandleResultsUpload(mockAE, resp, req)
			assert.Equal(t, resp.Code, http.StatusOK)
		})
	}
}

func TestHandleResultsUpload_retry_fetching(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	payload := url.Values{
		"result_url": {"http://wpt.fyi/test.json.gz"},
	}
	req := httptest.NewRequest("POST", "/api/results/upload", strings.NewReader(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("blade-runner", "123")
	resp := httptest.NewRecorder()

	f := &os.File{}
	errTimeout := fmt.Errorf("server timed out")
	task := &taskqueue.Task{Name: "task"}
	mockAE := NewMockAppEngineAPI(mockCtrl)
	mockAE.EXPECT().Context().Return(sharedtest.NewTestContext()).AnyTimes()
	gomock.InOrder(
		mockAE.EXPECT().IsAdmin().Return(false),
		mockAE.EXPECT().authenticateUploader("blade-runner", "123").Return(true),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(nil, errTimeout),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(nil, errTimeout),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(f, nil),
		mockAE.EXPECT().uploadToGCS(matchRegex(`^/wptd-results-buffer/blade-runner/.*\.json$`), f, true).Return(nil),
		mockAE.EXPECT().scheduleResultsTask("blade-runner", gomock.Any(), "single", emptyParams).Return(task, nil),
	)

	HandleResultsUpload(mockAE, resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
}

func TestHandleResultsUpload_fail_fetching(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	payload := url.Values{
		"result_url": {"http://wpt.fyi/test.json.gz"},
	}
	req := httptest.NewRequest("POST", "/api/results/upload", strings.NewReader(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("blade-runner", "123")
	resp := httptest.NewRecorder()

	errTimeout := fmt.Errorf("server timed out")
	mockAE := NewMockAppEngineAPI(mockCtrl)
	mockAE.EXPECT().Context().Return(sharedtest.NewTestContext()).AnyTimes()
	gomock.InOrder(
		mockAE.EXPECT().IsAdmin().Return(false),
		mockAE.EXPECT().authenticateUploader("blade-runner", "123").Return(true),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(nil, errTimeout),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(nil, errTimeout),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(nil, errTimeout),
	)

	HandleResultsUpload(mockAE, resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}

func TestHandleResultsUpload_fail_uploading(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	payload := url.Values{
		"result_url": {"http://wpt.fyi/test.json.gz"},
	}
	req := httptest.NewRequest("POST", "/api/results/upload", strings.NewReader(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("blade-runner", "123")
	resp := httptest.NewRecorder()

	f := &os.File{}
	errGCS := fmt.Errorf("failed to upload to GCS")
	mockAE := NewMockAppEngineAPI(mockCtrl)
	mockAE.EXPECT().Context().Return(sharedtest.NewTestContext()).AnyTimes()
	gomock.InOrder(
		mockAE.EXPECT().IsAdmin().Return(false),
		mockAE.EXPECT().authenticateUploader("blade-runner", "123").Return(true),
		mockAE.EXPECT().fetchWithTimeout("http://wpt.fyi/test.json.gz", DownloadTimeout).Return(f, nil),
		mockAE.EXPECT().uploadToGCS(matchRegex(`^/wptd-results-buffer/blade-runner/.*\.json$`), f, true).Return(errGCS),
	)

	HandleResultsUpload(mockAE, resp, req)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}
