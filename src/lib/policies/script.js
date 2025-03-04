import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 1500 }
  ],
};

const publicPolicyUrl = 'http://opa:8181/v1/data/tests/public/allow';
const zerotrustPolicyUrl = 'http://opa:8181/v1/data/tests/zerotrust/allow';

export default function () {
  // Define test inputs for OPA
  const payload = JSON.stringify({
    input: {
      method: 'PUT',
      path: ['pets', '1234'],
      user: 'john_doe',
      owner: 'john_doe',
    },
  });

  const headers = {
    'Content-Type': 'application/json',
  };

  // Test the public policy
  const publicRes = http.post(publicPolicyUrl, payload, { headers });
  check(publicRes, {
    'public policy - status was 200': (r) => r.status === 200,
    'public policy - access allowed': (r) => {
      const result = JSON.parse(r.body);
      return result.result === true; // Adjust as per expected OPA response
    },
  });

  // Test the zerotrust policy
  const zerotrustRes = http.post(zerotrustPolicyUrl, payload, { headers });
  check(zerotrustRes, {
    'zerotrust policy - status was 200': (r) => r.status === 200,
    'zerotrust policy - access allowed': (r) => {
      const result = JSON.parse(r.body);
      return result.result === true; // Adjust as per expected OPA response
    },
  });

  sleep(1);
}
