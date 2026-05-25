import http from 'k6/http';
import { check, sleep, fail } from 'k6';

export const options = {
    vus: 150,
    duration: '30s',
    summaryTrendStats: ['min', 'avg', 'med', 'max', 'p(50)', 'p(95)', 'p(99)'],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: [
            'p(50)<50',
            'p(95)<100',
            'p(99)<250'
        ]
    }
};

const BASE_URL = 'http://localhost:8080';

export function setup() {
    const login_payload = JSON.stringify({
        email: 'test@gmail.com',
        password: 'test123456',
    });

    const params = { headers: { 'Content-Type': 'application/json' } };
    const res = http.post(`${BASE_URL}/auth/login`, login_payload, params);

    if (res.status !== 200) {
        fail(`Setup failed! Unable to login to fetch refresh token. Status: ${res.status}`);
    }

    const res_json = res.json();

    // Adjusted to check for snake_case keys from your Go server response
    if (!res_json || !res_json.data || !res_json.data.refresh_token) {
        fail(`Setup failed! JSON structure mismatch. Body: ${res.body}`);
    }

    const refresh_token = res_json.data.refresh_token;
    return { token: refresh_token };
}

export default function (data) {
    // Resolve token context for current VU
    let current_token = globalThis.vu_rotated_token || (data && data.token);

    if (!current_token) {
        console.error(`VU ${__VU} Iteration ${__ITER} has no available token!`);
        sleep(0.5);
        return;
    }

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${current_token}`
        }
    };

    const res = http.get(`${BASE_URL}/auth/refresh`, params);

    const status_ok = check(res, {
        'status is 200': (r) => r.status === 200,
    });

    if (status_ok) {
        const json_body = res.json();

        // Save the newly rotated snake_case token to globalThis context
        if (json_body && json_body.data && json_body.data.refresh_token) {
            globalThis.vu_rotated_token = json_body.data.refresh_token;
        }
    } else {
        console.error(`VU ${__VU} Iteration ${__ITER} failed with status ${res.status}. Header used: Bearer ${current_token.substring(0, 15)}... Raw Body: ${res.body}`);
    }
}