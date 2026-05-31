import http from 'k6/http';
import {check} from 'k6';

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
}

export default function () {
    const current_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQGdtYWlsLmNvbSIsInVzZXJuYW1lIjoiZGFybGluZyIsImlzcyI6ImhlYWx0aGNhcmUtZ292Iiwic3ViIjoiOTdlN2YxMzYtMjBkYS00MmViLThmMTYtMmY0MjFhMTRjNWE2IiwiZXhwIjoxNzgwMjQ1NTE4LCJpYXQiOjE3ODAyNDE5MTgsImp0aSI6IjA4ZGUxMDdlLTU3YzMtNDdiNS04MDJjLTNlYWZhNDJkYTM1ZCJ9.C1T3lynX02oaKvk3v2hYRX23p_TYkiQB0LfM9rtYMZM"

    const payload = JSON.stringify({
        plan_id: '01107fe3-791c-48ae-a027-c360ce9448cf',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${current_token}`
        }
    }

    const res = http.post(
        'http://localhost:8080/enrollments',
        payload,
        params
    );

    check(res, {
        'status is 201': (r) => r.status === 201,
    });
}