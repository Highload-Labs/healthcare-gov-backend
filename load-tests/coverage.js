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
    const zipcode = "50000"

    const res = http.get(
        `http://13.212.195.147/coverage/${zipcode}`,
    );

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}