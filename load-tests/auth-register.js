import http from 'k6/http';
import {check} from 'k6';
import {uuidv4} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    vus: 150,
    duration: '30s',
    summaryTrendStats: ['min', 'avg', 'med', 'max', 'p(50)', 'p(95)', 'p(99)'],
    thresholds: {
        http_req_failed: ['rate<0.01'], // Require less than 1% failure rate (e.g., no 500s or unintended 400s)
        http_req_duration: [
            'p(50)<100', // Bcrypt password hashing takes CPU cycles; slightly relaxed targets are normal here
            'p(95)<300',
            'p(99)<500'
        ]
    }
};

const BASE_URL = 'http://localhost:8080';

export default function () {
    // 1. Generate high-entropy unique identifiers per execution iteration
    const uniqueId = uuidv4().substring(0, 8);
    const timestamp = Date.now();

    // Generates completely unique payloads like: user_abcd1234_1716654321@healthcare.gov
    const email = `user_${uniqueId}_${timestamp}@healthcare.gov`;
    const username = `user_${uniqueId}`;
    const password = 'SecurePassword123!';

    const payload = JSON.stringify({
        email: email,
        username: username,
        password: password,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // 2. Fire the POST execution target
    const res = http.post(`${BASE_URL}/auth/register`, payload, params);

    // 3. Assert your exact response structure mapping rules
    const statusOk = check(res, {
        'status is 201 Created': (r) => r.status === 201,
        'response marks success': (r) => {
            const body = r.json();
            return body && body.success === true;
        },
        'returns valid user_id': (r) => {
            const body = r.json();
            return body && body.data && body.data.user_id !== undefined;
        }
    });

    // Log unexpected non-201 structural bugs straight to your console out channel
    if (!statusOk) {
        console.error(`Execution fault! VU ${__VU} Iteration ${__ITER} yielded Status: ${res.status}. Body: ${res.body}`);
    }
}