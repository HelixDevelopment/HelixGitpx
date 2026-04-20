// load-tests/git-push-fanout.js
// Measures end-to-end push-to-replicated latency. Uses internal test
// endpoint /internal/testing/git/push-dummy which simulates a push
// without shelling out to real git. Authoritative latency is
// measured via subscription to the ref.updated events.
//
// Expected targets:
//   p99 push accept latency     ≤ 500 ms
//   p99 time to first upstream  ≤ 5 s
//   p99 time to all enabled     ≤ 15 s
//   error rate                  < 0.1 %
//
// Run:  k6 run --vus 50 --duration 10m tests/load/git-push-fanout.js
import http from 'k6/http';
import { Trend, Counter } from 'k6/metrics';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const TARGET = __ENV.TARGET || 'http://localhost:8080';
const PAT = __ENV.PAT;
const ORG = __ENV.ORG || 'acme';
const REPO = __ENV.REPO || 'bench';
const UPSTREAMS = parseInt(__ENV.UPSTREAMS || '3');

const pushLatency = new Trend('push_accept_ms', true);
const firstRepLatency = new Trend('first_replicated_ms', true);
const allRepLatency = new Trend('all_replicated_ms', true);
const errors = new Counter('errors');

export const options = {
  thresholds: {
    push_accept_ms: ['p(99)<500'],
    first_replicated_ms: ['p(99)<5000'],
    all_replicated_ms: ['p(99)<15000'],
    errors: ['count<10'],
  },
};

export default function () {
  const ref = `refs/heads/bench-${__VU}`;
  const idempotency = uuidv4();
  const start = Date.now();

  const res = http.post(
    `${TARGET}/internal/testing/git/push-dummy`,
    JSON.stringify({
      org: ORG,
      repo: REPO,
      ref: ref,
      new_sha: fakeSha(),
      old_sha: fakeSha(),
      fan_out_to: UPSTREAMS,
    }),
    {
      headers: {
        Authorization: `Bearer ${PAT}`,
        'Content-Type': 'application/json',
        'X-Idempotency-Key': idempotency,
      },
      timeout: '30s',
    }
  );

  const accept = Date.now() - start;
  pushLatency.add(accept);

  const ok = check(res, { 'push accepted': (r) => r.status === 202 });
  if (!ok) {
    errors.add(1);
    return;
  }

  // Poll the job status until completion (also exposes step-level timings).
  const jobId = res.json('job_id');
  const pollUntil = Date.now() + 20_000;
  let firstSeen = null;
  let allDone = null;

  while (Date.now() < pollUntil) {
    const s = http.get(`${TARGET}/api/v1/sync/jobs/${jobId}`, {
      headers: { Authorization: `Bearer ${PAT}` },
      timeout: '5s',
    });
    if (s.status !== 200) break;
    const body = s.json();
    if (!firstSeen && body.steps && body.steps.some((x) => x.status === 'succeeded')) {
      firstSeen = Date.now() - start;
      firstRepLatency.add(firstSeen);
    }
    if (body.status === 'succeeded' || body.status === 'partial') {
      allDone = Date.now() - start;
      allRepLatency.add(allDone);
      break;
    }
    sleep(0.25);
  }

  if (firstSeen === null) errors.add(1);
}

function fakeSha() {
  const hex = '0123456789abcdef';
  let s = '';
  for (let i = 0; i < 40; i++) s += hex[Math.floor(Math.random() * 16)];
  return s;
}
