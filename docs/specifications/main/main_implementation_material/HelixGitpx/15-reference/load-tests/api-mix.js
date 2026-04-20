// load-tests/api-mix.js
// k6 load test simulating a realistic API mix for HelixGitpx.
// Run:
//   k6 run --env TARGET=https://api.staging.helixgitpx.example.com \
//          --env PAT=hpxat_xxxxxx tests/load/api-mix.js
//
// Profiles:
//   smoke  — 1 VU, 30s — quick sanity
//   base   — 100 VUs, 5m — baseline
//   peak   — ramp to 2000 VUs over 10m, hold 20m
//   soak   — 500 VUs, 8h — leak / memory pressure
//
// Asserts SLOs: p99 read < 300ms, p99 write < 500ms, err-rate < 0.1%.
import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  scenarios: {
    reads: {
      executor: 'ramping-vus',
      exec: 'readFlow',
      startVUs: 0,
      stages: profile('reads'),
      gracefulRampDown: '1m',
    },
    writes: {
      executor: 'ramping-vus',
      exec: 'writeFlow',
      startVUs: 0,
      stages: profile('writes'),
      gracefulRampDown: '1m',
    },
    subscribe: {
      executor: 'constant-vus',
      exec: 'subscribeFlow',
      vus: subs(),
      duration: duration(),
    },
  },
  thresholds: {
    http_req_duration: [
      { threshold: 'p(95)<250', abortOnFail: false },
      { threshold: 'p(99)<500', abortOnFail: false },
    ],
    'http_req_duration{op:read}': ['p(99)<300'],
    'http_req_duration{op:write}': ['p(99)<500'],
    http_req_failed: ['rate<0.001'],
    errors: ['count<10'],
  },
  summaryTrendStats: ['avg', 'med', 'p(90)', 'p(95)', 'p(99)', 'p(99.9)'],
  noConnectionReuse: false,
  discardResponseBodies: false,
};

function profile(kind) {
  const p = (__ENV.PROFILE || 'base');
  const mul = kind === 'reads' ? 3 : 1; // more reads than writes
  if (p === 'smoke')
    return [{ duration: '30s', target: 1 }];
  if (p === 'peak')
    return [
      { duration: '3m', target: 100 * mul },
      { duration: '7m', target: 2000 * mul },
      { duration: '20m', target: 2000 * mul },
      { duration: '2m', target: 0 },
    ];
  if (p === 'soak')
    return [
      { duration: '2m', target: 500 * mul },
      { duration: '8h', target: 500 * mul },
    ];
  return [
    { duration: '2m', target: 100 * mul },
    { duration: '5m', target: 100 * mul },
    { duration: '1m', target: 0 },
  ];
}
function subs() { return parseInt(__ENV.SUBS || '50'); }
function duration() { return __ENV.SUB_DURATION || '5m'; }

const errors = new Counter('errors');
const readLatency = new Trend('read_latency_ms', true);
const writeLatency = new Trend('write_latency_ms', true);
const createdRepos = new Counter('repos_created');
const createdIssues = new Counter('issues_created');

const TARGET = __ENV.TARGET || 'http://localhost:8080';
const PAT = __ENV.PAT;
const ORG = __ENV.ORG || 'acme';

function authedHeaders(extra) {
  return Object.assign(
    {
      Authorization: `Bearer ${PAT}`,
      'Content-Type': 'application/json',
      'X-Idempotency-Key': `${__VU}-${Date.now()}-${Math.random()}`,
    },
    extra || {}
  );
}

export function readFlow() {
  group('list repos', () => {
    const r = http.get(`${TARGET}/api/v1/orgs/${ORG}/repos?page_size=25`, {
      headers: authedHeaders(),
      tags: { op: 'read' },
    });
    check(r, { 'status 200': (x) => x.status === 200 }) || errors.add(1);
    readLatency.add(r.timings.duration);
  });

  group('get a repo', () => {
    const r = http.get(`${TARGET}/api/v1/orgs/${ORG}/repos/demo`, {
      headers: authedHeaders(),
      tags: { op: 'read' },
    });
    check(r, { 'status 200': (x) => x.status === 200 }) || errors.add(1);
    readLatency.add(r.timings.duration);
  });

  group('list PRs', () => {
    const r = http.get(`${TARGET}/api/v1/orgs/${ORG}/repos/demo/prs?state=open`, {
      headers: authedHeaders(),
      tags: { op: 'read' },
    });
    check(r, { 'status 200': (x) => x.status === 200 }) || errors.add(1);
    readLatency.add(r.timings.duration);
  });

  group('search', () => {
    const q = `q=${encodeURIComponent(pick(['deploy','api','ci','auth','repo','storage']))}`;
    const r = http.get(`${TARGET}/api/v1/search/repos?${q}&limit=20`, {
      headers: authedHeaders(),
      tags: { op: 'read' },
    });
    check(r, { 'status 200': (x) => x.status === 200 }) || errors.add(1);
    readLatency.add(r.timings.duration);
  });

  sleep(randomIntBetween(1, 3));
}

export function writeFlow() {
  const slug = `loadtest-${__VU}-${__ITER}`;

  group('create issue', () => {
    const r = http.post(
      `${TARGET}/api/v1/orgs/${ORG}/repos/demo/issues`,
      JSON.stringify({
        title: `Load test issue ${slug}`,
        body: `Created at ${new Date().toISOString()} by k6 VU ${__VU}`,
        labels: ['load-test'],
      }),
      { headers: authedHeaders(), tags: { op: 'write' } }
    );
    check(r, { 'status 201': (x) => x.status === 201 }) || errors.add(1);
    writeLatency.add(r.timings.duration);
    if (r.status === 201) createdIssues.add(1);
  });

  if (__ITER % 10 === 0) {
    group('create repo', () => {
      const r = http.post(
        `${TARGET}/api/v1/orgs/${ORG}/repos`,
        JSON.stringify({
          slug: slug,
          display_name: `Load ${slug}`,
          visibility: 'internal',
          auto_bind_all_enabled_upstreams: false,
        }),
        { headers: authedHeaders(), tags: { op: 'write' } }
      );
      check(r, { 'status in 2xx': (x) => x.status < 300 }) || errors.add(1);
      writeLatency.add(r.timings.duration);
      if (r.status < 300) createdRepos.add(1);
    });
  }

  sleep(randomIntBetween(2, 5));
}

export function subscribeFlow() {
  const url = `${TARGET.replace(/^http/, 'ws')}/events/v1?resume_token=`;
  // Simplified: for true WS load, use k6 websocket module with
  // import ws from 'k6/ws'; kept here as HTTP SSE fallback for demo.
  const r = http.get(
    `${TARGET}/api/v1/events/stream?scope=org:${ORG}&accept=text/event-stream`,
    { headers: authedHeaders({ Accept: 'text/event-stream' }), timeout: '60s' }
  );
  check(r, { 'stream 200': (x) => x.status === 200 }) || errors.add(1);
  sleep(30);
}

function pick(arr) { return arr[randomIntBetween(0, arr.length - 1)]; }

export function handleSummary(data) {
  const file = __ENV.SUMMARY_OUT || 'summary.json';
  return { [file]: JSON.stringify(data, null, 2), stdout: renderText(data) };
}

function renderText(data) {
  const m = data.metrics;
  return [
    '\n=== HelixGitpx k6 summary ===',
    `Reads p99:  ${fmt(m['read_latency_ms'] && m['read_latency_ms'].values['p(99)'])} ms`,
    `Writes p99: ${fmt(m['write_latency_ms'] && m['write_latency_ms'].values['p(99)'])} ms`,
    `Errors:     ${m.errors ? m.errors.values.count : 0}`,
    `Fail-rate:  ${(m.http_req_failed.values.rate * 100).toFixed(3)} %`,
    '',
  ].join('\n');
}
function fmt(v) { return v ? v.toFixed(1) : '—'; }
