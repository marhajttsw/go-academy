import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 100,
  iterations: 5000,
};

const BASE = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const createRes = http.post(`${BASE}/movies`, JSON.stringify({ name: 'k6-movie', year: 2024 }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(createRes, { 'POST /movies 201': (r) => r.status === 201 });
  const created = createRes.json();
  const id = created?.id;

  const getRes = http.get(`${BASE}/movies/${id}`);
  check(getRes, { 'GET /movies/{id} 200': (r) => r.status === 200 });

  const updRes = http.put(
    `${BASE}/movies/${id}`,
    JSON.stringify({ name: 'k6-movie-upd', year: 2025 }),
    { headers: { 'Content-Type': 'application/json' } },
  );
  check(updRes, { 'PUT /movies/{id} 200': (r) => r.status === 200 });

  const delRes = http.del(`${BASE}/movies/${id}`);
  check(delRes, { 'DELETE /movies/{id} 204': (r) => r.status === 204 });

  sleep(0.5);
}
