import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
    vus: 10,
    duration: '30s',
};
export default function () {
    let data = {
        "name": "test",
        "surname": "test",
        "age": 20,
        "detail": "try to insert with time"
    }
    http.post('http://localhost:8080/demo', JSON.stringify(data), {
        headers: {'Content-Type': 'application/json'}
    });
    // check(res, { 'status was success': (r) => r.status >= 200 && r.status <= 299 });
}