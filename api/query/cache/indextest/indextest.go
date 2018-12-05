// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package indextest

import (
	"sync/atomic"

	"github.com/web-platform-tests/wpt.fyi/api/query/cache/index"
	"github.com/web-platform-tests/wpt.fyi/shared"
)

// CountingIndex counts the number of runs currently stored in the index.
type CountingIndex struct {
	index.ProxyIndex

	count int32
}

// IngestRun delegates and increments count iff no error is returned from the
// delegate.
func (i *CountingIndex) IngestRun(r shared.TestRun) error {
	err := i.ProxyIndex.IngestRun(r)
	if err != nil {
		return err
	}

	atomic.AddInt32(&i.count, 1)
	return nil
}

// EvictAnyRun delegates and decrements count iff no error is returned from the
// delegate.
func (i *CountingIndex) EvictAnyRun() error {
	err := i.ProxyIndex.EvictAnyRun()
	if err != nil {
		return err
	}

	atomic.AddInt32(&i.count, -1)
	return nil
}

// Count returns the current count of runs in the index.
func (i *CountingIndex) Count() int32 {
	return i.count
}

// NewCountingIndex constructs a new CountIndex bound to the given delegate.
func NewCountingIndex(delegate index.Index) *CountingIndex {
	return &CountingIndex{index.NewProxyIndex(delegate), 0}
}
