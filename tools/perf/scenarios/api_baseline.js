import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 50 },
    { duration: '5m', target: 200 },
    { duration: '5m', target: 500 },
    { duration: '3m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<100', 'p(99)<200'],
    http_req_failed: ['rate<0.001'],
  },
};

const base = __ENV.HELIXGITPX_API || 'https://api.helixgitpx.local';
const token = __ENV.HELIXGITPX_TOKEN || '';

export default function () {
  const headers = { Authorization: `Bearer ${token}` };
  const r = http.get(`${base}/api/v1/orgs`, { headers });
  check(r, { 'orgs 200': (x) => x.status === 200 });
  sleep(0.5);
}
