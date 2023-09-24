import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

import exec from 'k6/execution';
import { SharedArray } from "k6/data";
import { users } from "./data/users.js";

const BASE_URL = 'http://localhost:3000';

const mockUsers = new SharedArray("users", function () {
  return users;
});

export const options = {
  scenarios: {
    contacts: {
      executor: 'per-vu-iterations',
      vus: 10,
      iterations: mockUsers.length,
    },
  },
};
// 10 virtual users in each iteration

function caseDeleteUser() {
  const generateTokenResponse = http.get(`${BASE_URL}/generateToken`);
  if (!mockUsers[exec.scenario.iterationInTest]) {
    return
  }
  const userId = mockUsers[exec.scenario.iterationInTest].id;
  const token = JSON.parse(generateTokenResponse.body).token;
  const deleteUserPayload = {
    userId: userId,
    token: token,
  };
  const deleteUserResponse = http.post(`${BASE_URL}/deleteUser`, JSON.stringify(deleteUserPayload), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(deleteUserResponse, {
    'Delete User status is 200': (r) => r.status === 200,
  });
  sleep(1);
}

export default function () {
  caseDeleteUser()
}
