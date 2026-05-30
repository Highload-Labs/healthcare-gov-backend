import http from 'k6/http';
import { check, sleep } from 'k6';

const ZIP_CODES = ['90210', '75001', '10001', '30301', '60601', '33101', '19101', '94101', '02101', '98101'];

export const options = {
    vus: 150,
    duration: '30s',
    summaryTrendStats: ['min', 'avg', 'med', 'max', 'p(50)', 'p(95)', 'p(99)'],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: [
            'p(50)<60',
            'p(95)<150',
            'p(99)<300'
        ]
    }
};

export default function () {
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQGdtYWlsLmNvbSIsInVzZXJuYW1lIjoiZGFybGluZyIsImlzcyI6ImhlYWx0aGNhcmUtZ292Iiwic3ViIjoiOTdlN2YxMzYtMjBkYS00MmViLThmMTYtMmY0MjFhMTRjNWE2IiwiZXhwIjoxNzgwMTM0NTc5LCJpYXQiOjE3ODAxMzA5NzksImp0aSI6IjI0YWMwMzYyLTFhODAtNDg1Yi1iMTFjLTdhZTUxNjVhN2ExMyJ9.ExSpw0MMKU7yEXoeVuOdrA4dcogc5y-nAC7ZXf2D8m8' // Add a valid token here
        }
    };

    // --- STEP 1: BROWSE PLANS BY ZIPCODE ---
    const randomZip = ZIP_CODES[Math.floor(Math.random() * ZIP_CODES.length)];
    const listUrl = `http://localhost:8080/plans?zipcode=${randomZip}`;

    // Using a tag makes it easy to isolate List vs Detail performance in the summary
    const listRes = http.get(listUrl, { ...params, tags: { name: 'GetPlansList' } });

    const listPassed = check(listRes, {
        'list status is 200': (r) => r.status === 200,
        'has plans data': (r) => {
            const body = JSON.parse(r.body);
            return body.success && Array.isArray(body.data) && body.data.length > 0;
        }
    });

    // --- STEP 2: DRILL DOWN TO DETAILS (Only if Step 1 returned plans) ---
    if (listPassed) {
        try {
            const body = JSON.parse(listRes.body);
            const plans = body.data;

            // Pick a random plan ID from the array returned by your actual DB
            const randomPlan = plans[Math.floor(Math.random() * plans.length)];
            const planId = randomPlan.id;

            const detailUrl = `http://localhost:8080/plans/${planId}`;
            const detailRes = http.get(detailUrl, { ...params, tags: { name: 'GetPlanDetail' } });

            check(detailRes, {
                'detail status is 200': (r) => r.status === 200
            });
        } catch (err) {
            // Catch parsing errors just in case an unexpected HTML/Error response slips through
            console.error(`Failed to parse response JSON: ${err}`);
        }
    }

}