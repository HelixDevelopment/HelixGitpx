// E2E smoke test. Drives the public API surfaces. Constitution §II §2 —
// NO mocks. Requires HELIXGITPX_API (reachable target) + a seeded token.
//
// Run with: k6 run test/e2e/api_smoke.js
//
// Exit 0: all smoke checks passed.
// Exit non-zero: any check failed (CI treats that as an e2e regression).

import http from 'k6/http';
import { check, group } from 'k6';

export const options = {
    vus: 1,
    iterations: 1,
    thresholds: {
        checks: ['rate>=0.99'],
        http_req_duration: ['p(95)<2000'],
    },
};

const base = __ENV.HELIXGITPX_API || 'https://staging.helixgitpx.io';
const token = __ENV.HELIXGITPX_E2E_TOKEN;

if (!token) {
    throw new Error('HELIXGITPX_E2E_TOKEN required — real backend, no mocks');
}

export default function () {
    const h = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };

    group('healthz', () => {
        const r = http.get(`${base}/healthz`);
        check(r, { 'healthz 200': (x) => x.status === 200 });
    });

    group('orgs list', () => {
        const r = http.get(`${base}/api/v1/orgs`, { headers: h });
        check(r, { 'orgs 2xx': (x) => x.status >= 200 && x.status < 300 });
    });

    group('whoami', () => {
        const r = http.get(`${base}/api/v1/auth/whoami`, { headers: h });
        check(r, {
            'whoami 2xx': (x) => x.status >= 200 && x.status < 300,
            'whoami returns email': (x) => /"email"\s*:/.test(x.body || ''),
        });
    });
}
