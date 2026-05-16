import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 150,
    duration: '30s',

    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: [
            'p(50)<50',
            'p(95)<100',
            'p(99)<250'
        ]
    }
}

export default function () {
    const payload = JSON.stringify({
        email: 'test@gmail.com',
        password: 'test123456',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json'
        }
    }

    const res = http.post(
        'http://localhost:8080/auth/login',
        payload,
        params
    );

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}