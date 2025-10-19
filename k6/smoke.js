import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 10,
  duration: '10s',
};

const BASE = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const resMovies = http.get(`${BASE}/movies`);
  check(resMovies, {
    'GET /movies is 200': (r) => r.status === 200,
  });

  //404
  const resNotFound = http.get(`${BASE}/movies/999999`);
  check(resNotFound, {
    'GET /movies/999999 is 404': (r) => r.status === 404,
  });

  sleep(1);
}
