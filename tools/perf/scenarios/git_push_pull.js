import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    clone: { executor: 'ramping-vus', startVUs: 0, stages: [
      { duration: '1m', target: 20 },
      { duration: '5m', target: 100 },
      { duration: '1m', target: 0 },
    ]},
  },
  thresholds: {
    http_req_duration: ['p(95)<1500'],
    http_req_failed: ['rate<0.005'],
  },
};

const base = __ENV.HELIXGITPX_GIT || 'https://git.helixgitpx.local';

export default function () {
  const r = http.get(`${base}/testorg/testrepo.git/info/refs?service=git-upload-pack`);
  check(r, { 'smart refs 200': (x) => x.status === 200 });
}
