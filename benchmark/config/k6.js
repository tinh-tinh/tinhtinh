// k6 Load Testing Scenarios
// Run with: k6 run k6.js

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

// Test scenarios
export const options = {
  scenarios: {
    // Constant load test
    constant_load: {
      executor: 'constant-vus',
      vus: 50,
      duration: '1m',
      tags: { scenario: 'constant' },
    },

    // Ramping VUs test
    ramping_vus: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '30s', target: 200 },
        { duration: '1m', target: 200 },
        { duration: '30s', target: 0 },
      ],
      tags: { scenario: 'ramping' },
    },

    // Spike test
    spike_test: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: 10 },
        { duration: '10s', target: 500 },  // Spike
        { duration: '10s', target: 10 },
        { duration: '10s', target: 0 },
      ],
      tags: { scenario: 'spike' },
    },

    // Stress test
    stress_test: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 100 },
        { duration: '2m', target: 200 },
        { duration: '2m', target: 300 },
        { duration: '2m', target: 400 },
        { duration: '1m', target: 0 },
      ],
      tags: { scenario: 'stress' },
    },

    // Soak test (sustained load)
    soak_test: {
      executor: 'constant-vus',
      vus: 100,
      duration: '10m',
      tags: { scenario: 'soak' },
    },
  },

  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],  // 95% < 500ms, 99% < 1s
    http_req_failed: ['rate<0.01'],  // Error rate < 1%
    errors: ['rate<0.1'],  // Custom error rate < 10%
  },
};

// Test data
const testData = {
  id: 1,
  name: 'Test User',
  email: 'test@example.com',
  message: 'This is a test message for load testing',
};

export default function () {
  // Test 1: Simple GET request
  let res1 = http.get(`${BASE_URL}/`);
  check(res1, {
    'GET / status is 200': (r) => r.status === 200,
    'GET / response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);

  sleep(0.1);

  // Test 2: JSON GET request
  let res2 = http.get(`${BASE_URL}/json`);
  check(res2, {
    'GET /json status is 200': (r) => r.status === 200,
    'GET /json has valid JSON': (r) => {
      try {
        JSON.parse(r.body);
        return true;
      } catch (e) {
        return false;
      }
    },
  }) || errorRate.add(1);

  sleep(0.1);

  // Test 3: JSON POST request
  const payload = JSON.stringify(testData);
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res3 = http.post(`${BASE_URL}/json`, payload, params);
  check(res3, {
    'POST /json status is 200': (r) => r.status === 200,
    'POST /json response time < 300ms': (r) => r.timings.duration < 300,
  }) || errorRate.add(1);

  sleep(0.1);

  // Test 4: Path parameter request
  let res4 = http.get(`${BASE_URL}/user/123`);
  check(res4, {
    'GET /user/:id status is 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(0.1);

  // Test 5: Query parameter request
  let res5 = http.get(`${BASE_URL}/query?name=test`);
  check(res5, {
    'GET /query status is 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(1);
}

// Setup function (runs once before tests)
export function setup() {
  console.log('Starting load tests...');
  console.log(`Base URL: ${BASE_URL}`);
}

// Teardown function (runs once after tests)
export function teardown(data) {
  console.log('Load tests completed');
}

// Handle summary
export function handleSummary(data) {
  return {
    'summary.json': JSON.stringify(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function textSummary(data, options) {
  // Custom summary formatting
  let summary = '\n';
  summary += '='.repeat(50) + '\n';
  summary += 'Load Test Summary\n';
  summary += '='.repeat(50) + '\n\n';

  for (const [name, scenario] of Object.entries(data.metrics)) {
    if (scenario.values) {
      summary += `${name}:\n`;
      summary += `  min: ${scenario.values.min}\n`;
      summary += `  avg: ${scenario.values.avg}\n`;
      summary += `  max: ${scenario.values.max}\n`;
      summary += `  p95: ${scenario.values['p(95)']}\n`;
      summary += `  p99: ${scenario.values['p(99)']}\n\n`;
    }
  }

  return summary;
}
