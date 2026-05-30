import http from 'k6/http';
import {check, fail, sleep} from 'k6';

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
    const tokens = [];
    const login_payload = JSON.stringify({
        email: 'test@gmail.com',
        password: 'test123456',
    });

    const params = {headers: {'Content-Type': 'application/json'}};

    console.log(`Pre-generating ${options.vus} unique sessions for VUs...`);

    for (let i = 0; i < options.vus; i++) {
        const res = http.post(`${BASE_URL}/auth/login`, login_payload, params);

        if (res.status !== 200) {
            fail(`Setup failed at iteration ${i}! Status: ${res.status} (${res.body})`);
        }

        const res_json = res.json();
        if (!res_json || !res_json.data || !res_json.data.refresh_token) {
            fail(`Setup mismatch at iteration ${i}. Body: ${res.body}`);
        }

        tokens.push(res_json.data.refresh_token);
    }

    console.log(`Pre-generated done...`)

    return {tokens: tokens};
}

export default function (data) {
    if (!globalThis.vu_rotated_token) {
        const vu_index = __VU - 1;
        globalThis.vu_rotated_token = data.tokens[vu_index];
    }

    const current_token = globalThis.vu_rotated_token;

    if (!current_token) {
        console.error(`VU ${__VU} has no assigned token context.`);
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
        if (json_body && json_body.data && json_body.data.refresh_token) {
            globalThis.vu_rotated_token = json_body.data.refresh_token;
        }
    } else {
        console.error(`VU ${__VU} Iteration ${__ITER} failed [Status ${res.status}]. Raw Body: ${res.body}`);
        sleep(1);
    }
}