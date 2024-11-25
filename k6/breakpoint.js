// https://github.com/grafana/k6-learn
// https://grafana.com/docs/k6/latest/testing-guides/test-types/breakpoint-testing/

// K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_HOST=0.0.0.0 k6 run ./breakpoint.js
// dstat -tcm
// viddy -n 0.5 'curl -s http://127.0.0.1:2112/metrics | grep "requests_in_flight"'

// 發起測試的機器, 本身運行太多軟體, 會影響 target rps 不準確
// 不應該使用開發者生活用的機器執行 k6

import http from 'k6/http';
import {check, sleep} from 'k6';

const BASE_URL = "http://127.0.0.1:8800";

function generateStages(start, end, diff) {
  const stages = [];
  for (let t = start; t <= end; t += diff) {
    stages.push(
      { duration: '5s', target: t }, // up
      { duration: '15s', target: t }, // stable
    );
  }
  return stages;
}

const constantVU = {
  executor: 'constant-vus',
  vus: 10,
  duration: '10s',
};

const perVU = {
  executor: 'per-vu-iterations',
  vus: 1000,
  iterations: 1,
};

const constantRate = {
  executor: 'constant-arrival-rate',
  rate: 3000,
  timeUnit: '1s',
  duration: '10s',
  preAllocatedVUs: 3000,
  // preAllocatedVUs: 2,
  maxVUs: 3100,
};

const findSystemLimitByVU = {
  executor: 'ramping-vus',
  startVUs: 1500,
  stages: generateStages(1500, 5000, 500),
};

const findSystemLimitByRate = {
  executor: 'ramping-arrival-rate',
  startRate: 2000,
  timeUnit: '1s',
  preAllocatedVUs: 2000,
  maxVUs: 6100,
  stages: generateStages(2000, 6000, 500),
};

export let options = {
  discardResponseBodies: true,
  thresholds: {
    'http_req_waiting': [
      // {abortOnFail: true, threshold: 'p(99) < 500'},
    ],
  },

  scenarios: {
    constant_VU: constantVU,
    // per_VU: perVU,
    // constant_rate: constantRate,
    // find_system_limit_by_vu: findSystemLimitByVU,
    // find_system_limit_by_rate: findSystemLimitByRate
  },
};

const endpoints = [
  {method: 'GET', url: '/logger/level?level=debug'},
  {
    method: 'POST',
    url: '/logger/level',
    body: JSON.stringify({level: "debug"})
  },
];

function getRandomIndex(length) {
  return Math.floor(Math.random() * length)
}

function getRandomEndpoint() {
  const endpoint = endpoints[getRandomIndex(endpoints.length)];
  let url = BASE_URL + endpoint.url; // Add base URL to the endpoint

  return {
    method: endpoint.method,
    url: url,
    body: undefined,
  };
}

export default function () {
  const {method, url, body} = getRandomEndpoint();
  let resp;

  if (method === 'GET') {
    resp = http.get(url);
  } else if (method === 'POST') {
    resp = http.post(url, body, {headers: {'Content-Type': 'application/json'}});
  }

  check(resp, {
    'http status < 500': (resp) => resp.status < 500,
    'latency < 500ms': (resp) => resp.timings.duration < 500,
  });
  sleep(1);
}
