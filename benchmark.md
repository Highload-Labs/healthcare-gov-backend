Environment:
- EC2: t3.small
- Postgresql-16
- Redis-8.8.0
- Nginx-1.13.1
- 3 Backend Instance
- Region AP-Southeast-1
- Duration 30s
- Tool k6

## Register
100VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=28.83ms avg=40.35ms med=37.87ms max=701ms    p(50)=37.87ms p(95)=51.08ms p(99)=97.15ms
http_reqs...............................................................: 74019 2464.497292/s
```

150VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=29.09ms avg=53.11ms med=46.6ms  max=1.34s p(50)=46.6ms  p(95)=79.97ms p(99)=146.08ms
http_reqs...............................................................: 84459 2810.622048/s
```

Postgres Profiling Result:
```bash
healthcare_gov=# SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time FROM pg_stat_statements;
query                                      | calls |  total_exec_time  |   mean_exec_time    |   max_exec_time    
--------------------------------------------------------------------------------+-------+-------------------+---------------------+--------------------
SELECT id, email, username, password  FROM users WHERE email = $1              | 74019 | 4261.065569999972 | 0.05756718639808657 | 58.051632999999995
INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id | 74019 | 6695.518750999987 | 0.09045675773787683 |          64.153558
SELECT pg_stat_statements_reset()                                              |     1 |          0.094362 |            0.094362 |           0.094362
(3 rows)
```

## Refresh
100VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=26.26ms avg=31.27ms med=30.11ms max=181.51ms p(50)=30.11ms p(95)=36.57ms p(99)=55.73ms
http_reqs...............................................................: 95685 2878.380704/s
```

150VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=26.57ms avg=34.04ms med=32.22ms max=496.42ms p(50)=32.22ms p(95)=42.22ms p(99)=54.91ms
http_reqs...............................................................: 131898 3795.742677/s
```

Postgres Profiling Result:
```bash
healthcare_gov=# SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time FROM pg_stat_statements;
query                               | calls |  total_exec_time  |   mean_exec_time    |    max_exec_time    
-------------------------------------------------------------------+-------+-------------------+---------------------+---------------------
SELECT id, email, username, password  FROM users WHERE email = $1 |   100 |          6.248554 | 0.06248554000000002 | 0.14485499999999998
SELECT id, email, username, password  FROM users WHERE id = $1    | 95585 | 5976.758181000012 | 0.06252820192498788 |           35.930399
SELECT pg_stat_statements_reset()                                 |     1 |          0.105191 |            0.105191 |            0.105191
(3 rows)
```

Observations:
Needs to re-design and make it 100% redis, but low priority EP and not on Load Model EP

## Coverage
400VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=25.1ms avg=31.47ms med=29.87ms max=265.81ms p(50)=29.87ms p(95)=40.32ms p(99)=57.48ms
http_reqs...............................................................: 379788 12636.399823/s
✓ status is 200
```

500VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=25.33ms avg=37.6ms  med=34.44ms max=467.7ms p(50)=34.44ms p(95)=48.85ms p(99)=87.03ms
http_reqs...............................................................: 397373 13231.735426/s
✗ status is 200
↳  99% — ✓ 397372 / ✗ 1
```

Postgres Profiling Result:
```bash
healthcare_gov=# SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time FROM pg_stat_statements;
query               | calls | total_exec_time | mean_exec_time | max_exec_time
-----------------------------------+-------+-----------------+----------------+---------------
SELECT pg_stat_statements_reset() |     1 |        0.109637 |       0.109637 |      0.109637
(1 row)
```

Observations:
Nginx may be the causes of 1 dropped requests. As noticed on the metrics, Nginx node has 96.6% with 500VUs and 92.8% with 400VUs. Therefore, Nginx node is highly speculated to become the bottleneck.

## Get All Plans:
100VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=26.96ms avg=35.98ms med=34.41ms max=134.64ms p(50)=34.41ms p(95)=46.82ms p(99)=61.15ms
http_reqs...............................................................: 83085 2757.763785/s
✓ list status is 200
✓ has plans data
```

150VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=28.9ms  avg=59.72ms med=50.44ms max=680.17ms p(50)=50.44ms p(95)=109.01ms p(99)=154.42ms
http_reqs...............................................................: 75190 2496.454455/s
✓ list status is 200
✓ has plans
```

Postgres Profiling Result:
```bash
healthcare_gov=# SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time FROM pg_stat_statements;
query                                                                                     | calls |   total_exec_time    |    mean_exec_time    |    max_exec_time     
-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+-------+----------------------+----------------------+----------------------
SELECT id, email, username, password  FROM users WHERE email = $1                                                                                                             |     1 | 0.047018000000000004 | 0.047018000000000004 | 0.047018000000000004
SELECT COUNT(*) as total_data FROM plans WHERE state = $1                                                                                                                     | 75189 |   16209.542792999993 |  0.21558396564656926 |            21.689046
SELECT id, name, provider, tier, monthly_premium, deductible, out_of_pocket_max, state, created_at, updated_at FROM plans WHERE state = $1 ORDER BY id ASC LIMIT $2 OFFSET $3 | 75189 |   16660.423648000113 |  0.22158059886419787 |            66.247581
SELECT pg_stat_statements_reset()                                                                                                                                             |     1 |             0.097913 |             0.097913 |             0.097913
(4 rows)
```

Observations:
Notice that despite VUs increased, the RPS going down and increase on latency. Metrics also showed that Backend 1,2,3 Node CPU usage still below 50% while Postgres & Redis node at 100% CPU Usage. This is likely the sign that Postgres becomes the bottleneck.

## Get Plan By ID:
500VUs:
```bash
✓ 'rate<0.01' rate=0.00%
http_req_duration.......................................................: min=25.65ms avg=46.92ms med=43.24ms max=320.99ms p(50)=43.24ms p(95)=63.05ms p(99)=167.92ms
http_reqs...............................................................: 317954 10559.324553/s
✓ detail status is 200
✓ has valid response
```

600VUs:
```bash
✓ 'rate<0.01' rate=0.08%
http_req_duration.......................................................: min=25.44ms avg=55.49ms med=51.85ms max=1.76s p(50)=51.85ms p(95)=62.21ms p(99)=123.13ms
http_reqs...............................................................: 322819 10717.871135/s
✗ detail status is 200
↳  99% — ✓ 322550 / ✗ 268
✗ has valid response
↳  99% — ✓ 322550 / ✗ 268
```

Postgres Profiling Result:
```bash
healthcare_gov=# SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time FROM pg_stat_statements;
query                               | calls | total_exec_time | mean_exec_time | max_exec_time
-------------------------------------------------------------------+-------+-----------------+----------------+---------------
SELECT id, email, username, password  FROM users WHERE email = $1 |     1 |        0.053625 |       0.053625 |      0.053625
SELECT pg_stat_statements_reset()                                 |     1 |        0.145615 |       0.145615 |      0.145615
(2 rows)
```

Observations:
Despite VUs increased, the RPS is not meaningfully increased and error rate increased from 0.00% to 0.08%. Metrics also showed that Backend 1,2,3 still below 80% and Postgres & Redis Node also stay at 24.6%. The one that has highest peak was Nginx & Monitoring Node. This is may indicate that Nginx host became the bottleneck.

## Enrollments
Please refer to [this pr comment](https://github.com/Highload-Labs/healthcare-gov-backend/pull/43#issuecomment-4587519213) for concurrent enrollment correctness. The Benchmark was produced with local environment which is powered by Ryzen 7 5700X and benchmarks above were produced with EC2 t3.small.