/**
 * Copyright 2018 The WPT Dashboard Project. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

import '../node_modules/@polymer/paper-toggle-button/paper-toggle-button.js';
import '../node_modules/@polymer/polymer/lib/elements/dom-if.js';
import { html } from '../node_modules/@polymer/polymer/lib/utils/html-tag.js';
import { PolymerElement } from '../node_modules/@polymer/polymer/polymer-element.js';
import { TestRunsUIQuery } from './test-runs-query.js';
import { TestRunsQueryLoader } from './test-runs.js';
import './wpt-colors.js';

class TestFileResults extends TestRunsUIQuery(
  TestRunsQueryLoader(PolymerElement, TestRunsUIQuery.Computer)) {
  static get template() {
    return html`
    <style include="wpt-colors">
      :host {
        display: block;
        font-size: 16px;
      }
      h1 {
        font-size: 1.5em;
      }
      .right {
        display: flex;
        justify-content: flex-end;
      }
      .right .pad {
        padding: 8px;
      }
      paper-toggle-button {
        --paper-toggle-button-checked-bar-color:  var(--paper-blue-500);
        --paper-toggle-button-checked-button-color:  var(--paper-blue-700);
        --paper-toggle-button-checked-ink-color: var(--paper-blue-300);
      }
    </style>

    <div class="right">
      <label class="pad">Expand</label>
      <paper-toggle-button class="pad" checked="{{isVerbose}}">
      </paper-toggle-button>
    </div>

    <template is="dom-if" if="{{!isVerbose}}">
      <test-file-results-table-terse test-runs="[[testRuns]]" results-table="[[resultsTable]]">
      </test-file-results-table-terse>
    </template>

    <template is="dom-if" if="{{isVerbose}}">
      <test-file-results-table-verbose test-runs="[[testRuns]]" results-table="[[resultsTable]]">
      </test-file-results-table-verbose>
    </template>
`;
  }

  static get is() {
    return 'test-file-results';
  }

  static get properties() {
    return {
      resultsTable: {
        type: Array,
        value: [],
      },
      isVerbose: {
        type: Boolean,
        value: false,
      },
    };
  }

  async connectedCallback() {
    await super.connectedCallback();
    console.assert(this.path);
    console.assert(this.path[0] === '/');
  }

  static get observers() {
    return ['fetchTestFile(path, testRuns)'];
  }

  async fetchTestFile(path, testRuns) {
    this.resultsTable = []; // Clear any existing rows.
    if (!path || !testRuns) {
      return;
    }
    const resultsPerTestRun = await Promise.all(
      testRuns.map(tr => this.loadResultFile(tr)));

    // resultsTable[0].name set after discovering subtests.
    let resultsTable = [{
      results: resultsPerTestRun.map(data => {
        return {
          status: data && data.status,
          message: data && data.message,
        };
      }),
    }];

    // Setup test name order according to when they appear in run results.
    let allNames = [];
    for (const runResults of resultsPerTestRun) {
      if (runResults && runResults.subtests) {
        this.mergeNamesInto(runResults.subtests.map(s => s.name), allNames);
      }
    }

    // Copy results into resultsTable.
    for (const name of allNames) {
      let results = [];
      for (const runResults of resultsPerTestRun) {
        const result = runResults && runResults.subtests &&
          runResults.subtests.find(sub => sub.name === name);
        results.push(result ? {
          status: result.status,
          message: result.message,
        } : {status: null, message: null});
      }
      resultsTable.push({
        name,
        results,
      });
    }

    // Set name for test-level status entry after subtests discovered.
    // Parameter is number of subtests.
    resultsTable[0].name = this.statusName(resultsTable.length - 1);

    this.resultsTable = resultsTable;
  }

  async loadResultFile(testRun) {
    const url = this.resultsURL(testRun, this.path);
    const response = await window.fetch(url);
    if (!response.ok) {
      return null;
    }
    return response.json();
  }

  mergeNamesInto(names, allNames) {
    if (!allNames.length) {
      allNames.splice(0, 0, ...names);
      return;
    }
    let lastOffset = 0;
    let lastMatch = 0;
    names.forEach((name, i) => {
      // Optimization for "next item matches too".
      let offset;
      if (i === lastMatch + 1 && allNames[lastOffset + 1] === name) {
        offset = lastOffset + 1;
      } else {
        offset = allNames.findIndex(n => n === name);
      }
      if (offset >= 0) {
        lastOffset = offset;
        lastMatch = i;
      } else {
        allNames.splice(lastOffset + i - lastMatch, 0, name);
      }
    });
  }

  resultsURL(testRun, path) {
    path = this.encodeTestPath(path);
    // This is relying on the assumption that result files end with '-summary.json.gz'.
    const resultsBase = testRun.results_url.slice(0, testRun.results_url.lastIndexOf('-summary.json.gz'));
    return `${resultsBase}${path}`;
  }

  statusName(numSubtests) {
    return numSubtests > 0 ? 'Harness status' : 'Test status';
  }
}

window.customElements.define(TestFileResults.is, TestFileResults);

export { TestFileResults };
