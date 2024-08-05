import http from 'k6/http';
import { check, sleep } from 'k6';

const URL = `http://localhost/api/rate`

export const options = {
    stages: [
        {duration: '30s', target: 100},
        {duration: '20s', target: 200},
        {duration: '20s', target: 50}
    ]
};

export default function() {
    let res = http.get(URL);
    check(res, { "status is 200": (res) => res.status === 200 });
    sleep(1);
}
