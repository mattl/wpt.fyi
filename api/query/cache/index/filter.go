// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package index

import (
	"errors"
	"strings"

	"github.com/web-platform-tests/wpt.fyi/api/query"
)

type filter interface {
	Filter(TestID) bool
	Tests() Tests
}

type index struct {
	tests      Tests
	runResults map[RunID]RunResults
}

type TestNamePattern struct {
	index
	q query.TestNamePattern
}

type RunTestStatusConstraint struct {
	index
	q query.RunTestStatusConstraint
}

type And struct {
	index
	args []filter
}

type Or struct {
	index
	args []filter
}

type Not struct {
	index
	arg filter
}

type ShardedFilter []filter

var errUnknownConcreteQuery = errors.New("Unknown ConcreteQuery type")

func (i index) Tests() Tests { return i.tests }

func (tnp TestNamePattern) Filter(t TestID) bool {
	name, _, err := tnp.tests.GetName(t)
	if err != nil {
		return false
	}
	return strings.Contains(name, tnp.q.Pattern)
}

func (rtsc RunTestStatusConstraint) Filter(t TestID) bool {
	return rtsc.runResults[RunID(rtsc.q.Run)].GetResult(t) == ResultID(rtsc.q.Status)
}

func (a And) Filter(t TestID) bool {
	args := a.args
	for _, arg := range args {
		if !arg.Filter(t) {
			return false
		}
	}
	return true
}

func (o Or) Filter(t TestID) bool {
	args := o.args
	for _, arg := range args {
		if arg.Filter(t) {
			return true
		}
	}
	return false
}

func (n Not) Filter(t TestID) bool {
	return !n.arg.Filter(t)
}

func NewFilter(idx index, q query.ConcreteQuery) (filter, error) {
	switch v := q.(type) {
	case query.TestNamePattern:
		return TestNamePattern{idx, v}, nil
	case query.RunTestStatusConstraint:
		return RunTestStatusConstraint{idx, v}, nil
	case query.And:
		fs, err := filters(idx, v.Args)
		if err != nil {
			return nil, err
		}
		return And{idx, fs}, nil
	case query.Or:
		fs, err := filters(idx, v.Args)
		if err != nil {
			return nil, err
		}
		return Or{idx, fs}, nil
	case query.Not:
		f, err := NewFilter(idx, v.Arg)
		if err != nil {
			return nil, err
		}
		return Not{idx, f}, nil
	default:
		return nil, errUnknownConcreteQuery
	}
}

func (fs ShardedFilter) Execute() (interface{}, error) {
	res := make(chan []TestID, len(fs))
	for _, f := range fs {
		go func(f filter) {
			ts := make([]TestID, 0)
			f.Tests().Range(func(t TestID) bool {
				if f.Filter(t) {
					ts = append(ts, t)
				}
				return true
			})
			res <- ts
		}(f)
	}

	ret := make([]TestID, 0)
	for ts := range res {
		ret = append(ret, ts...)
	}
	return ret, nil
}

func filters(idx index, qs []query.ConcreteQuery) ([]filter, error) {
	fs := make([]filter, len(qs))
	var err error
	for i := range qs {
		fs[i], err = NewFilter(idx, qs[i])
		if err != nil {
			return nil, err
		}
	}
	return fs, nil
}
