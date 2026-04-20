import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
  vus: 500,
  duration: '10m',
  thresholds: {
    ws_connecting: ['p(95)<500'],
    ws_session_duration: ['p(99)<600000'],
  },
};

const base = __ENV.HELIXGITPX_WS || 'wss://live.helixgitpx.local/ws';

export default function () {
  const res = ws.connect(base, {}, (socket) => {
    socket.on('open', () => socket.send(JSON.stringify({ subscribe: 'org.*.repo.*.events' })));
    socket.setTimeout(() => socket.close(), 60000);
  });
  check(res, { 'ws 101': (r) => r && r.status === 101 });
}
