import http from 'k6/http';
import { check, sleep } from 'k6';

function makeEmail(length) {
  let result = '';
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}

const URL = `http://localhost/api/subscribe`

export const options = {
  stages: [
    {duration: '30s', target: 100},
    {duration: '1m30s', target: 200},
    {duration: '20s', target: 50}
  ]
};

export default function() {
  let payload = {
    "email": `${makeEmail(5)}@gmail.com`
  }
  let res = http.post(URL, JSON.stringify(payload), {headers:{
      'Content-Type': 'application/json'
    }});
  check(res, { "status is 200": (res) => res.status === 200 });
  sleep(1);
}
