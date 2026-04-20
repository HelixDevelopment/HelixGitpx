import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 10 },
    { duration: '5m', target: 50 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    'http_req_duration{endpoint:chat}': ['p(95)<8000'],
    http_req_failed: ['rate<0.01'],
  },
};

const base = __ENV.HELIXGITPX_API || 'https://api.helixgitpx.local';
const token = __ENV.HELIXGITPX_TOKEN || '';

export default function () {
  const body = JSON.stringify({
    messages: [{ role: 'user', content: 'Summarize PR #42' }],
  });
  const r = http.post(`${base}/api/v1/ai/chat`, body, {
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    tags: { endpoint: 'chat' },
  });
  check(r, { 'chat 200': (x) => x.status === 200 });
  sleep(1);
}
