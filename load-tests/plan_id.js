import http from 'k6/http';
import { check, sleep, fail } from 'k6';

const BASE_URL = 'http://13.212.195.147';
const PLAN_IDS = ['8db0f547-a5e6-474e-87fb-91ec665806aa'];

export const options = {
    vus: 600,
    duration: '30s',
    summaryTrendStats: ['min', 'avg', 'med', 'max', 'p(50)', 'p(95)', 'p(99)'],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        'http_req_duration{name:GetPlanDetail}': ['p(95)<150'],
    }
};

export function setup() {
    const login_payload = JSON.stringify({
        email: 'test2@gmail.com',
        password: 'test123456',
    });

    const params = { headers: { 'Content-Type': 'application/json' } };

    console.log('Logging in once to retrieve shared token...');
    const res = http.post(`${BASE_URL}/auth/login`, login_payload, params);

    if (res.status !== 200) {
        fail(`Setup failed! Status: ${res.status} (${res.body})`);
    }

    const res_json = res.json();
    const token = res_json.data?.access_token || res_json.data?.refresh_token;

    if (!token) {
        fail(`Setup mismatch. Token not found in response body: ${res.body}`);
    }

    return { authToken: token };
}

export default function (data) {
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${data.authToken}` // Reusing the setup token
        }
    };

    const randomPlanId = PLAN_IDS[Math.floor(Math.random() * PLAN_IDS.length)];
    const detailUrl = `${BASE_URL}/plans/${randomPlanId}`;

    const detailRes = http.get(detailUrl, { ...params, tags: { name: 'GetPlanDetail' } });

    check(detailRes, {
        'detail status is 200': (r) => r.status === 200,
        'has valid response': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.hasOwnProperty('data') || body.success === true;
            } catch (e) {
                return false;
            }
        }
    });
}