import { check, sleep } from 'k6';
import http from 'k6/http';

const BASE_URL = 'http://localhost:3000';

export const options = {
  stages: [
    { duration: '1m', target: 10 },  // Ramp-up to 10 virtual users over 1 minute
    { duration: '3m', target: 10 },  // Maintain 10 virtual users for 3 minutes
    { duration: '1m', target: 0 },   // Ramp-down to 0 virtual users over 1 minute
  ],
};

function randomUsername() {
  return 'user' + Math.floor(Math.random() * 1000);
}

function caseCreateUser() {
  const username = randomUsername();
  const createUserPayload = {
    username: username,
    password: 'password',
    email: `${username}@example.com`,
  };
  const createUserResponse = http.post(`${BASE_URL}/createUser`, JSON.stringify(createUserPayload), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(createUserResponse, {
    'Create User status is 200': (r) => r.status === 200,
  });
  sleep(1);
}

export default function () {
  caseCreateUser()
}
